<template>
  <n-card title="菜单管理">
    <template #header-extra>
      <n-button type="primary" ghost @click="openCreate(0)">新增顶级菜单</n-button>
    </template>
    <n-data-table :columns="columns" :data="tree" :loading="loading" default-expand-all />
  </n-card>

  <n-modal v-model:show="showModal" preset="dialog" title="编辑菜单" style="width: 460px">
    <n-form :model="form" label-placement="left" label-width="80">
      <n-form-item label="名称"><n-input v-model:value="form.name" /></n-form-item>
      <n-form-item label="编码" v-if="!editing"><n-input v-model:value="form.code" placeholder="如 menu:xxx" /></n-form-item>
      <n-form-item label="路由"><n-input v-model:value="form.path" placeholder="如 /system/xxx" /></n-form-item>
      <n-form-item label="父级">
        <n-select v-model:value="form.parent_id" :options="parentOptions" clearable placeholder="顶级" />
      </n-form-item>
      <n-form-item label="图标"><n-input v-model:value="form.icon" placeholder="如 User" /></n-form-item>
      <n-form-item label="排序"><n-input-number v-model:value="form.sort" /></n-form-item>
    </n-form>
    <template #action>
      <n-button @click="showModal = false">取消</n-button>
      <n-button type="primary" @click="save">确定</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue'
import { NCard, NSpace, NButton, NDataTable, NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, useMessage, type DataTableColumns } from 'naive-ui'
import { createPermission, deletePermission, listPermissions, updatePermission } from '../../api'
import type { Permission } from '../../api/types'

const message = useMessage()
const loading = ref(false)
const all = ref<Permission[]>([])
const showModal = ref(false)
const editing = ref(false)
const editId = ref(0)
const form = reactive({ name: '', code: '', path: '', parent_id: null as number | null, icon: '', sort: 0 })

async function load() {
  loading.value = true
  try {
    const { data } = await listPermissions({ page: 1, page_size: 200, type: 'menu' })
    all.value = data.data.list
  } finally {
    loading.value = false
  }
}

// 平铺 -> 树（n-data-table rowKey/children）
const tree = ref<any[]>([])
function buildTree() {
  const map = (items: Permission[], parentID: number): any[] =>
    items.filter((p) => p.parent_id === parentID).map((p) => {
      const children = map(items, p.id)
      const n: any = { id: p.id, name: p.name, code: p.code, path: p.path, icon: p.icon, sort: p.sort }
      if (children.length) n.children = children
      return n
    })
  tree.value = map(all.value, 0)
}

const parentOptions = ref<{ label: string; value: number }[]>([])
function loadParents() {
  parentOptions.value = all.value.filter((p) => p.parent_id === 0).map((p) => ({ label: p.name, value: p.id }))
}

function openCreate(parentID: number) {
  editing.value = false
  Object.assign(form, { name: '', code: '', path: '', parent_id: parentID || null, icon: '', sort: 0 })
  showModal.value = true
}

function openEdit(row: any) {
  editing.value = true
  editId.value = row.id
  const p = all.value.find((x) => x.id === row.id)!
  Object.assign(form, { name: p.name, code: p.code, path: p.path, parent_id: p.parent_id || null, icon: p.icon, sort: p.sort })
  showModal.value = true
}

async function save() {
  if (editing.value) {
    await updatePermission(editId.value, { name: form.name, path: form.path, icon: form.icon, sort: form.sort })
  } else {
    await createPermission({ name: form.name, code: form.code, type: 'menu', path: form.path, parent_id: form.parent_id ?? 0, icon: form.icon, sort: form.sort })
  }
  message.success('保存成功（刷新页面后生效路由）')
  showModal.value = false
  await load()
  buildTree()
  loadParents()
}

const columns: DataTableColumns<any> = [
  { title: '名称', key: 'name' },
  { title: '编码', key: 'code' },
  { title: '路由', key: 'path' },
  { title: '图标', key: 'icon', width: 100 },
  { title: '排序', key: 'sort', width: 70 },
  {
    title: '操作', key: 'actions', width: 220,
    render(row) {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { size: 'small', onClick: () => openEdit(row) }, { default: () => '编辑' }),
          h(NButton, { size: 'small', type: 'info', onClick: () => openCreate(row.id) }, { default: () => '加子菜单' }),
          h(NButton, {
            size: 'small', type: 'error',
            onClick: async () => { await deletePermission(row.id); message.success('已删除'); await load(); buildTree(); loadParents() },
          }, { default: () => '删除' }),
        ],
      })
    },
  },
]

onMounted(async () => {
  await load()
  buildTree()
  loadParents()
})
</script>
