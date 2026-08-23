<template>
  <!-- 搜索栏独立卡片：可折叠，重置/搜索按钮在卡片右下角 -->
  <SearchCard storage-key="users" @search="load" @reset="resetQuery">
    <n-input v-model:value="query.username" placeholder="用户名" clearable style="width: 180px" @keyup.enter="load" />
  </SearchCard>

  <n-card>
    <template #header>
      <div class="page-actions">
        <n-button ghost :loading="exporting" @click="doExport" v-permission="['user:export']">导出</n-button>
        <n-button type="primary" ghost @click="openCreate" v-permission="['user:create']">新增用户</n-button>
      </div>
    </template>

    <n-data-table :columns="columns" :data="rows" :loading="loading" :pagination="pagination" remote />
  </n-card>

  <!-- 新增/编辑 -->
  <n-modal v-model:show="showModal" preset="dialog" :title="editing ? '编辑用户' : '新增用户'" style="width: 480px">
    <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80">
      <n-form-item label="用户名" path="username" v-if="!editing">
        <n-input v-model:value="form.username" placeholder="3-64 个字符" />
      </n-form-item>
      <n-form-item label="密码" path="password" v-if="!editing">
        <n-input v-model:value="form.password" type="password" :maxlength="20" placeholder="6-20 位" />
      </n-form-item>
      <n-form-item label="昵称" path="nickname">
        <n-input v-model:value="form.nickname" :maxlength="20" show-word-limit placeholder="最多 20 个字符" />
      </n-form-item>
      <n-form-item label="邮箱" path="email"><n-input v-model:value="form.email" placeholder="选填" /></n-form-item>
      <n-form-item label="角色">
        <n-select v-model:value="form.role_ids" multiple :options="roleOptions" :disabled="editing && !userStore.has('user:setRoles')" />
      </n-form-item>
      <n-form-item label="状态" v-if="editing">
        <n-switch v-model:value="form.statusOn" :checked-value="1" :unchecked-value="0" />
      </n-form-item>
    </n-form>
    <template #action>
      <n-button @click="showModal = false">取消</n-button>
      <n-button type="primary" :loading="saving" @click="save">确定</n-button>
    </template>
  </n-modal>

  <!-- 重置密码 -->
  <n-modal v-model:show="showPwd" preset="dialog" title="重置密码" style="width: 420px">
    <n-input v-model:value="newPassword" type="password" :maxlength="20" placeholder="新密码（6-20 位）" />
    <template #action>
      <n-button @click="showPwd = false">取消</n-button>
      <n-button type="primary" @click="savePwd">确定</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, h, onMounted, reactive, ref, type VNode } from 'vue'
import { NCard, NInput, NButton, NDataTable, NModal, NForm, NFormItem, NSelect, NSwitch, NTag, useMessage, useDialog, type DataTableColumns, type FormInst, type FormRules } from 'naive-ui'
import { renderActions, type TableAction } from '../../utils/tableActions'
import SearchCard from '../../components/SearchCard.vue'
import { createUser, createExport, deleteUser, listRoles, listUsers, resetPassword, setUserRoles, updateUser } from '../../api'
import { useUserStore } from '../../stores/user'
import type { UserInfo } from '../../api/types'

const message = useMessage()
const dialog = useDialog()
const userStore = useUserStore()

// admin（id=1）超管账号：仅其本人可见操作按钮
const SUPER_ADMIN_ID = 1
const canOperate = (row: UserInfo) => row.id !== SUPER_ADMIN_ID || userStore.user?.id === SUPER_ADMIN_ID
const loading = ref(false)
const saving = ref(false)
const rows = ref<UserInfo[]>([])
const total = ref(0)
const query = reactive({ username: '', page: 1, page_size: 10 })
const roleOptions = ref<{ label: string; value: number }[]>([])

const showModal = ref(false)
const showPwd = ref(false)
const editing = ref(false)
const editId = ref(0)
const newPassword = ref('')
const form = reactive({ username: '', password: '', nickname: '', email: '', role_ids: [] as number[], statusOn: 1 })
const formRef = ref<FormInst | null>(null)

// 与后端 binding 规则保持一致：用户名 3-64、密码 6-20、昵称最长 20、邮箱格式（选填）
const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: ['blur', 'input'] },
    { min: 3, max: 64, message: '用户名长度 3-64 个字符', trigger: ['blur', 'input'] },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: ['blur', 'input'] },
    { min: 6, max: 20, message: '密码长度 6-20 个字符', trigger: ['blur', 'input'] },
  ],
  nickname: [
    { max: 20, message: '昵称不能超过 20 个字符', trigger: ['blur', 'input'] },
  ],
  email: [
    {
      trigger: ['blur', 'input'],
      validator: (_rule, value: string) => !value || /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value),
      message: '邮箱格式不正确',
    },
  ],
}

