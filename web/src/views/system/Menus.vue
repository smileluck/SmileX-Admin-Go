<template>
  <n-card>
    <template #header>
      <div class="page-actions">
        <n-button type="primary" ghost @click="openCreate(0, 'menu')" v-permission="['menu:create']">新增顶级菜单</n-button>
        <n-button type="primary" ghost @click="openCreate(0, 'button')" v-permission="['menu:create']">新增权限点</n-button>
      </div>
    </template>

    <!-- 表格控制行：与用户/角色页的搜索栏同层，保持三页节奏一致 -->
    <div class="search-bar">
      <n-button @click="toggleExpand">{{ expandAll ? '全部收起' : '全部展开' }}</n-button>
    </div>

    <n-data-table :columns="columns" :data="tree" :loading="loading" :row-key="rowKey"
      :expanded-row-keys="expandedKeys" @update:expanded-row-keys="onExpandUpdate" />
  </n-card>

  <!-- 菜单 / 权限点编辑：type 区分表单 -->
  <n-modal v-model:show="showModal" preset="dialog" :title="modalTitle" style="width: 480px">
    <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80">
      <n-form-item label="名称" path="name">
        <n-input v-model:value="form.name" :maxlength="20" show-word-limit placeholder="最多 20 个字符" />
      </n-form-item>
      <n-form-item label="编码" path="code" v-if="!editing">
        <n-input v-model:value="form.code" :maxlength="64" show-word-limit :placeholder="form.type === 'menu' ? '如 menu:xxx' : '如 user:create'" />
      </n-form-item>
      <template v-if="form.type === 'menu'">
        <n-form-item label="路由"><n-input v-model:value="form.path" placeholder="如 /system/xxx" /></n-form-item>
        <n-form-item label="父级">
          <n-tree-select
            v-model:value="form.parent_id" :options="parentOptions"
            clearable placeholder="顶级" key-field="key" label-field="label" children-field="children"
          />
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
      </template>
      <template v-else>
        <n-form-item label="Method">
          <n-select v-model:value="form.method" :options="methodOptions" />
        </n-form-item>
        <n-form-item label="接口路径">
          <n-input v-model:value="form.path" placeholder="选填，如 /api/v1/users/*；留空则仅控制前端显隐" />
        </n-form-item>
        <n-form-item label="所属菜单">
          <n-tree-select
            v-model:value="form.parent_id" :options="parentOptions"
            clearable placeholder="挂载到菜单" key-field="key" label-field="label" children-field="children"
          />
        </n-form-item>
      </template>
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
import { computed, h, onMounted, reactive, ref, type VNode } from 'vue'
import { NCard, NButton, NDataTable, NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, NInputGroup, NTabs, NTabPane, NTreeSelect, NTag, useMessage, useDialog, type DataTableColumns, type FormInst, type FormRules } from 'naive-ui'
import { renderActions, type TableAction } from '../../utils/tableActions'
import * as icons from '@vicons/ionicons5'
import { createPermission, deletePermission, listAllPermissions, updatePermission } from '../../api'
import { renderMenuIcon } from '../../utils/menuIcon'
import { useUserStore } from '../../stores/user'
import type { Permission } from '../../api/types'

// 模板中预览用的函数式组件
const IconPreview = (props: { icon?: string; size?: number }) => renderMenuIcon(props.icon, props.size ?? 16)

const message = useMessage()
const dialog = useDialog()
const userStore = useUserStore()
const loading = ref(false)
const all = ref<Permission[]>([])
const showModal = ref(false)
const editing = ref(false)
const editId = ref(0)
const formRef = ref<FormInst | null>(null)

// 超管通配权限（id=1，code=all）禁止删除
const WILDCARD_PERM_ID = 1

// form.type 区分菜单 / 按钮权限点两种表单
const form = reactive({
  name: '', code: '', type: 'menu' as 'menu' | 'button',
  method: 'GET', path: '', parent_id: null as number | null, icon: '', sort: 0,
})
const methodOptions = ['GET', 'POST', 'PUT', 'DELETE', '*'].map((m) => ({ label: m, value: m }))
const modalTitle = computed(() => `${editing.value ? '编辑' : '新增'}${form.type === 'menu' ? '菜单' : '权限点'}`)

