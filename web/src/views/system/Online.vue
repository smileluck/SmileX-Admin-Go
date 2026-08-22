<template>
  <!-- 搜索栏独立卡片：可折叠，重置/搜索按钮在卡片右下角 -->
  <SearchCard storage-key="online" @search="load" @reset="resetQuery">
    <n-input v-model:value="query.username" placeholder="用户名" clearable style="width: 180px" @keyup.enter="load" />
    <n-select v-model:value="query.device" :options="deviceOptions" clearable placeholder="设备端" style="width: 140px" />
  </SearchCard>

  <n-card>
    <template #header>
      <div class="page-actions">
        <n-button ghost @click="load">
          <template #icon><n-icon :component="RefreshOutline" /></template>
          刷新
        </n-button>
      </div>
    </template>

    <n-data-table :columns="columns" :data="rows" :loading="loading" :pagination="pagination" remote />
  </n-card>
</template>

<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NEllipsis, NIcon, NInput, NSelect, NTag, useMessage, useDialog, type DataTableColumns } from 'naive-ui'
import { RefreshOutline } from '@vicons/ionicons5'
import { renderActions, type TableAction } from '../../utils/tableActions'
import SearchCard from '../../components/SearchCard.vue'
import { kickOnlineSession, kickUserSessions, listOnlineUsers } from '../../api'
import { useUserStore } from '../../stores/user'
import type { OnlineSession } from '../../api/types'

const message = useMessage()
const dialog = useDialog()
const userStore = useUserStore()

// admin（id=1）超管账号：仅其本人可见操作按钮（与用户管理保护规则一致）
const SUPER_ADMIN_ID = 1
const canOperate = (row: OnlineSession) => row.user_id !== SUPER_ADMIN_ID || userStore.user?.id === SUPER_ADMIN_ID

const loading = ref(false)
const rows = ref<OnlineSession[]>([])
const total = ref(0)
const query = reactive({ username: '', device: null as string | null, page: 1, page_size: 10 })

const deviceOptions = [
  { label: '网页端', value: 'web' },
  { label: '移动端', value: 'app' },
]

const pagination = reactive({
  page: 1, pageSize: 10, pageCount: 1, showSizePicker: true,
  onChange: (p: number) => { query.page = p; load() },
  onUpdatePageSize: (s: number) => { query.page_size = s; load() },
})

async function load() {
  loading.value = true
  try {
    const { data } = await listOnlineUsers({
      page: query.page,
      page_size: query.page_size,
      username: query.username || undefined,
      device: query.device || undefined,
    })
    rows.value = data.data.list
    total.value = data.data.page.total
    pagination.page = query.page
    pagination.pageSize = query.page_size
    pagination.pageCount = Math.max(1, Math.ceil(total.value / query.page_size))
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.username = ''
  query.device = null
  query.page = 1
  load()
}

// 下线单个会话（该用户此设备端）
function confirmKick(row: OnlineSession) {
  const self = row.is_current ? '（这是你当前登录的会话，下线后需重新登录）' : ''
  dialog.warning({
    title: '下线确认',
    content: `确定下线「${row.username}」的${deviceLabel(row.device)}会话吗？${self}`,
    positiveText: '下线',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await kickOnlineSession(row.sid)
        message.success('已下线')
        load()
      } catch (e: any) {
        message.error(e?.response?.data?.msg || '操作失败')
      }
    },
  })
}

// 下线某用户全部端
function confirmKickUser(row: OnlineSession) {
  dialog.warning({
    title: '全部下线确认',
    content: `确定将「${row.username}」的全部在线端（含网页端、移动端）一并下线吗？`,
    positiveText: '全部下线',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await kickUserSessions(row.user_id)
        message.success('已全部下线')
        load()
      } catch (e: any) {
        message.error(e?.response?.data?.msg || '操作失败')
      }
    },
  })
}

function deviceLabel(device: string) {
  return device === 'app' ? '移动端' : '网页端'
}

// 操作列依赖按钮权限，computed 使权限变化后重新渲染
const columns = computed<DataTableColumns<OnlineSession>>(() => [
  {
    title: '用户名', key: 'username', width: 160,
    render: (row) =>
      h('span', { style: 'display:inline-flex;align-items:center;gap:6px' }, [
        row.username,
        row.is_current ? h(NTag, { type: 'primary', size: 'small', bordered: false }, { default: () => '当前会话' }) : null,
      ]),
  },
  { title: '昵称', key: 'nickname', width: 120, render: (row) => row.nickname || '—' },
  {
    title: '设备端', key: 'device', width: 90,
    render: (row) =>
      h(NTag, { type: row.device === 'app' ? 'success' : 'info', size: 'small' }, { default: () => deviceLabel(row.device) }),
  },
  { title: 'IP', key: 'ip', width: 130 },
  {
    title: '设备信息', key: 'user_agent',
    render: (row) => h(NEllipsis, { style: 'max-width: 220px' }, { tooltip: true, default: () => row.user_agent || '—' }),
  },
  { title: '登录时间', key: 'login_at', width: 170 },
  { title: '最近活跃', key: 'last_active_at', width: 170 },
  {
    title: '操作', key: 'actions', width: 150,
    render(row) {
      if (!canOperate(row)) {
        return renderActions([])
      }
      const actions: Array<TableAction> = []
      if (userStore.has('online:kick')) {
        actions.push({ label: '下线', danger: true, onClick: () => confirmKick(row) })
      }
      if (userStore.has('online:kickUser')) {
        actions.push({ label: '全部下线', danger: true, onClick: () => confirmKickUser(row) })
      }
      return renderActions(actions)
    },
  },
])

onMounted(() => { load() })
</script>

<style scoped>
/* 卡头只放操作按钮（页面标题由顶栏展示） */
.page-actions {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}
</style>
