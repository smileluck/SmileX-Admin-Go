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
      <n-form-item label="图标">
        <n-input-group>
          <n-input v-model:value="form.icon" placeholder="图标名如 HomeOutline，或图片 URL">
            <template #prefix>
              <IconPreview v-if="form.icon" :icon="form.icon" :size="16" />
            </template>
          </n-input>
          <n-button @click="openPicker">选择图标</n-button>
        </n-input-group>
      </n-form-item>
      <n-form-item label="排序"><n-input-number v-model:value="form.sort" /></n-form-item>
    </n-form>
    <template #action>
      <n-button @click="showModal = false">取消</n-button>
      <n-button type="primary" @click="save">确定</n-button>
    </template>
  </n-modal>

  <!-- 图标选择器：本地图标搜索 + 自定义图片 URL -->
  <n-modal v-model:show="showPicker" preset="card" title="选择图标" style="width: 640px">
    <n-tabs type="line" size="small">
      <n-tab-pane name="local" tab="本地图标">
        <n-input v-model:value="iconSearch" placeholder="搜索图标名，如 Home / Setting / User" clearable style="margin-bottom: 12px" />
        <div class="icon-grid">
          <div v-for="name in filteredIconNames" :key="name" class="icon-cell" :title="name"
            :class="{ active: form.icon === name }" @click="pickIcon(name)">
            <component :is="(icons as any)[name]" />
            <span class="icon-cell-name">{{ name }}</span>
          </div>
        </div>
      </n-tab-pane>
      <n-tab-pane name="url" tab="网络图片">
        <n-space vertical>
          <n-input v-model:value="iconUrl" placeholder="https://example.com/icon.png" clearable />
          <div class="url-preview">
            <span>预览：</span>
            <IconPreview v-if="iconUrl" :icon="iconUrl" :size="20" />
            <span v-else>—</span>
          </div>
          <n-button type="primary" :disabled="!iconUrl" @click="pickIcon(iconUrl)">使用该图片</n-button>
        </n-space>
      </n-tab-pane>
    </n-tabs>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue'
import { NCard, NSpace, NButton, NDataTable, NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, NInputGroup, NTabs, NTabPane, useMessage, type DataTableColumns } from 'naive-ui'
import * as icons from '@vicons/ionicons5'
import { createPermission, deletePermission, listPermissions, updatePermission } from '../../api'
import { renderMenuIcon } from '../../utils/menuIcon'
import type { Permission } from '../../api/types'

// 模板中预览用的函数式组件
const IconPreview = (props: { icon?: string; size?: number }) => renderMenuIcon(props.icon, props.size ?? 16)

const message = useMessage()
const loading = ref(false)
const all = ref<Permission[]>([])
const showModal = ref(false)
const editing = ref(false)
const editId = ref(0)
const form = reactive({ name: '', code: '', path: '', parent_id: null as number | null, icon: '', sort: 0 })

// 图标选择器
const showPicker = ref(false)
const iconSearch = ref('')
const iconUrl = ref('')
const iconNames = Object.keys(icons).filter((k) => k !== 'default')
const filteredIconNames = computed(() => {
  const kw = iconSearch.value.trim().toLowerCase()
  const list = kw ? iconNames.filter((n) => n.toLowerCase().includes(kw)) : iconNames
  return list.slice(0, 200) // 避免一次渲染上千节点
})

function openPicker() {
  iconSearch.value = ''
  iconUrl.value = ''
  showPicker.value = true
}

function pickIcon(v: string) {
  form.icon = v
  showPicker.value = false
}

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
  {
    title: '图标', key: 'icon', width: 100,
    render(row) {
      return h('span', { style: 'display:inline-flex;align-items:center;gap:6px' }, [
        renderMenuIcon(row.icon, 16),
        h('span', { style: 'font-size:12px;color:var(--sx-muted, #888)' }, row.icon || ''),
      ])
    },
  },
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

<style scoped>
.icon-grid {
  max-height: 360px;
  overflow: auto;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(88px, 1fr));
  gap: 6px;
}
.icon-cell {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 4px;
  border: 1px solid var(--sx-line, #e0e0e0);
  border-radius: 6px;
  cursor: pointer;
  font-size: 18px;
  transition: border-color 0.15s ease, background 0.15s ease;
}
.icon-cell:hover {
  border-color: var(--sx-accent, #2f5d50);
}
.icon-cell.active {
  border-color: var(--sx-accent, #2f5d50);
  background: var(--sx-accent-soft, #eef5f0);
}
.icon-cell :deep(svg) {
  width: 20px;
  height: 20px;
}
.icon-cell-name {
  font-size: 10px;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.url-preview {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 28px;
}
</style>