const pagination = reactive({
  page: 1, pageSize: 10, pageCount: 1, showSizePicker: true,
  onChange: (p: number) => { query.page = p; load() },
  onUpdatePageSize: (s: number) => { query.page_size = s; load() },
})

async function load() {
  loading.value = true
  try {
    const { data } = await listUsers(query)
    rows.value = data.data.list
    total.value = data.data.page.total
    pagination.page = query.page
    pagination.pageSize = query.page_size
    pagination.pageCount = Math.ceil(total.value / query.page_size)
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.username = ''
  query.page = 1
  load()
}

// 异步导出：提交当前过滤条件（page/page_size 与空值由 createExport 剔除）
const exporting = ref(false)
async function doExport() {
  exporting.value = true
  try {
    await createExport('users', { ...query })
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

async function loadRoles() {
  const { data } = await listRoles({ page: 1, page_size: 100 })
  roleOptions.value = data.data.list.map((r) => ({ label: r.name, value: r.id }))
}

function openCreate() {
  editing.value = false
  Object.assign(form, { username: '', password: '', nickname: '', email: '', role_ids: [], statusOn: 1 })
  showModal.value = true
}

function openEdit(row: UserInfo) {
  editing.value = true
  editId.value = row.id
  Object.assign(form, { username: row.username, password: '', nickname: row.nickname, email: row.email, role_ids: row.role_ids ?? [], statusOn: row.status })
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
      await updateUser(editId.value, { nickname: form.nickname, email: form.email, status: form.statusOn })
      // 分配角色是独立接口权限，未授权时不发起该调用（角色保持不变）
      if (userStore.has('user:setRoles')) {
        await setUserRoles(editId.value, form.role_ids)
      }
    } else {
      await createUser({ username: form.username.trim(), password: form.password, nickname: form.nickname.trim(), email: form.email.trim(), role_ids: form.role_ids })
    }
    message.success('保存成功')
    showModal.value = false
    await load()
  } catch (e: any) {
    message.error(e?.response?.data?.msg || '保存失败')
  } finally {
    saving.value = false
  }
}

function confirmDelete(row: UserInfo) {
  if (row.id === SUPER_ADMIN_ID) { message.error('超级管理员账号禁止删除'); return }
  dialog.warning({
    title: '删除确认',
    content: `确定删除用户「${row.username}」吗？该操作不可恢复。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteUser(row.id)
        message.success('已删除')
        load()
      } catch (e: any) {
        message.error(e?.response?.data?.msg || '删除失败')
      }
    },
  })
}

async function savePwd() {
  if (newPassword.value.length < 6 || newPassword.value.length > 20) { message.warning('密码长度 6-20 位'); return }
  try {
    await resetPassword(editId.value, newPassword.value)
    message.success('密码已重置')
    showPwd.value = false
  } catch (e: any) {
    message.error(e?.response?.data?.msg || '重置失败')
  }
}

// 操作列依赖按钮权限，computed 使权限变化后重新渲染
const columns = computed<DataTableColumns<UserInfo>>(() => [
  { title: 'ID', key: 'id', width: 60 },
  { title: '用户名', key: 'username' },
  { title: '昵称', key: 'nickname' },
  { title: '邮箱', key: 'email' },
  {
    title: '状态', key: 'status', width: 80,
    render: (row) => h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' }, { default: () => (row.status === 1 ? '启用' : '禁用') }),
  },
  { title: '创建时间', key: 'created_at', width: 170 },
  {
    title: '操作', key: 'actions', width: 160,
    render(row) {
      if (!canOperate(row)) {
        return renderActions([])
      }
      const actions: Array<TableAction | VNode> = []
      if (userStore.has('user:update')) {
        actions.push({ label: '编辑', accent: true, onClick: () => openEdit(row) })
      }
      if (userStore.has('user:resetPassword')) {
        actions.push({ label: '重置密码', onClick: () => { editId.value = row.id; newPassword.value = ''; showPwd.value = true } })
      }
      // 超管账号禁止删除（即使超管本人），不展示删除按钮
      if (row.id !== SUPER_ADMIN_ID && userStore.has('user:delete')) {
        actions.push({ label: '删除', danger: true, onClick: () => confirmDelete(row) })
      }
      return renderActions(actions)
    },
  },
])

onMounted(() => { load(); loadRoles() })
</script>

<style scoped>
/* 卡头只放操作按钮（页面标题由顶栏展示） */
.page-actions {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
