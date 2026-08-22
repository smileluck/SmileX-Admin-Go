<template>
  <div class="err-page">
    <div class="err-code mono">404</div>
    <div class="err-title">页面不存在</div>
    <div class="err-desc">你访问的地址不存在或已被移除。</div>
    <n-button type="primary" @click="goHome">返回首页</n-button>
  </div>
</template>

<script setup lang="ts">
import { NButton } from 'naive-ui'
import { useRouter } from 'vue-router'

const router = useRouter()

function goHome() {
  router.push('/')
}
</script>

<style scoped>
.err-page {
  height: 100%;
  min-height: 320px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  text-align: center;
}
.err-code {
  font-size: 72px;
  font-weight: 700;
  line-height: 1;
  color: var(--sx-accent);
  opacity: 0.85;
  margin-bottom: 10px;
}
.err-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--sx-ink);
}
.err-desc {
  font-size: 12px;
  color: var(--sx-muted);
  margin-bottom: 18px;
}
</style>
