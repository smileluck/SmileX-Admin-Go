<template>
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider class="sider" collapse-mode="width" :collapsed-width="64" :width="200" :collapsed="collapsed">
      <div class="logo" :class="{ 'logo-collapsed': collapsed }">
        <div class="seal">S</div>
        <div v-if="!collapsed" class="logo-text">
          <span class="logo-name">SmileX</span>
          <span class="logo-sub mono">admin console</span>
        </div>
      </div>
      <n-menu class="sider-menu" :collapsed="collapsed" :collapsed-width="64" :root-indent="16" :indent="20"
        :options="menuOptions" :value="activeKey" @update:value="onMenuSelect" />
      <div class="sider-foot mono">{{ collapsed ? 'v1.0' : 'v1.0 · internal' }}</div>
    </n-layout-sider>

    <n-layout class="main">
      <n-layout-header bordered class="header">
        <div class="header-left">
          <n-button class="sider-trigger" quaternary circle :focusable="false" aria-label="折叠/展开侧边栏"
            @click="toggleCollapsed">
            <template #icon>
              <n-icon :component="MenuOutline" />
            </template>
          </n-button>
          <div class="crumb">
            <span class="crumb-eyebrow mono">section</span>
            <span class="crumb-title">{{ route.meta?.title || '首页' }}</span>
          </div>
        </div>
        <div class="header-right">
          <n-button
            class="search-trigger"
            quaternary
            circle
            :focusable="false"
            :aria-label="`搜索菜单页面 ${searchKbd}`"
            @click="openSearch"
          >
            <template #icon>
              <n-icon :component="SearchOutline" />
            </template>
          </n-button>
          <n-dropdown :options="userOptions" @select="onUserAction">
            <div class="user-chip">
              <div class="avatar">{{ avatarChar }}</div>
              <span class="user-name">{{ userStore.user?.nickname || userStore.user?.username }}</span>
            </div>
          </n-dropdown>
        </div>
      </n-layout-header>
      <n-layout-content class="content" content-style="padding: 8px;" :native-scrollbar="false">
        <router-view />
      </n-layout-content>
    </n-layout>

    <!-- 修改密码（右上角用户下拉触发；成功后强制重新登录） -->
    <n-modal v-model:show="showPwd" preset="dialog" title="修改密码" style="width: 420px">
      <n-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-placement="left" label-width="90">
        <n-form-item label="原密码" path="oldPassword">
          <n-input v-model:value="pwdForm.oldPassword" type="password" show-password-on="click" placeholder="请输入原密码" />
        </n-form-item>
        <n-form-item label="新密码" path="newPassword">
          <n-input v-model:value="pwdForm.newPassword" type="password" show-password-on="click" placeholder="6-64 位" />
        </n-form-item>
        <n-form-item label="确认新密码" path="confirmPassword">
          <n-input v-model:value="pwdForm.confirmPassword" type="password" show-password-on="click" placeholder="再次输入新密码" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showPwd = false">取消</n-button>
        <n-button type="primary" :loading="pwdSaving" @click="savePwd">确定</n-button>
      </template>
    </n-modal>

    <!-- 顶栏搜索：点击搜索图标弹出，选中后跳转对应页面 -->
    <n-modal v-model:show="showSearch" preset="card" :bordered="false" style="width: 520px">
      <n-auto-complete
        ref="searchRef"
        v-model:value="searchKw"
        class="search-palette"
        size="large"
        :options="searchOptions"
        :render-label="renderSearchLabel"
        :input-props="{ autocomplete: 'off' }"
        :placeholder="`搜索菜单页面 ${searchKbd}`"
        clearable
        @select="onSearchSelect"
      />
    </n-modal>
  </n-layout>
</template>

<script setup lang="ts">
import { computed, h, nextTick, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NDropdown, NButton, NIcon,
  NAutoComplete, NModal, NForm, NFormItem, NInput, useMessage,
  type DropdownOption, type FormInst, type FormRules,
} from 'naive-ui'
import { MenuOutline, SearchOutline } from '@vicons/ionicons5'
import { useUserStore } from '../stores/user'
import { renderMenuIcon } from '../utils/menuIcon'
import { changePassword } from '../api'
import type { MenuNode } from '../api/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const message = useMessage()

// 折叠状态持久化到 localStorage，刷新后保持
const COLLAPSED_KEY = 'sider_collapsed'
const collapsed = ref(localStorage.getItem(COLLAPSED_KEY) === '1')

function toggleCollapsed() {
  collapsed.value = !collapsed.value
  localStorage.setItem(COLLAPSED_KEY, collapsed.value ? '1' : '0')
}

// 图标支持本地 ionicons5 名称 / 网络图片 URL，统一走 menuIcon 渲染
const renderIcon = (icon?: string) => () => h('span', { style: 'display:inline-flex;align-items:center' }, renderMenuIcon(icon))

