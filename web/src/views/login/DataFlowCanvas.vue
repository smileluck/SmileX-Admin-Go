<template>
  <canvas ref="canvasRef" class="dataflow-canvas"></canvas>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

/*
 * 抽象网络数据流：静态点线拓扑 + 沿边双向往返的光点。
 * 光点到达节点后消亡，节点短暂增亮并扩散一圈涟漪。
 */
const MINT = '127, 212, 184' // 与本页 seal/pulse 同源的薄荷绿点缀色

interface FlowNode {
  nx: number // 归一化坐标（0~1），容器缩放时无需重建图
  ny: number
  r: number
  phase: number // 漂移相位
  amp: number // 漂移幅度 px
  speed: number // 漂移角速度
  activity: number // 数据到达后的余亮 0~1
}
interface FlowEdge { a: number; b: number }
interface Packet { edge: number; forward: boolean; t: number; speed: number }
interface Ripple { node: number; t: number }

const MAX_PACKETS = 6
const PACKET_SPEED = 105 // px/s
const SPAWN_MIN = 0.5
const SPAWN_MAX = 0.9

const canvasRef = ref<HTMLCanvasElement | null>(null)
let ctx: CanvasRenderingContext2D | null = null
let rafId = 0
let observer: ResizeObserver | null = null
let width = 0
let height = 0
let nodes: FlowNode[] = []
let edges: FlowEdge[] = []
let packets: Packet[] = []
let ripples: Ripple[] = []
let spawnTimer = 0
let lastTime = 0
const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

/* 拒绝采样布点（保证最小间距）+ 每点连 2 个最近邻（去重） */
function buildGraph() {
  const count = 9
  const minDist = 0.24
  const pts: { nx: number; ny: number }[] = []
  let guard = 0
  while (pts.length < count && guard++ < 500) {
    const p = { nx: 0.06 + Math.random() * 0.88, ny: 0.1 + Math.random() * 0.8 }
    if (pts.every(q => Math.hypot(q.nx - p.nx, q.ny - p.ny) > minDist)) pts.push(p)
  }
  nodes = pts.map(p => ({
    nx: p.nx,
    ny: p.ny,
    r: 2.4 + Math.random() * 1.8,
    phase: Math.random() * Math.PI * 2,
    amp: 2 + Math.random() * 2,
    speed: 0.4 + Math.random() * 0.5,
    activity: 0,
  }))

  const seen = new Set<string>()
  edges = []
  nodes.forEach((n, i) => {
    nodes
      .map((q, j) => ({ j, d: Math.hypot(q.nx - n.nx, q.ny - n.ny) }))
      .filter(o => o.j !== i)
      .sort((a, b) => a.d - b.d)
      .slice(0, 2)
      .forEach(o => {
        const key = Math.min(i, o.j) + '-' + Math.max(i, o.j)
        if (!seen.has(key)) {
          seen.add(key)
          edges.push({ a: Math.min(i, o.j), b: Math.max(i, o.j) })
        }
      })
  })
}

/* 节点实时坐标：基准位置 + 微小正弦漂移 */
function nodePos(n: FlowNode, time: number): [number, number] {
  return [
    n.nx * width + Math.sin(time * n.speed + n.phase) * n.amp,
    n.ny * height + Math.cos(time * n.speed * 0.8 + n.phase) * n.amp,
  ]
}

function edgeLenPx(e: FlowEdge, time: number): number {
  const a = nodePos(nodes[e.a], time)
  const b = nodePos(nodes[e.b], time)
  return Math.hypot(b[0] - a[0], b[1] - a[1]) || 1
}

function spawnPacket() {
  if (!edges.length || packets.length >= MAX_PACKETS) return
  packets.push({
    edge: Math.floor(Math.random() * edges.length),
    forward: Math.random() < 0.5,
    t: 0,
    speed: PACKET_SPEED + Math.random() * 30,
  })
}

