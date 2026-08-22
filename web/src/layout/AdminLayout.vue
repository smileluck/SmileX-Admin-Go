<template>
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider class="sider" collapse-mode="width" :collapsed-width="64" :width="200" :collapsed="collapsed">
      <div class="logo" :class="{ 'logo-collapsed': collapsed }">
        <div class="seal">S</div>
        <div class="logo-text">
          <span class="logo-name">SmileX</span>
          <span class="logo-sub mono">admin console</span>
        </div>
      </div>
      <n-menu class="sider-menu" :collapsed="collapsed" :collapsed-width="64" :root-indent="16" :indent="20"
        :options="menuOptions" :value="activeKey" :expanded-keys="expandedKeys"
        @update:expanded-keys="expandedKeys = $event" @update:value="onMenuSelect" />
      <div class="sider-foot mono" :class="{ 'foot-collapsed': collapsed }">v1.0<span class="foot-ext"> · internal</span></div>
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
      <n-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-placement="left" label-width="110">
        <n-form-item label="原密码" path="oldPassword">
          <n-input v-model:value="pwdForm.oldPassword" type="password" show-password-on="click" placeholder="请输入原密码" />
        </n-form-item>
        <n-form-item label="新密码" path="newPassword">
          <n-input v-model:value="pwdForm.newPassword" type="password" show-password-on="click" :maxlength="20" placeholder="6-20 位" />
        </n-form-item>
        <n-form-item label="确认新密码" path="confirmPassword">
          <n-input v-model:value="pwdForm.confirmPassword" type="password" show-password-on="click" :maxlength="20" placeholder="再次输入新密码" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showPwd = false">取消</n-button>
        <n-button type="primary" :loading="pwdSaving" @click="savePwd">确定</n-button>
      </template>
    </n-modal>

    <!-- 顶栏搜索：命令面板（输入行 / 结果列表 / 快捷键提示栏） -->
    <n-modal v-model:show="showSearch" :bordered="false" style="width: auto">
      <div class="cmd-palette">
        <div class="cmd-input-row">
          <n-icon :component="SearchOutline" :size="18" class="cmd-search-icon" />
          <input
            ref="searchRef"
            v-model="searchKw"
            class="cmd-input"
            type="text"
            autocomplete="off"
            placeholder="搜索菜单页面…"
            @keydown="onCmdKeydown"
          />
          <span class="kbd">esc</span>
        </div>

        <div ref="cmdListRef" class="cmd-list">
          <div v-if="searchLoading" class="cmd-empty">搜索中…</div>
          <div v-else-if="!searchKw.trim()" class="cmd-empty">输入关键词搜索菜单</div>
          <div v-else-if="!menuCount" class="cmd-empty">未找到匹配的菜单</div>
          <template v-else>
            <template v-for="(it, i) in searchMatches" :key="it.path + it.name">
              <!-- 目录：灰色不可选中的分组标题，按 depth 缩进形成树形层级（无路由、不可跳转） -->
              <div
                v-if="it.dir"
                class="cmd-group cmd-group-dir"
                :style="{ paddingLeft: 10 + (it.depth || 0) * 16 + 'px' }"
              >
                <span class="cmd-item-icon"><MenuIcon :icon="it.icon" /></span>
                <span>{{ it.name }}</span>
                <span class="dir-badge mono">目录</span>
              </div>
              <!-- 菜单：可选中跳转，按 depth 缩进展示树形结构 -->
              <div
                v-else
                class="cmd-item"
                :class="{ active: i === activeIndex }"
                :style="{ marginLeft: (it.depth || 0) * 16 + 'px' }"
                :data-idx="i"
                @mouseenter="activeIndex = i"
                @click="gotoSearchItem(it)"
              >
                <span class="cmd-item-icon"><MenuIcon :icon="it.icon" /></span>
                <span class="cmd-item-name">{{ it.name }}</span>
                <span v-if="it.parents" class="cmd-item-trail mono">{{ it.parents }}</span>
                <span v-show="i === activeIndex" class="kbd">↵</span>
              </div>
            </template>
          </template>
        </div>

        <div class="cmd-foot">
          <span><span class="kbd">↑</span><span class="kbd">↓</span> 切换</span>
          <span><span class="kbd">↵</span> 跳转</span>
          <span><span class="kbd">esc</span> 关闭</span>
          <span v-if="!searchLoading && menuCount" class="cmd-foot-count mono">{{ menuCount }} 项</span>
        </div>
      </div>
    </n-modal>
  </n-layout>