function toOptions(nodes: MenuNode[]): any[] {
  return (nodes ?? []).map((m) => ({
    label: m.name,
    key: m.path,
    icon: renderIcon(m.icon),
    children: m.children?.length ? toOptions(m.children) : undefined,
  }))
}

const menuOptions = computed(() => toOptions(userStore.menus))
const activeKey = computed(() => route.path)

// 侧边栏点击跳转：key 即菜单 path（叶子节点对应已注册的动态路由）
function onMenuSelect(key: string) {
  if (key && key !== route.path) {
    router.push(key)
  }
}

const avatarChar = computed(() => (userStore.user?.nickname || userStore.user?.username || 'U').charAt(0).toUpperCase())

const userOptions: DropdownOption[] = [
  { label: '个人中心', key: 'profile' },
  { label: '修改密码', key: 'password' },
  { type: 'divider', key: 'divider' },
  { label: '退出登录', key: 'logout' },
]

async function onUserAction(key: string | number) {
  if (key === 'logout') {
    await userStore.logout()
    router.push('/login')
  } else if (key === 'profile') {
    router.push('/profile')
  } else if (key === 'password') {
    openPwdModal()
  }
}

// ---- 修改密码（本人，校验原密码；成功后强制重新登录） ----
const showPwd = ref(false)
const pwdSaving = ref(false)
const pwdFormRef = ref<FormInst | null>(null)
const pwdForm = reactive({ oldPassword: '', newPassword: '', confirmPassword: '' })
// 与后端 binding 保持一致：6-64 位，两次输入须一致
const pwdRules: FormRules = {
  oldPassword: [{ required: true, message: '请输入原密码', trigger: ['blur', 'input'] }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: ['blur', 'input'] },
    { min: 6, max: 64, message: '密码长度 6-64 位', trigger: ['blur', 'input'] },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: ['blur', 'input'] },
    { validator: (_rule, v: string) => v === pwdForm.newPassword, message: '两次输入的密码不一致', trigger: ['blur', 'input'] },
  ],
}

function openPwdModal() {
  Object.assign(pwdForm, { oldPassword: '', newPassword: '', confirmPassword: '' })
  showPwd.value = true
}

async function savePwd() {
  try {
    await pwdFormRef.value?.validate()
  } catch {
    return // 校验失败，错误已在表单项上展示
  }
  pwdSaving.value = true
  try {
    await changePassword({ old_password: pwdForm.oldPassword, new_password: pwdForm.newPassword })
    showPwd.value = false
    message.success('密码已修改，请重新登录')
    await userStore.logout()
    router.push('/login')
  } catch (e: any) {
    message.error(e?.response?.data?.msg || '修改失败')
  } finally {
    pwdSaving.value = false
  }
}

// ---- 顶栏全局搜索：按菜单名称/路径模糊匹配并快速跳转 ----
interface SearchItem { name: string; path: string; icon?: string; parents: string }

// 父级菜单无组件，命中时跳到其第一个叶子节点对应的路由
function firstLeafPath(n: MenuNode): string {
  return n.children?.length ? firstLeafPath(n.children[0]) : n.path
}

function flattenMenus(nodes: MenuNode[] | null, parent: string): SearchItem[] {
  const out: SearchItem[] = []
  for (const n of nodes ?? []) {
    out.push({ name: n.name, path: firstLeafPath(n), icon: n.icon, parents: parent })
    if (n.children?.length) {
      out.push(...flattenMenus(n.children, parent ? `${parent} / ${n.name}` : n.name))
    }
  }
  return out
}

const searchItems = computed(() => flattenMenus(userStore.menus, ''))
const searchKw = ref('')
const searchRef = ref<{ focus: () => void } | null>(null)
// 快捷键提示按平台显示（Mac 显示 ⌘K）
const searchKbd = /Mac|iPhone|iPad/.test(navigator.platform || navigator.userAgent) ? '⌘K' : 'Ctrl K'

// ---- 搜索弹窗（点击顶栏图标 / ⌘K 唤起） ----
const showSearch = ref(false)

function openSearch() {
  searchKw.value = ''
  showSearch.value = true
}

// 弹窗打开后自动聚焦输入框
watch(showSearch, async (v) => {
  if (v) {
    await nextTick()
    searchKw.value = ''
    searchRef.value?.focus()
  }
})

// 命中项（value 用数组下标，渲染与跳转都从 searchMatches 取）
const searchMatches = computed(() => {
  const kw = searchKw.value.trim().toLowerCase()
  if (!kw) return []
  return searchItems.value
    .filter((it) => it.name.toLowerCase().includes(kw) || it.path.toLowerCase().includes(kw))
    .slice(0, 8)
})
const searchOptions = computed(() => searchMatches.value.map((it, i) => ({ label: it.name, value: String(i) })))