function draw(time: number) {
  const g = ctx // 局部常量让 TS 的非空收窄能传入下方回调
  if (!g || width === 0 || height === 0) return
  g.clearRect(0, 0, width, height)
  const pos = nodes.map(n => nodePos(n, time))

  // 连线
  g.lineWidth = 1
  g.strokeStyle = `rgba(${MINT}, 0.14)`
  edges.forEach(e => {
    g.beginPath()
    g.moveTo(pos[e.a][0], pos[e.a][1])
    g.lineTo(pos[e.b][0], pos[e.b][1])
    g.stroke()
  })

  // 节点：柔光外圈 + 实心点（activity 增亮）
  nodes.forEach((n, i) => {
    const [x, y] = pos[i]
    g.beginPath()
    g.arc(x, y, n.r * 3, 0, Math.PI * 2)
    g.fillStyle = `rgba(${MINT}, ${0.12 + n.activity * 0.3})`
    g.fill()
    g.beginPath()
    g.arc(x, y, n.r, 0, Math.PI * 2)
    g.fillStyle = `rgba(${MINT}, ${0.55 + n.activity * 0.45})`
    g.fill()
  })

  // 到达涟漪
  ripples.forEach(rp => {
    const [x, y] = pos[rp.node]
    g.beginPath()
    g.arc(x, y, 4 + rp.t * 12, 0, Math.PI * 2)
    g.strokeStyle = `rgba(${MINT}, ${(1 - rp.t) * 0.45})`
    g.stroke()
  })

  // 数据包：短尾迹 + 亮头
  g.lineWidth = 1.5
  packets.forEach(p => {
    const e = edges[p.edge]
    const from = p.forward ? pos[e.a] : pos[e.b]
    const to = p.forward ? pos[e.b] : pos[e.a]
    const hx = from[0] + (to[0] - from[0]) * p.t
    const hy = from[1] + (to[1] - from[1]) * p.t
    const tail = Math.max(0, p.t - 12 / edgeLenPx(e, time))
    g.beginPath()
    g.moveTo(from[0] + (to[0] - from[0]) * tail, from[1] + (to[1] - from[1]) * tail)
    g.lineTo(hx, hy)
    g.strokeStyle = `rgba(${MINT}, 0.4)`
    g.stroke()
    g.beginPath()
    g.arc(hx, hy, 2.2, 0, Math.PI * 2)
    g.fillStyle = `rgba(${MINT}, 0.95)`
    g.fill()
  })
}

function frame(now: number) {
  const dt = Math.min((now - lastTime) / 1000, 0.05)
  lastTime = now
  const time = now / 1000

  spawnTimer -= dt
  if (spawnTimer <= 0) {
    spawnPacket()
    spawnTimer = SPAWN_MIN + Math.random() * (SPAWN_MAX - SPAWN_MIN)
  }

  // 推进数据包；到达则点亮目标节点并生成涟漪
  packets = packets.filter(p => {
    const e = edges[p.edge]
    const dest = p.forward ? e.b : e.a
    p.t += (p.speed * dt) / edgeLenPx(e, time)
    if (p.t >= 1) {
      nodes[dest].activity = 1
      ripples.push({ node: dest, t: 0 })
      return false
    }
    return true
  })
  ripples = ripples.filter(r => (r.t += dt / 0.6) < 1)
  nodes.forEach(n => { n.activity *= Math.exp(-3 * dt) })

  draw(time)
  rafId = requestAnimationFrame(frame)
}

function resize() {
  const canvas = canvasRef.value
  if (!canvas) return
  const rect = canvas.getBoundingClientRect()
  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  width = rect.width
  height = rect.height
  canvas.width = Math.round(width * dpr)
  canvas.height = Math.round(height * dpr)
  ctx?.setTransform(dpr, 0, 0, dpr, 0, 0)
}

/* 标签页隐藏时挂起循环，回来自动续播 */
function onVisibility() {
  if (document.hidden) {
    cancelAnimationFrame(rafId)
    rafId = 0
  } else if (!rafId) {
    lastTime = performance.now()
    rafId = requestAnimationFrame(frame)
  }
}

onMounted(() => {
  const canvas = canvasRef.value
  if (!canvas) return
  ctx = canvas.getContext('2d')
  buildGraph()
  resize()
  observer = new ResizeObserver(() => {
    resize()
    if (reducedMotion) draw(0) // 静态模式：尺寸变化后重绘单帧
  })
  observer.observe(canvas)
  if (reducedMotion) {
    draw(0) // 降级：仅静态点线，不启动循环
  } else {
    lastTime = performance.now()
    rafId = requestAnimationFrame(frame)
    document.addEventListener('visibilitychange', onVisibility)
  }
})

onBeforeUnmount(() => {
  cancelAnimationFrame(rafId)
  observer?.disconnect()
  document.removeEventListener('visibilitychange', onVisibility)
})
</script>

<style scoped>
.dataflow-canvas {
  flex: 1;
  min-height: 0;
  width: 100%;
  display: block;
}
</style>
