<template>
  <!-- 搜索栏独立卡片：可折叠，重置/搜索按钮在卡片右下角 -->
  <SearchCard storage-key="roles" @search="load" @reset="resetQuery">
    <n-input v-model:value="query.name" placeholder="角色名" clearable style="width: 180px" @keyup.enter="load" />
  </SearchCard>

  <n-card>
    <template #header>
      <div class="page-actions">
        <n-button type="primary" ghost @click="openCreate" v-permission="['role:create']">新增角色</n-button>
      </div>
    </template>

    <n-data-table :columns="columns" :data="rows" :loading="loading" :pagination="pagination" remote />
  </n-card>

  <!-- 新增/编辑 -->
  <n-modal v-model:show="showModal" preset="dialog" :title="editing ? '编辑角色' : '新增角色'" style="width: 460px">
    <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80">
      <n-form-item label="名称" path="name">
        <n-input v-model:value="form.name" :maxlength="20" show-word-limit placeholder="最多 20 个字符" />
      </n-form-item>
      <n-form-item label="备注">
        <n-input v-model:value="form.remark" :maxlength="200" show-word-limit placeholder="最多 200 个字符" />
      </n-form-item>
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
import { NCard, NInput, NButton, NDataTable, NModal, NForm, NFormItem, NTree, NTag, useMessage, useDialog, type DataTableColumns, type FormInst, type FormRules } from 'naive-ui'
import { renderActions, type TableAction } from '../../utils/tableActions'
import SearchCard from '../../components/SearchCard.vue'
import { createRole, deleteRole, getRole, listAllPermissions, listRoles, setRolePermissions, updateRole } from '../../api'
import { useUserStore } from '../../stores/user'
import type { Permission, Role } from '../../api/types'

const message = useMessage()
const dialog = useDialog()
const userStore = useUserStore()
const loading = ref(false)
const saving = ref(false)
const rows = ref<Role[]>([])
const query = reactive({ name: '', page: 1, page_size: 10 })

// 超管角色（id=1）系统内置：禁止修改和操作
const SUPER_ROLE_ID = 1
const SUPER_ROLE_MSG = '超级管理员角色为系统内置，禁止修改和操作'

const showModal = ref(false)
const showPerm = ref(false)
const editing = ref(false)
const editId = ref(0)
const form = reactive({ name: '', remark: '' })
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: ['blur', 'input'] },
    { max: 20, message: '角色名称不能超过 20 个字符', trigger: ['blur', 'input'] },
  ],
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

function resetQuery() {
  query.name = ''
  query.page = 1
  load()
}

function openCreate() {
  editing.value = false
  Object.assign(form, { name: '', remark: '' })
  showModal.value = true
}

function openEdit(row: Role) {
  if (row.id === SUPER_ROLE_ID) { message.error(SUPER_ROLE_MSG); return }
  editing.value = true
  editId.value = row.id
  Object.assign(form, { name: row.name, remark: row.remark })
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
      await updateRole(editId.value, { name: form.name.trim(), remark: form.remark.trim() })
    } else {
      await createRole({ name: form.name.trim(), remark: form.remark.trim() })
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
  if (row.id === SUPER_ROLE_ID) { message.error(SUPER_ROLE_MSG); return }
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
  if (row.id === SUPER_ROLE_ID) { message.error(SUPER_ROLE_MSG); return }
  editId.value = row.id
  try {
    // 权限树需整表构建，走全量接口（page_size=0），分页会截断子节点
    const [{ data: allResp }, { data: roleResp }] = await Promise.all([
      listAllPermissions(),
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
  { title: '备注', key: 'remark' },
  {
    title: '操作', key: 'actions', width: 170,
    render(row) {
      // 超管角色系统内置：禁止修改和操作，仅展示「内置」标记
      if (row.id === SUPER_ROLE_ID) {
        return h(NTag, { size: 'small', bordered: false }, { default: () => '内置' })
      }
      const actions: Array<TableAction | VNode> = []
      if (userStore.has('role:update')) {
        actions.push({ label: '编辑', accent: true, onClick: () => openEdit(row) })
      }
      if (userStore.has('role:setPermissions')) {
        actions.push({ label: '分配权限', onClick: () => openPerms(row) })
      }
      if (userStore.has('role:delete')) {
        actions.push({ label: '删除', danger: true, onClick: () => confirmDelete(row) })
      }
      return renderActions(actions)
    },
  },
])

onMounted(load)
</script>

<style scoped>
/* 卡头只放操作按钮（页面标题由顶栏展示） */
.page-actions {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}
</style>