// 自定义渲染下拉项：图标 + 名称 + 父级路径（下拉在 teleport 层，用内联样式避免 scoped CSS 失效）
function renderSearchLabel(option: { value?: string | number; label?: string | number }) {
  const it = searchMatches.value[Number(option.value)]
  if (!it) return String(option.label ?? '')
  return h('div', { style: 'display:flex;align-items:center;gap:8px;min-width:230px;padding:2px 0' }, [
    h('span', { style: 'display:inline-flex;align-items:center;color:var(--sx-muted)' }, renderMenuIcon(it.icon, 15)),
    h('span', { style: 'color:var(--sx-ink)' }, it.name),
    it.parents ? h('span', { style: 'margin-left:auto;font-size:11px;color:var(--sx-muted);white-space:nowrap' }, it.parents) : null,
  ])
}

function onSearchSelect(value: string | number) {
  const it = searchMatches.value[Number(value)]
  // naive-ui 选中后会把 label 写回输入框，nextTick 后再清空才能生效
  nextTick(() => { searchKw.value = '' })
  showSearch.value = false
  if (it && it.path !== route.path) {
    router.push(it.path)
  }
}

// Ctrl/Cmd + K 唤起搜索弹窗
function onGlobalKeydown(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    openSearch()
  }
}

onMounted(() => window.addEventListener('keydown', onGlobalKeydown))
onUnmounted(() => window.removeEventListener('keydown', onGlobalKeydown))
</script>

<style scoped>
/* 侧边栏：浅色线框式，安静克制 */
.sider {
  background: var(--sx-surface) !important;
}
/* flex 作用到 naive 原生滚动容器，保证 logo/菜单/脚注三段式撑满整列 */
.sider :deep(.n-layout-sider-scroll-container) {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden auto;
  scrollbar-width: thin;
}
.logo {
  height: 64px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
  border-bottom: 1px solid var(--sx-line);
  margin-bottom: 4px;
}
.seal {
  width: 34px;
  height: 34px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 9px;
  font-family: var(--sx-font-mono);
  font-weight: 700;
  font-size: 17px;
  color: #fff;
  background: var(--sx-accent);
}
.logo-text {
  display: flex;
  flex-direction: column;
  line-height: 1.25;
  white-space: nowrap;
}
.logo-name {
  font-weight: 700;
  font-size: 15px;
  color: var(--sx-ink);
}
.logo-sub {
  font-size: 10px;
}

.sider-menu {
  flex: 1;
  background: transparent !important;
  padding: 4px 8px;
}
/* 收起态：清零水平内边距。naive-ui 按 collapsedWidth/2 - 图标宽/2 计算缩进居中图标，
   前提是菜单占满整个窄栏；保留 8px 内边距会让图标整体偏右 8px */
.sider-menu.n-menu--collapsed {
  padding-left: 0;
  padding-right: 0;
}
.sider-menu :deep(.n-menu .n-menu-item-content::before) {
  left: 8px;
  right: 8px;
  border-radius: 7px;
}
/* 收起态 hover 背景占满整行（内边距已清零） */
.sider-menu.n-menu--collapsed :deep(.n-menu-item-content::before) {
  left: 0;
  right: 0;
  border-radius: 0;
}
/* 收起态：仅显示印章，在 64px 窄栏内居中 */
.logo.logo-collapsed {
  padding: 0;
  justify-content: center;
}

.sider-foot {
  flex-shrink: 0;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-top: 1px solid var(--sx-line);
  font-size: 10px;
  white-space: nowrap;
  overflow: hidden;
}

/* 顶部栏：eyebrow + 标题的层级读法 */
.main {
  background: var(--sx-bg) !important;
}
.header {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: var(--sx-surface);
  box-shadow: var(--sx-shadow);
  position: relative;
  z-index: 1;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}
.sider-trigger {
  flex-shrink: 0;
}
.crumb {
  display: flex;
  flex-direction: column;
  line-height: 1.3;
}
.crumb-eyebrow {
  font-size: 10px;
}
.crumb-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--sx-ink);
}

/* 顶栏右侧：搜索图标 + 用户区 */
.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}
.search-trigger {
  flex-shrink: 0;
  color: var(--sx-muted);
}

.user-chip {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 5px 14px 5px 6px;
  border-radius: 999px;
  border: 1px solid var(--sx-line);
  background: var(--sx-surface);
  cursor: pointer;
  transition: border-color 0.2s ease, background 0.2s ease;
}
.user-chip:hover {
  border-color: var(--sx-accent);
  background: var(--sx-accent-soft);
}
.avatar {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #fff;
  font-weight: 700;
  font-size: 13px;
  background: var(--sx-accent);
}
.user-name {
  font-size: 13px;
  color: var(--sx-ink);
}

/* 内容区 */
.content {
  background: transparent !important;
}

/* 搜索弹窗：输入框撑满弹窗宽度 */
.search-palette {
  width: 100%;
}
</style>
