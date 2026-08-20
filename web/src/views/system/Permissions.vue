<template>
  <n-card title="权限管理（API 权限）">
    <template #header-extra>
      <n-space>
        <n-button type="primary" ghost @click="openCreate">新增权限</n-button>
        <n-button @click="toggleType">{{ typeFilter === 'api' ? '查看菜单' : '查看 API 权限' }}</n-button>
      </n-space>
    </template>
    <n-data-table :columns="columns" :data="rows" :loading="loading" :pagination="pagination" remote />
  </n-card>

  <n-modal v-model:show="showModal" preset="dialog" :title="editing ? '编辑权限' : '新增权限'" style="width: 460px">
    <n-form :model="form" label-placement="left" label-width="80">
      <n-form-item label="名称"><n-input v-model:value="form.name" /></n-form-item>
      <n-form-item label="编码" v-if="!editing"><n-input v-model:value="form.code" placeholder="如 api:user:list" /></n-form-item>
      <n-form-item label="类型" v-if="!editing">
        <n-radio-group v-model:value="form.type">
          <n-radio value="api">API</n-radio>
          <n-radio value="menu">菜单</n-radio>
        </n-radio-group>
      </n-form-item>
      <n-form-item label="Method" v-if="form.type === 'api'">
        <n-select v-model:value="form.method" :options="methodOptions" />
      </n-form-item>
      <n-form-item :label="form.type === 'api' ? '路径' : '路由'">
        <n-input v-model:value="form.path" :placeholder="form.type === 'api' ? '如 /api/v1/users 或 /api/v1/users/*' : '如 /system/users'" />
      </n-form-item>
      <template v-if="form.type === 'menu'">
        <n-form-item label="父级">
          <n-select v-model:value="form.parent_id" :options="parentOptions" clearable />
        </n-form-item>
        <n-form-item label="图标"><n-input v-model:value="form.icon" /></n-form-item>
        <n-form-item label="排序"><n-input-number v-model:value="form.sort" /></n-form-item>
      </template>
    </n-form>
    <template #action>
      <n-button @click="showModal = false">取消</n-button>
      <n-button type="primary" @click="save">确定</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { h, onMounted, reactive, ref, computed } from 'vue'
import { NCard, NSpace, NButton, NDataTable, NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, NRadioGroup, NRadio, NTag, useMessage, type DataTableColumns } from 'naive-ui'
import { createPermission, deletePermission, listPermissions, updatePermission } from '../../api'
import type { Permission } from '../../api/types'

const message = useMessage()
const loading = ref(false)
const rows = ref<Permission[]>([])
const typeFilter = ref('api')
const query = reactive({ page: 1, page_size: 10, type: 'api' })

const showModal = ref(false)
const editing = ref(false)
const editId = ref(0)
const form = reactive({ name: '', code: '', type: 'api' as 'api' | 'menu', method: 'GET', path: '', parent_id: null as number | null, icon: '', sort: 0 })

const methodOptions = ['GET', 'POST', 'PUT', 'DELETE', '*'].map((m) => ({ label: m, value: m }))
const parentOptions = ref<{ label: string; value: number }[]>([])

const pagination = reactive({
  page: 1, pageSize: 10, pageCount: 1, showSizePicker: true,
  onChange: (p: number) => { query.page = p; load() },
  onUpdatePageSize: (s: number) => { query.page_size = s; load() },
})

async function load() {
  loading.value = true
  try {
    const { data } = await listPermissions(query)
    rows.value = data.data.list
    pagination.page = query.page
    pagination.pageCount = Math.ceil(data.data.page.total / query.page_size)
  } finally {
    loading.value = false
  }
}

async function loadParents() {
  const { data } = await listPermissions({ page: 1, page_size: 100, type: 'menu' })
  parentOptions.value = data.data.list.map((p) => ({ label: p.name, value: p.id }))
}

function toggleType() {
  typeFilter.value = typeFilter.value === 'api' ? 'menu' : 'api'
  query.type = typeFilter.value
  query.page = 1
  load()
}

function openCreate() {
  editing.value = false
  Object.assign(form, { name: '', code: '', type: typeFilter.value, method: 'GET', path: '', parent_id: null, icon: '', sort: 0 })
  showModal.value = true
}

function openEdit(row: Permission) {
  editing.value = true
  editId.value = row.id
  Object.assign(form, { name: row.name, code: row.code, type: row.type, method: row.method || 'GET', path: row.path, parent_id: row.parent_id || null, icon: row.icon, sort: row.sort })
  showModal.value = true
}

async function save() {
  if (editing.value) {
    await updatePermission(editId.value, { name: form.name, method: form.method, path: form.path, icon: form.icon, sort: form.sort })
  } else {
    await createPermission({ ...form, parent_id: form.parent_id ?? 0 })
  }
  message.success('保存成功')
  showModal.value = false
  load()
}

const columns = computed<DataTableColumns<Permission>>(() => [
  { title: 'ID', key: 'id', width: 60 },
  { title: '名称', key: 'name' },
  { title: '编码', key: 'code' },
  { title: '类型', key: 'type', width: 80, render: (row) => h(NTag, { size: 'small', type: row.type === 'api' ? 'info' : 'success' }, { default: () => row.type }) },
  { title: 'Method', key: 'method', width: 80 },
  { title: '路径', key: 'path' },
  ...(typeFilter.value === 'menu' ? [{ title: '父级', key: 'parent_id', width: 70 }] : []),
  ...(typeFilter.value === 'menu' ? [{ title: '排序', key: 'sort', width: 70 }] : []),
  {
    title: '操作', key: 'actions', width: 150,
    render(row) {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { size: 'small', onClick: () => openEdit(row) }, { default: () => '编辑' }),
          h(NButton, {
            size: 'small', type: 'error',
            onClick: async () => { await deletePermission(row.id); message.success('已删除'); load() },
          }, { default: () => '删除' }),
        ],
      })
    },
  },
])

onMounted(() => { load(); loadParents() })
</script>