const rules: FormRules = {
  name: [
    { required: true, message: '请输入名称', trigger: ['blur', 'input'] },
    { max: 20, message: '名称不能超过 20 个字符', trigger: ['blur', 'input'] },
  ],
  code: [
    { required: true, message: '请输入编码', trigger: ['blur', 'input'] },
    { max: 64, message: '编码不能超过 64 个字符', trigger: ['blur', 'input'] },
  ],
}

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
    // 树目录需整表构建，走全量接口（page_size=0），分页会截断子节点
    const { data } = await listAllPermissions()
    all.value = data.data.list
  } finally {
    loading.value = false
  }
}

// 平铺 -> 树（菜单 + 按钮权限点统一按 parent_id 组树，n-data-table rowKey/children）
const tree = ref<any[]>([])
function buildTree() {
  const build = (parentID: number): any[] =>
    all.value
      .filter((p) => p.parent_id === parentID)
      .sort((a, b) => a.sort - b.sort)
      .map((p) => {
        const children = build(p.id)
        const n: any = { id: p.id, name: p.name, code: p.code, type: p.type, method: p.method, path: p.path, icon: p.icon, sort: p.sort }
        if (children.length) n.children = children
        return n
      })
  tree.value = build(0)
}

// 树表受控展开：default-expand-all 对异步加载后才渲染的行不生效，改为手动维护展开行键
const rowKey = (row: any) => row.id
const expandedKeys = ref<number[]>([])
const expandAll = ref(true)
function collectParentIds(nodes: any[]): number[] {
  return nodes.flatMap((n) => (n.children?.length ? [n.id, ...collectParentIds(n.children)] : []))
}
function onExpandUpdate(keys: Array<string | number>) {
  expandedKeys.value = keys as number[]
}
function toggleExpand() {
  expandAll.value = !expandAll.value
  expandedKeys.value = expandAll.value ? collectParentIds(tree.value) : []
}

// 父级选项：仅菜单可作父级；编辑菜单时排除自身及其后代（防环）
const parentOptions = ref<any[]>([])
function buildParentOptions(excludeID?: number) {
  const excluded = new Set<number>(excludeID ? [excludeID] : [])
  let changed = true
  while (changed) {
    changed = false
    for (const p of all.value) {
      if (!excluded.has(p.id) && excluded.has(p.parent_id)) {
        excluded.add(p.id)
        changed = true
      }
    }
  }
  const menus = all.value.filter((p) => p.type === 'menu' && !excluded.has(p.id))
  const build = (parentID: number): any[] =>
    menus
      .filter((p) => p.parent_id === parentID)
      .sort((a, b) => a.sort - b.sort)
      .map((p) => {
        const children = build(p.id)
        const n: any = { label: p.name, key: p.id }
        if (children.length) n.children = children
        return n
      })
  parentOptions.value = build(0)
}

function openCreate(parentID: number, type: 'menu' | 'button') {
  editing.value = false
  Object.assign(form, { name: '', code: '', type, method: 'GET', path: '', parent_id: parentID || null, icon: '', sort: 0 })
  buildParentOptions()
  showModal.value = true
}

function openEdit(row: any) {
  editing.value = true
  editId.value = row.id
  const p = all.value.find((x) => x.id === row.id)!
  Object.assign(form, {
    name: p.name, code: p.code, type: p.type as 'menu' | 'button',
    method: p.method || 'GET', path: p.path, parent_id: p.parent_id || null, icon: p.icon, sort: p.sort,
  })
  // 编辑菜单时排除自身及后代，避免把自己挂到子级下
  buildParentOptions(p.type === 'menu' ? p.id : undefined)
  showModal.value = true
}

