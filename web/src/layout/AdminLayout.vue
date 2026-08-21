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
        <n-dropdown :options="userOptions" @select="onUserAction">
          <div class="user-chip">
            <div class="avatar">{{ avatarChar }}</div>
            <span class="user-name">{{ userStore.user?.nickname || userStore.user?.username }}</span>
          </div>
        </n-dropdown>
      </n-layout-header>
      <n-layout-content class="content" content-style="padding: 8px;" :native-scrollbar="false">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NDropdown, NButton, NIcon } from 'naive-ui'
import { MenuOutline } from '@vicons/ionicons5'
import { useUserStore } from '../stores/user'
import { renderMenuIcon } from '../utils/menuIcon'
import type { MenuNode } from '../api/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

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

const userOptions = [
  { label: '退出登录', key: 'logout' },
]

async function onUserAction(key: string) {
  if (key === 'logout') {
    await userStore.logout()
    router.push('/login')
  }
}
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
</style>
