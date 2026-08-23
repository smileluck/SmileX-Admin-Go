<template>
  <!-- 搜索栏独立卡片：可折叠，重置/搜索按钮在卡片右下角 -->
  <SearchCard storage-key="loginLogs" @search="load" @reset="resetQuery">
    <n-input v-model:value="query.username" placeholder="用户名" clearable style="width: 160px" @keyup.enter="load" />
    <n-input v-model:value="query.ip" placeholder="IP 地址" clearable style="width: 150px" @keyup.enter="load" />
    <n-select v-model:value="query.status" :options="statusOptions" clearable placeholder="登录结果" style="width: 120px" />
    <n-date-picker v-model:value="range" type="datetimerange" clearable style="width: 340px; max-width: 100%" />
  </SearchCard>

  <n-card>
    <template #header>
      <div class="page-header">
        <span class="retention-hint">{{ retentionHint }}</span>
        <div class="page-actions">
          <n-button ghost :loading="exporting" v-permission="['log:login:export']" @click="doExport">导出</n-button>
          <n-button type="error" ghost v-permission="['log:login:clear']" @click="confirmClear">清理</n-button>
        </div>
      </div>
    </template>

    <n-data-table :columns="columns" :data="rows" :loading="loading" :pagination="pagination" remote />
  </n-card>
</template>

<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NDatePicker, NEllipsis, NInput, NSelect, NTag, useDialog, useMessage, type DataTableColumns } from 'naive-ui'
import SearchCard from '../../components/SearchCard.vue'
import { clearLoginLogs, createExport, listLoginLogs } from '../../api'
import type { LoginLogInfo } from '../../api/types'

const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const rows = ref<LoginLogInfo[]>([])
const retentionDays = ref(0)
const query = reactive({ username: '', ip: '', status: null as number | null, page: 1, page_size: 10 })
// 时间范围（毫秒时间戳二元组，提交时转秒）
const range = ref<[number, number] | null>(null)

const statusOptions = [
  { label: '成功', value: 1 },
  { label: '失败', value: 0 },
]

const pagination = reactive({
  page: 1, pageSize: 10, pageCount: 1, showSizePicker: true,
  onChange: (p: number) => { query.page = p; load() },
  onUpdatePageSize: (s: number) => { query.page_size = s; load() },
})

async function load() {
  loading.value = true
  try {
    const { data } = await listLoginLogs({
      page: query.page,
      page_size: query.page_size,
      username: query.username || undefined,
      ip: query.ip || undefined,
      status: query.status ?? undefined,
      start: range.value ? Math.floor(range.value[0] / 1000) : undefined,
      end: range.value ? Math.floor(range.value[1] / 1000) : undefined,
    })
    rows.value = data.data.list
    retentionDays.value = data.data.retention_days
    pagination.page = query.page
    pagination.pageSize = query.page_size
    pagination.pageCount = Math.max(1, Math.ceil(data.data.page.total / query.page_size))
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.username = ''
  query.ip = ''
  query.status = null
  range.value = null
  query.page = 1
  load()
}

// 异步导出：提交当前过滤条件（与列表查询一致，剔除分页参数）
const exporting = ref(false)
async function doExport() {
  exporting.value = true
  try {
    await createExport('login-logs', {
      username: query.username,
      ip: query.ip,
      status: query.status,
      start: range.value ? Math.floor(range.value[0] / 1000) : undefined,
      end: range.value ? Math.floor(range.value[1] / 1000) : undefined,
    })
    message.success('已加入导出队列，可在右上角导出图标查看进度')
  } catch (e: any) {
    if (e?.response?.status === 429) {
      message.warning('导出任务过多，请稍后再试')
    } else {
      message.error(e?.response?.data?.msg || '导出失败')
    }
  } finally {
    exporting.value = false
  }
}

function confirmClear() {
  // 与后端保留期自动清理同一截止时间（retentionDays 天前）；未启用保留期时退化为清空全部
  const content =
    retentionDays.value > 0
      ? `确定立即清理 ${retentionDays.value} 天前的登录日志吗？与自动保留策略一致，近期日志将保留。`
      : '当前未启用保留期，确定清空全部登录日志吗？该操作不可恢复。'
  dialog.warning({
    title: '清理确认',
    content,
    positiveText: '清理',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const { data } = await clearLoginLogs()
        message.success(`已清理 ${data.data.deleted} 条过期日志`)
        load()
      } catch (e: any) {
        message.error(e?.response?.data?.msg || '清理失败')
      }
    },
  })
}

function deviceLabel(device: string) {
  return device === 'app' ? '移动端' : '网页端'
}

const columns = computed<DataTableColumns<LoginLogInfo>>(() => [
  { title: 'ID', key: 'id', width: 70 },
  { title: '用户名', key: 'username', width: 140, render: (row) => row.username || '—' },
  { title: 'IP', key: 'ip', width: 130 },
  {
    title: '设备端', key: 'device', width: 90,
    render: (row) => (row.device ? h(NTag, { type: row.device === 'app' ? 'success' : 'info', size: 'small' }, { default: () => deviceLabel(row.device) }) : '—'),
  },
  {
    title: '浏览器 / 终端', key: 'user_agent',
    render: (row) => h(NEllipsis, { style: 'max-width: 200px' }, { tooltip: true, default: () => row.user_agent || '—' }),
  },
  {
    title: '结果', key: 'status', width: 80,
    render: (row) => h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' }, { default: () => (row.status === 1 ? '成功' : '失败') }),
  },
  { title: '说明', key: 'msg', render: (row) => row.msg || '—' },
  { title: '时间', key: 'created_at', width: 170 },
])

// 保留策略说明（与后端保留期自动清理共用同一配置）
const retentionHint = computed(() =>
  retentionDays.value > 0 ? `日志自动保留 ${retentionDays.value} 天，超期自动清理` : '日志永久保留（未启用自动清理）',
)

onMounted(load)
</script>

<style scoped>
/* 卡头左侧保留说明、右侧清空按钮（页面标题由顶栏展示） */
.page-header {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.retention-hint {
  font-size: 13px;
  color: var(--sx-muted);
}
.page-actions {
  display: flex;
  gap: 10px;
}
</style>
