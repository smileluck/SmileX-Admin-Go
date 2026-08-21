<template>
  <n-card title="角色管理">
    <template #header-extra>
      <n-space>
        <n-input v-model:value="query.name" placeholder="角色名" clearable style="width: 160px" @keyup.enter="load" />
        <n-button type="primary" @click="load">查询</n-button>
        <n-button type="primary" ghost @click="openCreate" v-permission="['role:create']">新增角色</n-button>
      </n-space>
    </template>
    <n-data-table :columns="columns" :data="rows" :loading="loading" :pagination="pagination" remote />
  </n-card>

  <!-- 新增/编辑 -->
  <n-modal v-model:show="showModal" preset="dialog" :title="editing ? '编辑角色' : '新增角色'" style="width: 460px">
    <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80">
      <n-form-item label="名称" path="name"><n-input v-model:value="form.name" /></n-form-item>
      <n-form-item label="编码" path="code" v-if="!editing">
        <n-input v-model:value="form.code" placeholder="唯一，如 editor" />
      </n-form-item>
      <n-form-item label="备注"><n-input v-model:value="form.remark" /></n-form-item>
    </n-form>
    <template #action>
      <n-button @click="showModal = false">取消</n-button>
      <n-button type="primary" :loading="saving" @click="save">确定</n-button>
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
import { computed, h, onMounted, reactive, ref, type VNode } from 'vue'
import { NCard, NSpace, NInput, NButton, NDataTable, NModal, NForm, NFormItem, NTree, NTag, useMessage, useDialog, type DataTableColumns, type FormInst, type FormRules } from 'naive-ui'
import { createRole, deleteRole, getRole, listPermissions, listRoles, setRolePermissions, updateRole } from '../../api'
import { useUserStore } from '../../stores/user'
import type { Permission, Role } from '../../api/types'

const message = useMessage()
const dialog = useDialog()
const userStore = useUserStore()
const loading = ref(false)
const saving = ref(false)
const rows = ref<Role[]>([])
const query = reactive({ name: '', page: 1, page_size: 10 })

// 超管角色（id=1）禁止删除
const SUPER_ROLE_ID = 1

const showModal = ref(false)
const showPerm = ref(false)
const editing = ref(false)
const editId = ref(0)
const form = reactive({ name: '', code: '', remark: '' })
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  name: [{ required: true, message: '请输入角色名称', trigger: ['blur', 'input'] }],
  code: [{ required: true, message: '请输入角色编码', trigger: ['blur', 'input'] }],
}

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
  try {
    await formRef.value?.validate()
  } catch {
    return // 校验失败，错误已在表单项上展示
  }
  saving.value = true
  try {
    if (editing.value) {
      await updateRole(editId.value, { name: form.name.trim(), remark: form.remark })
    } else {
      await createRole({ name: form.name.trim(), code: form.code.trim(), remark: form.remark })
    }
    message.success('保存成功')
    showModal.value = false
    load()
  } catch (e: any) {
    message.error(e?.response?.data?.msg || '保存失败')
  } finally {
    saving.value = false
  }
}

function confirmDelete(row: Role) {
  if (row.id === SUPER_ROLE_ID) { message.error('超级管理员角色禁止删除'); return }
  dialog.warning({
    title: '删除确认',
    content: `确定删除角色「${row.name}」吗？该操作不可恢复。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteRole(row.id)
        message.success('已删除')
        load()
      } catch (e: any) {
        message.error(e?.response?.data?.msg || '删除失败')
      }
    },
  })
}

// 构造权限树：菜单 + 按钮权限点统一按 parent_id 组树（按钮自然挂在所属菜单下）
async function openPerms(row: Role) {
  editId.value = row.id
  try {
    const [{ data: allResp }, { data: roleResp }] = await Promise.all([
      listPermissions({ page: 1, page_size: 500 }),
      getRole(row.id),
    ])
    const all = allResp.data.list
    permTree.value = buildTree(all)
    expandedKeys.value = all.filter((p) => p.parent_id === 0).map((p) => p.id)
    checkedKeys.value = roleResp.data.permission_ids ?? []
    showPerm.value = true
  } catch (e: any) {
    message.error(e?.response?.data?.msg || '加载权限数据失败')
  }
}

function buildTree(items: Permission[], parentID = 0): any[] {
  return items
    .filter((p) => p.parent_id === parentID)
    .sort((a, b) => a.sort - b.sort)
    .map((p) => {
      const children = buildTree(items, p.id)
      const n: any = { key: p.id, label: `${p.name}（${p.code}）${p.type === 'button' ? '［按钮］' : ''}` }
      if (children.length) n.children = children
      return n
    })
}

async function savePerms() {
  try {
    await setRolePermissions(editId.value, checkedKeys.value as number[])
    message.success('权限已更新')
    showPerm.value = false
  } catch (e: any) {
    message.error(e?.response?.data?.msg || '保存失败')
  }
}

// 操作列依赖按钮权限，computed 使权限变化后重新渲染
const columns = computed<DataTableColumns<Role>>(() => [
  { title: 'ID', key: 'id', width: 60 },
  { title: '名称', key: 'name' },
  { title: '编码', key: 'code' },
  { title: '备注', key: 'remark' },
  {
    title: '操作', key: 'actions', width: 220,
    render(row) {
      const actions: VNode[] = []
      if (userStore.has('role:update')) {
        actions.push(h(NButton, { size: 'small', onClick: () => openEdit(row) }, { default: () => '编辑' }))
      }
      if (userStore.has('role:setPermissions')) {
        actions.push(h(NButton, { size: 'small', type: 'info', onClick: () => openPerms(row) }, { default: () => '分配权限' }))
      }
      // 超管角色禁止删除，不展示删除按钮
      if (row.id !== SUPER_ROLE_ID && userStore.has('role:delete')) {
        actions.push(h(NButton, { size: 'small', type: 'error', onClick: () => confirmDelete(row) }, { default: () => '删除' }))
      } else if (row.id === SUPER_ROLE_ID) {
        actions.push(h(NTag, { size: 'small', bordered: false }, { default: () => '内置' }))
      }
      if (actions.length === 0) {
        return h('span', { style: 'color: #999; font-size: 12px' }, '—')
      }
      return h(NSpace, {}, { default: () => actions })
    },
  },
])

onMounted(load)
</script>