async function save() {
  try {
    await formRef.value?.validate()
  } catch {
    return // 校验失败，错误已在表单项上展示
  }
  try {
    if (editing.value) {
      await updatePermission(editId.value, {
        name: form.name.trim(),
        method: form.type === 'button' ? form.method : '',
        path: form.path.trim(),
        icon: form.type === 'menu' ? form.icon : '',
        sort: form.sort,
        parent_id: form.parent_id ?? 0,
      })
    } else {
      await createPermission({
        name: form.name.trim(), code: form.code.trim(), type: form.type,
        method: form.type === 'button' ? form.method : '',
        path: form.path.trim(), parent_id: form.parent_id ?? 0,
        icon: form.type === 'menu' ? form.icon : '', sort: form.sort,
      })
    }
    message.success('保存成功（菜单变更刷新页面后生效路由）')
    showModal.value = false
    await refresh()
  } catch (e: any) {
    message.error(e?.response?.data?.msg || '保存失败')
  }
}

async function remove(row: any) {
  if (row.id === WILDCARD_PERM_ID) { message.error('超管通配权限禁止删除'); return }
  if (all.value.some((p) => p.parent_id === row.id)) {
    message.warning('该节点下存在子级，请先删除子级节点')
    return
  }
  dialog.warning({
    title: '删除确认',
    content: `确定删除「${row.name}（${row.code}）」吗？该操作不可恢复。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deletePermission(row.id)
        message.success('已删除')
        await refresh()
      } catch (e: any) {
        message.error(e?.response?.data?.msg || '删除失败')
      }
    },
  })
}

async function refresh() {
  await load()
  buildTree()
  // 数据刷新后按当前展开状态重算行键，保证新增的父节点也默认展开
  if (expandAll.value) expandedKeys.value = collectParentIds(tree.value)
}

// 操作列依赖按钮权限，computed 使权限变化后重新渲染
const columns = computed<DataTableColumns<any>>(() => [
  {
    title: '名称', key: 'name',
    render(row) {
      return h('span', { style: 'display:inline-flex;align-items:center;gap:6px' }, [
        row.type === 'menu' ? renderMenuIcon(row.icon, 16) : null,
        h('span', {}, row.name),
      ])
    },
  },
  {
    title: '类型', key: 'type', width: 80,
    render(row) {
      return h(NTag, { size: 'small', type: row.type === 'menu' ? 'success' : 'info' }, { default: () => (row.type === 'menu' ? '菜单' : '按钮') })
    },
  },
  { title: '编码', key: 'code' },
  {
    title: '路由/接口', key: 'path',
    render(row) {
      if (row.type === 'menu') return row.path || '—'
      return row.method && row.path ? `${row.method}  ${row.path}` : '—'
    },
  },
  { title: '排序', key: 'sort', width: 70 },
  {
    title: '操作', key: 'actions', width: 210,
    render(row) {
      const actions: Array<TableAction | VNode> = []
      if (userStore.has('menu:update')) {
        actions.push({ label: '编辑', accent: true, onClick: () => openEdit(row) })
      }
      if (row.type === 'menu' && userStore.has('menu:create')) {
        actions.push(
          { label: '加子菜单', onClick: () => openCreate(row.id, 'menu') },
          { label: '加权限点', onClick: () => openCreate(row.id, 'button') },
        )
      }
      // 超管通配权限（all）禁止删除，不展示删除按钮
      if (row.id !== WILDCARD_PERM_ID && userStore.has('menu:delete')) {
        actions.push({ label: '删除', danger: true, onClick: () => remove(row) })
      }
      return renderActions(actions)
    },
  },
])

onMounted(refresh)
</script>

<style scoped>
/* 卡头只放操作按钮（页面标题由顶栏展示） */
.page-actions {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
/* 表格控制行：发丝线与表格分区，与其他系统页的搜索栏同层 */
.search-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  padding-bottom: 12px;
  margin-bottom: 12px;
  border-bottom: 1px solid var(--sx-line);
}
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
  border-color: var(--sx-accent, #3F75AB);
}
.icon-cell.active {
  border-color: var(--sx-accent, #3F75AB);
  background: var(--sx-accent-soft, #E4EDF6);
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