</template>

<script setup lang="ts">
import { computed, h, nextTick, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NDropdown, NButton, NIcon,
  NModal, NForm, NFormItem, NInput, useMessage,
  type DropdownOption, type FormInst, type FormRules,
} from 'naive-ui'
import { MenuOutline, SearchOutline } from '@vicons/ionicons5'
import { useUserStore } from '../stores/user'
import { renderMenuIcon } from '../utils/menuIcon'
import { changePassword, searchMenus } from '../api'
import type { MenuHit, MenuNode } from '../api/types'

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
    // 目录无路由 path，用 code 作 key 保证唯一；菜单叶子仍以 path 为 key，点击即跳转
    key: m.path || m.code,
    icon: renderIcon(m.icon),
    children: m.children?.length ? toOptions(m.children) : undefined,
  }))
}

const menuOptions = computed(() => toOptions(userStore.menus))
const activeKey = computed(() => route.path)

// 侧边栏展开的分组 key 列表（受控）：路由变化时并入当前菜单的祖先目录链，保留手动展开/收起状态
const expandedKeys = ref<string[]>([])

// 在菜单树中查找 key 的祖先链（不含自身），找不到返回 null
function findAncestorKeys(nodes: any[], key: string, trail: string[] = []): string[] | null {
  for (const n of nodes) {
    if (n.key === key) return trail
    if (n.children?.length) {
      const hit = findAncestorKeys(n.children, key, [...trail, n.key])
      if (hit) return hit
    }
  }
  return null
}

// 顶栏搜索跳转/刷新深链接等场景下展开侧边栏对应目录；
// 同时 watch menuOptions：刷新时菜单异步加载，activeKey 不变但菜单树就绪后仍需展开
watch([activeKey, menuOptions], ([p]) => {
  const ancestors = findAncestorKeys(menuOptions.value, p)
  if (!ancestors?.length) return
  expandedKeys.value = [...new Set([...expandedKeys.value, ...ancestors])]
}, { immediate: true })

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
// 与后端 binding 保持一致：新密码 6-20 位（旧密码不限上限，兼容历史长密码），两次输入须一致
const pwdRules: FormRules = {
  oldPassword: [{ required: true, message: '请输入原密码', trigger: ['blur', 'input'] }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: ['blur', 'input'] },
    { min: 6, max: 20, message: '密码长度 6-20 位', trigger: ['blur', 'input'] },
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

// ---- 顶栏全局搜索：命令面板（后端接口搜索 + 自管键盘导航与高亮） ----

// 函数式组件：菜单图标（模板中渲染 renderMenuIcon 的 VNode）
const MenuIcon = (props: { icon?: string }) => renderMenuIcon(props.icon, 16)

const searchKw = ref('')
const searchRef = ref<HTMLInputElement | null>(null)
// 快捷键提示按平台显示（Mac 显示 ⌘K）
const searchKbd = /Mac|iPhone|iPad/.test(navigator.platform || navigator.userAgent) ? '⌘K' : 'Ctrl K'

const showSearch = ref(false)
const activeIndex = ref(0)
const cmdListRef = ref<HTMLElement | null>(null)
// 命中项：空态保持空列表，输入后由防抖接口搜索填充
const searchMatches = ref<MenuHit[]>([])
const searchLoading = ref(false)
// 可选项（菜单）数量：目录仅作分组标题，不参与计数与选中
const menuCount = computed(() => searchMatches.value.filter((m) => !m.dir).length)

function openSearch() {
  searchKw.value = ''
  showSearch.value = true
}

// 弹窗打开后自动聚焦输入框
watch(showSearch, async (v) => {
  if (v) {
    await nextTick()
    searchRef.value?.focus()
  }
})

// 关键词变化：防抖 300ms 调接口搜索；序号防过期响应竞态
let searchTimer: number | undefined
let searchSeq = 0
watch(searchKw, (kw) => {
  activeIndex.value = 0
  window.clearTimeout(searchTimer)
  if (!kw.trim()) {
    searchSeq++
    searchLoading.value = false
    searchMatches.value = []
    return
  }
  searchLoading.value = true
  const seq = ++searchSeq
  searchTimer = window.setTimeout(async () => {
    try {
      const { data } = await searchMenus(kw.trim())
      if (seq === searchSeq) {
        searchMatches.value = data.data
        activeIndex.value = firstSelectable()
      }
    } catch {
      if (seq === searchSeq) searchMatches.value = []
    } finally {
      if (seq === searchSeq) searchLoading.value = false
    }
  }, 300)
})

// 第一个可选项（非目录）下标；无则 -1
function firstSelectable(): number {
  return searchMatches.value.findIndex((m) => !m.dir)
}

// ↑↓ 步进：跳过目录分组标题，只在可选菜单项间移动
function stepSelectable(from: number, delta: 1 | -1): number {
  let i = from
  for (;;) {
    i += delta
    if (i < 0 || i >= searchMatches.value.length) return from
    if (!searchMatches.value[i].dir) return i
  }
}

// 键盘导航：↑↓ 切换高亮、Enter 跳转（Esc 由 n-modal 的 close-on-esc 处理）
function onCmdKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    activeIndex.value = stepSelectable(activeIndex.value, 1)
    scrollActiveIntoView()
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    activeIndex.value = stepSelectable(activeIndex.value, -1)
    scrollActiveIntoView()
  } else if (e.key === 'Enter') {
    e.preventDefault()
    const it = searchMatches.value[activeIndex.value]
    if (it) gotoSearchItem(it)
  }
}

