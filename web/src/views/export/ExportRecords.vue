<template>
  <n-card>
    <n-data-table :columns="columns" :data="rows" :loading="loading" :pagination="pagination" remote />
  </n-card>
</template>

<script setup lang="ts">
import { computed, h, onMounted, onBeforeUnmount, reactive, ref, type VNode } from 'vue'
import { NButton, NCard, NDataTable, NTag, NTooltip, useMessage, useDialog, type DataTableColumns, type TagProps } from 'naive-ui'
import { renderActions, type TableAction } from '../../utils/tableActions'
import { listExportRecords, getExportBlob, deleteExport } from '../../api'
import { saveBlob, parseDispositionFilename } from '../../utils/download'
import type { ExportRecord } from '../../api/types'

const message = useMessage()
const dialog = useDialog()
const loading = ref(false)
const rows = ref<ExportRecord[]>([])
const downloadingId = ref(0)
const query = reactive({ page: 1, page_size: 10 })

const pagination = reactive({
  page: 1, pageSize: 10, pageCount: 1, showSizePicker: true,
  onChange: (p: number) => { query.page = p; load() },
  onUpdatePageSize: (s: number) => { query.page_size = s; load() },
})

async function load() {
  loading.value = true
  try {
    const { data } = await listExportRecords(query)
    rows.value = data.data.list
    pagination.page = query.page
    pagination.pageSize = query.page_size
    pagination.pageCount = Math.max(1, Math.ceil(data.data.page.total / query.page_size))
  } finally {
    loading.value = false
  }
  syncPolling()
}

// 有进行中任务时每 5s 自动刷新当前页；无进行中任务即停表
let pollTimer: number | undefined
function syncPolling() {
  window.clearInterval(pollTimer)
  pollTimer = undefined
  if (rows.value.some((r) => r.status === 'pending' || r.status === 'running')) {
    pollTimer = window.setInterval(load, 5000)
  }
}

// 文件名优先取 Content-Disposition（RFC5987），取不到用记录名加 .csv 兜底
async function download(row: ExportRecord) {
  downloadingId.value = row.id
  try {
    const resp = await getExportBlob(row.id)
    const filename = parseDispositionFilename(resp.headers['content-disposition']) || `${row.name}.csv`
    saveBlob(resp.data, filename)
  } catch {
    message.error('下载失败')
  } finally {
    downloadingId.value = 0
  }
}

function confirmDelete(row: ExportRecord) {
  dialog.warning({
    title: '删除确认',
    content: `确定删除导出记录「${row.name}」吗？导出文件将一并删除，不可恢复。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteExport(row.id)
        message.success('已删除')
        load()
      } catch (e: any) {
        message.error(e?.response?.data?.msg || '删除失败')
      }
    },
  })
}

const bizNames: Record<string, string> = { user: '用户列表', login_log: '登录日志', op_log: '操作日志' }

const statusType = (s: ExportRecord['status']): TagProps['type'] =>
  s === 'pending' ? 'info' : s === 'running' ? 'warning' : s === 'done' ? 'success' : 'error'
const statusText = (s: ExportRecord['status']) =>
  s === 'pending' ? '排队中' : s === 'running' ? '导出中' : s === 'done' ? '已完成' : '失败'

function fmtSize(n: number) {
  if (!n) return '—'
  if (n < 1024) return `${n} B`
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`
  return `${(n / 1024 / 1024).toFixed(1)} MB`
}

// 状态标签：失败行 hover 显示错误原因
function renderStatus(row: ExportRecord): VNode {
  const tag = h(NTag, { type: statusType(row.status), size: 'small' }, { default: () => statusText(row.status) })
  if (row.status === 'failed' && row.error) {
    return h(NTooltip, null, { trigger: () => tag, default: () => row.error })
  }
  return tag
}

const columns = computed<DataTableColumns<ExportRecord>>(() => [
  { title: '文件名', key: 'name', ellipsis: { tooltip: true } },
  {
    title: '业务类型', key: 'biz', width: 110,
    render: (row) => h(NTag, { size: 'small', bordered: false }, { default: () => bizNames[row.biz] || row.biz }),
  },
  { title: '状态', key: 'status', width: 90, render: renderStatus },
  {
    title: '行数', key: 'rows', width: 120,
    render: (row) =>
      h('span', { style: 'display: inline-flex; align-items: center; gap: 6px' }, [
        String(row.rows),
        row.truncated
          ? h(NTag, { type: 'warning', size: 'tiny', bordered: false }, { default: () => '已截断' })
          : null,
      ]),
  },
  { title: '大小', key: 'size', width: 100, render: (row) => fmtSize(row.size) },
  { title: '创建时间', key: 'created_at', width: 170 },
  { title: '完成时间', key: 'finished_at', width: 170, render: (row) => row.finished_at || '—' },
  {
    title: '操作', key: 'actions', width: 110,
    render(row) {
      const actions: Array<TableAction | VNode> = []
      if (row.status === 'done') {
        actions.push({
          label: downloadingId.value === row.id ? '下载中' : '下载',
          accent: true,
          onClick: () => { if (!downloadingId.value) download(row) },
        })
      }
      actions.push({ label: '删除', danger: true, onClick: () => confirmDelete(row) })
      return renderActions(actions)
    },
  },
])

onMounted(load)
onBeforeUnmount(() => window.clearInterval(pollTimer))
</script>
