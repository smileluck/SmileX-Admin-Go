<template>
  <n-card title="用户管理">
    <template #header-extra>
      <n-space>
        <n-input v-model:value="query.username" placeholder="用户名" clearable style="width: 160px" @keyup.enter="load" />
        <n-button type="primary" @click="load">查询</n-button>
        <n-button type="primary" ghost @click="openCreate" v-permission="['menu:user']">新增用户</n-button>
      </n-space>
    </template>

    <n-data-table :columns="columns" :data="rows" :loading="loading" :pagination="pagination" remote />
  </n-card>

  <!-- 新增/编辑 -->
  <n-modal v-model:show="showModal" preset="dialog" :title="editing ? '编辑用户' : '新增用户'" style="width: 480px">
    <n-form :model="form" label-placement="left" label-width="80">
      <n-form-item label="用户名" v-if="!editing">
        <n-input v-model:value="form.username" />
      </n-form-item>
      <n-form-item label="密码" v-if="!editing">
        <n-input v-model:value="form.password" type="password" />
      </n-form-item>
      <n-form-item label="昵称"><n-input v-model:value="form.nickname" /></n-form-item>
      <n-form-item label="邮箱"><n-input v-model:value="form.email" /></n-form-item>
      <n-form-item label="角色">
        <n-select v-model:value="form.role_ids" multiple :options="roleOptions" />
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
    <n-input v-model:value="newPassword" type="password" placeholder="新密码（至少 6 位）" />
    <template #action>
      <n-button @click="showPwd = false">取消</n-button>
      <n-button type="primary" @click="savePwd">确定</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { h, onMounted, reactive, ref } from 'vue'
import { NCard, NSpace, NInput, NButton, NDataTable, NModal, NForm, NFormItem, NSelect, NSwitch, NTag, useMessage, useDialog, type DataTableColumns } from 'naive-ui'
import { createUser, deleteUser, listRoles, listUsers, resetPassword, setUserRoles, updateUser } from '../../api'
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
  saving.value = true
  try {
    if (editing.value) {
      await updateUser(editId.value, { nickname: form.nickname, email: form.email, status: form.statusOn })
      await setUserRoles(editId.value, form.role_ids)
    } else {
      await createUser({ username: form.username, password: form.password, nickname: form.nickname, email: form.email, role_ids: form.role_ids })
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
  if (newPassword.value.length < 6) { message.warning('密码至少 6 位'); return }
  await resetPassword(editId.value, newPassword.value)
  message.success('密码已重置')
  showPwd.value = false
}

const columns: DataTableColumns<UserInfo> = [
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
    title: '操作', key: 'actions', width: 220,
    render(row) {
      if (!canOperate(row)) {
        return h('span', { style: 'color: #999; font-size: 12px' }, '—')
      }
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { size: 'small', onClick: () => openEdit(row) }, { default: () => '编辑' }),
          h(NButton, { size: 'small', type: 'warning', onClick: () => { editId.value = row.id; newPassword.value = ''; showPwd.value = true } }, { default: () => '重置密码' }),
          h(NButton, {
            size: 'small', type: 'error',
            onClick: () => confirmDelete(row),
          }, { default: () => '删除' }),
        ],
      })
    },
  },
]

onMounted(() => { load(); loadRoles() })
</script>