// 高亮项滚动跟随（列表超出高度时保持可见）
function scrollActiveIntoView() {
  nextTick(() => {
    cmdListRef.value?.querySelector(`[data-idx="${activeIndex.value}"]`)?.scrollIntoView({ block: 'nearest' })
  })
}

function gotoSearchItem(it: MenuHit) {
  // 目录无对应路由且渲染为分组标题，正常不可达；防御性拦截
  if (it.dir) return
  showSearch.value = false
  if (it.path !== route.path) {
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
onUnmounted(() => {
  window.removeEventListener('keydown', onGlobalKeydown)
  window.clearTimeout(searchTimer)
})
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
/* 折叠过渡与 naive-ui sider 的宽度动画同步：.3s cubic-bezier(.4, 0, .2, 1) */
.logo {
  height: 64px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
  border-bottom: 1px solid var(--sx-line);
  margin-bottom: 4px;
  transition: padding .3s cubic-bezier(.4, 0, .2, 1), gap .3s cubic-bezier(.4, 0, .2, 1);
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
/* 文字固定宽度，收起时宽度过渡收缩（v-if 瞬时移除会让宽度动画跳变） */
.logo-text {
  display: flex;
  flex-direction: column;
  line-height: 1.25;
  white-space: nowrap;
  flex-shrink: 0;
  width: 92px;
  overflow: hidden;
  transition: width .3s cubic-bezier(.4, 0, .2, 1), opacity .2s cubic-bezier(.4, 0, .2, 1), margin-left .3s cubic-bezier(.4, 0, .2, 1);
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
  transition: padding .3s cubic-bezier(.4, 0, .2, 1);
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
  transition: left .3s cubic-bezier(.4, 0, .2, 1), right .3s cubic-bezier(.4, 0, .2, 1);
}
/* 收起态 hover 背景占满整行（内边距已清零） */
.sider-menu.n-menu--collapsed :deep(.n-menu-item-content::before) {
  left: 0;
  right: 0;
  border-radius: 0;
}
/* 收起态：印章经 padding 精确居中（(64-34)/2 = 15px），文字宽度收缩至 0 */
.logo.logo-collapsed {
  padding: 0 0 0 15px;
  gap: 0;
}
.logo.logo-collapsed .logo-text {
  width: 0;
  opacity: 0;
  margin-left: -10px;
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
/* 后缀与主版本号一起过渡收缩，避免文字瞬时跳变 */
.foot-ext {
  display: inline-block;
  max-width: 80px;
  overflow: hidden;
  white-space: nowrap;
  vertical-align: bottom;
  transition: max-width .3s cubic-bezier(.4, 0, .2, 1), opacity .25s cubic-bezier(.4, 0, .2, 1);
}
.sider-foot.foot-collapsed .foot-ext {
  max-width: 0;
  opacity: 0;
}

/* 顶部栏：eyebrow + 标题的层级读法 */
.main {
  background: var(--sx-bg) !important;
  height: 100%;
}
/* 顶栏固定：main 层自身不滚（默认滚动发生在含 header 的这一层），
   锁为 flex 列 —— header 固定高 + content 吃满剩余并内部滚动 */
.main :deep(.n-layout-scroll-container) {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.main :deep(.n-layout-header) {
  flex-shrink: 0;
}
.main :deep(.n-layout-content) {
  flex: 1;
  min-height: 0;
}
.main :deep(.n-layout-content > .n-scrollbar) {
  height: 100%;
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

/* 搜索命令面板：输入行 / 结果列表 / 快捷键提示栏 */
.cmd-palette {
  width: 520px;
  background: var(--sx-surface);
  border-radius: 12px;
  box-shadow: var(--sx-shadow);
  overflow: hidden;
}
.cmd-input-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 13px 16px;
  border-bottom: 1px solid var(--sx-line);
}
.cmd-search-icon {
  color: var(--sx-muted);
  flex-shrink: 0;
}
.cmd-input {
  flex: 1;
  min-width: 0;
  border: none;
  outline: none;
  background: transparent;
  font-size: 14px;
  color: var(--sx-ink);
}
.cmd-input::placeholder {
  color: var(--sx-muted);
}
.cmd-list {
  max-height: 330px;
  overflow: auto;
  scrollbar-width: thin;
  padding: 6px;
}
.cmd-group {
  font-size: 11px;
  color: var(--sx-muted);
  padding: 7px 10px 3px;
}
/* 目录分组标题：灰色不可选中，仅视觉分层，无 hover/点击交互 */
.cmd-group-dir {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-top: 9px;
  font-size: 12px;
  color: var(--sx-muted);
  user-select: none;
  cursor: default;
}
.cmd-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 7px;
  cursor: pointer;
}
.cmd-item.active {
  background: var(--sx-accent-soft);
}
.cmd-item-icon {
  display: inline-flex;
  align-items: center;
  color: var(--sx-muted);
  flex-shrink: 0;
}
.cmd-item-name {
  font-size: 13px;
  color: var(--sx-ink);
  white-space: nowrap;
}
/* 目录分组标题名称旁的小徽标 */
.dir-badge {
  flex-shrink: 0;
  font-size: 10px;
  color: var(--sx-muted);
  border: 1px solid var(--sx-line);
  border-radius: 4px;
  padding: 0 4px;
  line-height: 16px;
}
.cmd-item-trail {
  margin-left: auto;
  font-size: 11px;
  color: var(--sx-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.cmd-item .kbd {
  flex-shrink: 0;
  opacity: 0;
}
.cmd-item.active .kbd {
  opacity: 1;
}
.cmd-empty {
  padding: 26px 0;
  text-align: center;
  font-size: 12px;
  color: var(--sx-muted);
}
.cmd-foot {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 9px 16px;
  border-top: 1px solid var(--sx-line);
  font-size: 11px;
  color: var(--sx-muted);
}
.cmd-foot-count {
  margin-left: auto;
}

/* 键帽：快捷键提示 */
.kbd {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  padding: 0 4px;
  margin-right: 3px;
  border: 1px solid var(--sx-line);
  border-bottom-width: 2px;
  border-radius: 4px;
  font-family: var(--sx-font-mono);
  font-size: 10px;
  color: var(--sx-muted);
  background: var(--sx-bg);
}
</style>
