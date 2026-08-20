<template>
  <n-card title="角色管理">
    <template #header-extra>
      <n-space>
        <n-input v-model:value="query.name" placeholder="角色名" clearable style="width: 160px" @keyup.enter="load" />
        <n-button type="primary" @click="load">查询</n-button>
        <n-button type="primary" ghost @click="openCreate">新增角色</n-button>
      </n-space>
    </template>
    <n-data-table :columns="columns" :data="rows" :loading="loading" :pagination="pagination" remote />
  </n-card>

  <!-- 新增/编辑 -->
  <n-modal v-model:show="showModal" preset="dialog" :title="editing ? '编辑角色' : '新增角色'" style="width: 460px">
    <n-form :model="form" label-placement="left" label-width="80">
      <n-form-item label="名称"><n-input v-model:value="form.name" /></n-form-item>
      <n-form-item label="编码" v-if="!editing"><n-input v-model:value="form.code" /></n-form-item>
      <n-form-item label="备注"><n-input v-model:value="form.remark" /></n-form-item>
    </n-form>
    <template #action>
      <n-button @click="showModal = false">取消</n-button>
      <n-button type="primary" @click="save">确定</n-button>
    </template>
  </n-modal>

  <!-- 分配权限 -->
  <n-modal v-model:show="showPerm" preset="dialog" title="分配权限" style="width: 520px">
    <n-tree
      :data="permTree"
      checkable cascade :checked-keys="checkedKeys" :default-expanded-keys="expandedKeys"
      key-field="key" label-field="label" children-field="children"
      @update:checked-keys="(keys: any) => (checkedKeys = keys)"
    />
    <template #action>
      <n-button @click="showPerm = false">取消</n-button>
      <n-button type="primary" @click="savePerms">确定</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue'
import { NCard, NSpace, NInput, NButton, NDataTable, NModal, NForm, NFormItem, NTree, useMessage, type DataTableColumns } from 'naive-ui'
import { createRole, deleteRole, getRole, listPermissions, listRoles, setRolePermissions, updateRole } from '../../api'
import type { Role } from '../../api/types'

const message = useMessage()
const loading = ref(false)
const rows = ref<Role[]>([])
const query = reactive({ name: '', page: 1, page_size: 10 })

const showModal = ref(false)
const showPerm = ref(false)
const editing = ref(false)
const editId = ref(0)
const form = reactive({ name: '', code: '', remark: '' })

const permTree = ref<any[]>([])
const expandedKeys = ref<number[]>([])
const checkedKeys = ref<number[]>([])

const pagination = reactive({
  page: 1, pageSize: 10, pageCount: 1, showSizePicker: true,
  onChange: (p: number) => { query.page = p; load() },
  onUpdatePageSize: (s: number) => { query.page_size = s; load() },
})

async function load() {
  loading.value = true
  try {
    const { data } = await listRoles(query)
    rows.value = data.data.list
    pagination.page = query.page
    pagination.pageCount = Math.ceil(data.data.page.total / query.page_size)
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = false
  Object.assign(form, { name: '', code: '', remark: '' })
  showModal.value = true
}

function openEdit(row: Role) {
  editing.value = true
  editId.value = row.id
  Object.assign(form, { name: row.name, code: row.code, remark: row.remark })
  showModal.value = true
}

async function save() {
  if (editing.value) {
    await updateRole(editId.value, { name: form.name, remark: form.remark })
  } else {
    await createRole(form)
  }
  message.success('保存成功')
  showModal.value = false
  load()
}

// 构造权限树（API + 菜单平铺到树根）
async function openPerms(row: Role) {
  editId.value = row.id
  const [{ data: allResp }, { data: roleResp }] = await Promise.all([
    listPermissions({ page: 1, page_size: 500 }),
    getRole(row.id),
  ])
  const all = allResp.data.list
  const menus = all.filter((p) => p.type === 'menu')
  const apis = all.filter((p) => p.type === 'api')
  const node = (p: any): any => ({ key: p.id, label: `${p.name}（${p.code}）` })
  permTree.value = [
    ...buildTree(menus),
    { key: -1, label: '接口权限', children: apis.map(node) },
  ]
  expandedKeys.value = [-1, ...menus.filter((m) => m.parent_id === 0).map((m) => m.id)]
  checkedKeys.value = roleResp.data.permission_ids ?? []
  showPerm.value = true
}

function buildTree(items: any[], parentID = 0): any[] {
  return items
    .filter((p) => p.parent_id === parentID)
    .map((p) => {
      const children = buildTree(items, p.id)
      const n: any = { key: p.id, label: `${p.name}（${p.code}）` }
      if (children.length) n.children = children
      return n
    })
}

async function savePerms() {
  const ids = (checkedKeys.value as number[]).filter((k) => k > 0)
  await setRolePermissions(editId.value, ids)
  message.success('权限已更新')
  showPerm.value = false
}

const columns: DataTableColumns<Role> = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '名称', key: 'name' },
  { title: '编码', key: 'code' },
  { title: '备注', key: 'remark' },
  {
    title: '操作', key: 'actions', width: 220,
    render(row) {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { size: 'small', onClick: () => openEdit(row) }, { default: () => '编辑' }),
          h(NButton, { size: 'small', type: 'info', onClick: () => openPerms(row) }, { default: () => '分配权限' }),
          h(NButton, {
            size: 'small', type: 'error',
            onClick: async () => { await deleteRole(row.id); message.success('已删除'); load() },
          }, { default: () => '删除' }),
        ],
      })
    },
  },
]

onMounted(load)
</script>
