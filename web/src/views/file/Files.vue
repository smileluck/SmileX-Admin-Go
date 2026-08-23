<template>
  <!-- 搜索栏独立卡片：可折叠，重置/搜索按钮在卡片右下角 -->
  <SearchCard storage-key="files" @search="load" @reset="resetQuery">
    <n-input v-model:value="query.name" placeholder="文件名" clearable style="width: 200px" @keyup.enter="load" />
  </SearchCard>

  <n-card>
    <template #header>
      <div class="page-actions">
        <!-- 隐藏原生文件选择框，按钮触发；选择后立即上传 -->
        <input ref="fileInput" type="file" style="display: none" @change="onPick" />
        <n-button type="primary" ghost :loading="uploading" @click="pickFile" v-permission="['file:upload']">上传文件</n-button>
      </div>
    </template>

    <n-data-table :columns="columns" :data="rows" :loading="loading" :pagination="pagination" remote />
  </n-card>

  <!-- 图片预览 -->
  <n-modal v-model:show="showPreview" preset="card" :title="previewName" style="width: 720px">
    <div class="preview-body">
      <n-spin :show="previewLoading">
        <img v-if="previewUrl" :src="previewUrl" class="preview-img" alt="预览" />
      </n-spin>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, h, onMounted, onBeforeUnmount, reactive, ref, type VNode } from 'vue'
import { NCard, NInput, NButton, NDataTable, NModal, NSpin, NTag, useMessage, useDialog, type DataTableColumns } from 'naive-ui'
import { renderActions, type TableAction } from '../../utils/tableActions'
import SearchCard from '../../components/SearchCard.vue'
import { listFiles, uploadFile, deleteFile, getFileBlob } from '../../api'
import { saveBlob } from '../../utils/download'
import { useUserStore } from '../../stores/user'
import type { FileInfo } from '../../api/types'

const message = useMessage()
const dialog = useDialog()
const userStore = useUserStore()
const loading = ref(false)
const uploading = ref(false)
const rows = ref<FileInfo[]>([])
const query = reactive({ name: '', page: 1, page_size: 10 })
const fileInput = ref<HTMLInputElement | null>(null)

const showPreview = ref(false)
const previewLoading = ref(false)
const previewUrl = ref('')
const previewName = ref('')

const pagination = reactive({
  page: 1, pageSize: 10, pageCount: 1, showSizePicker: true,
  onChange: (p: number) => { query.page = p; load() },
  onUpdatePageSize: (s: number) => { query.page_size = s; load() },
})

async function load() {
  loading.value = true
  try {
    const { data } = await listFiles(query)
    rows.value = data.data.list
    pagination.page = query.page
    pagination.pageCount = Math.ceil(data.data.page.total / query.page_size)
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.name = ''
  query.page = 1
  load()
}

function pickFile() {
  fileInput.value?.click()
}

async function onPick(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = '' // 允许重复选择同一文件
  if (!file) return
  uploading.value = true
  try {
    await uploadFile(file)
    message.success('上传成功')
    load()
  } catch (e: any) {
    message.error(e?.response?.data?.msg || '上传失败')
  } finally {
    uploading.value = false
  }
}

function isImage(row: FileInfo) {
  return /^image\/(png|jpe?g|gif|webp|bmp)$/.test(row.content_type)
}

function fmtSize(n: number) {
  if (n < 1024) return `${n} B`
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`
  return `${(n / 1024 / 1024).toFixed(1)} MB`
}

// 下载需带 JWT，不能用 <a href> 直连：blob 取回后触发浏览器保存
async function download(row: FileInfo) {
  try {
    const { data } = await getFileBlob(row.id, true)
    saveBlob(data, row.name)
  } catch (e: any) {
    message.error(e?.response?.data?.msg || '下载失败')
  }
}

async function preview(row: FileInfo) {
  previewName.value = row.name
  showPreview.value = true
  previewLoading.value = true
  try {
    const { data } = await getFileBlob(row.id)
    if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = URL.createObjectURL(data)
  } catch (e: any) {
    showPreview.value = false
    message.error(e?.response?.data?.msg || '预览失败')
  } finally {
    previewLoading.value = false
  }
}

function confirmDelete(row: FileInfo) {
  dialog.warning({
    title: '删除确认',
    content: `确定删除文件「${row.name}」吗？存储对象将一并删除，不可恢复。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteFile(row.id)
        message.success('已删除')
        load()
      } catch (e: any) {
        message.error(e?.response?.data?.msg || '删除失败')
      }
    },
  })
}

const driverNames: Record<string, string> = {
  local: '本地', oss: '阿里云OSS', cos: '腾讯云COS', tos: '火山云TOS', minio: 'MinIO',
}

// 操作列依赖按钮权限，computed 使权限变化后重新渲染
const columns = computed<DataTableColumns<FileInfo>>(() => [
  { title: 'ID', key: 'id', width: 60 },
  { title: '文件名', key: 'name', ellipsis: { tooltip: true } },
  { title: '大小', key: 'size', width: 100, render: (row) => fmtSize(row.size) },
  {
    title: '存储', key: 'driver', width: 110,
    render: (row) => h(NTag, { size: 'small', bordered: false }, { default: () => driverNames[row.driver] || row.driver }),
  },
  { title: '上传人', key: 'uploader_name', width: 110 },
  { title: '上传时间', key: 'created_at', width: 170 },
  {
    title: '操作', key: 'actions', width: 150,
    render(row) {
      const actions: Array<TableAction | VNode> = []
      if (userStore.has('file:view')) {
        if (isImage(row)) actions.push({ label: '预览', onClick: () => preview(row) })
        actions.push({ label: '下载', accent: true, onClick: () => download(row) })
      }
      if (userStore.has('file:delete')) {
        actions.push({ label: '删除', danger: true, onClick: () => confirmDelete(row) })
      }
      return renderActions(actions)
    },
  },
])

onMounted(load)
// 释放预览 blob URL，避免内存泄漏
onBeforeUnmount(() => { if (previewUrl.value) URL.revokeObjectURL(previewUrl.value) })
</script>

<style scoped>
/* 卡头只放操作按钮（页面标题由顶栏展示） */
.page-actions {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}
.preview-body {
  display: flex;
  justify-content: center;
  min-height: 120px;
}
.preview-img {
  max-width: 100%;
  max-height: 70vh;
}
</style>
