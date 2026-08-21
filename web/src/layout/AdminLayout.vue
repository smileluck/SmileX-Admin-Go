<template>
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider class="sider" collapse-mode="width" :collapsed-width="64" :width="232" :collapsed="collapsed"
      :native-scrollbar="false" show-trigger @collapse="collapsed = true" @expand="collapsed = false">
      <div class="logo">
        <div class="logo-mark">S</div>
        <span v-if="!collapsed" class="logo-text">SmileX Admin</span>
      </div>
      <n-menu class="sider-menu" :collapsed="collapsed" :collapsed-width="64" :inverted="true"
        :options="menuOptions" :value="activeKey" @update:value="onMenuSelect" />
    </n-layout-sider>

    <n-layout class="main">
      <n-layout-header bordered class="header">
        <div class="header-left">
          <n-breadcrumb>
            <n-breadcrumb-item>{{ route.meta?.title || '首页' }}</n-breadcrumb-item>
          </n-breadcrumb>
        </div>
        <n-dropdown :options="userOptions" @select="onUserAction">
          <div class="user-chip">
            <div class="avatar">{{ avatarChar }}</div>
            <span>{{ userStore.user?.nickname || userStore.user?.username }}</span>
          </div>
        </n-dropdown>
      </n-layout-header>
      <n-layout-content class="content" content-style="padding: 20px;" :native-scrollbar="false">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NDropdown, NBreadcrumb, NBreadcrumbItem } from 'naive-ui'
import { useUserStore } from '../stores/user'
import type { MenuNode } from '../api/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const collapsed = ref(false)

const iconMap: Record<string, string> = {
  HomeFilled: '🏠', Setting: '⚙️', User: '👤', Avatar: '🎭', Lock: '🔒', Menu: '📋',
}
const renderIcon = (icon?: string) => () => h('span', { style: 'margin-right:6px' }, icon ? (iconMap[icon] || '📄') : '📄')

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
/* 侧边栏：深色渐变 */
.sider {
  background: linear-gradient(180deg, #1e1b4b 0%, #0f172a 60%, #0c1322 100%) !important;
}
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  margin-bottom: 8px;
}
.logo-mark {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 9px;
  font-size: 16px;
  font-weight: 800;
  color: #fff;
  background: linear-gradient(135deg, #6366f1, #a855f7);
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.45);
}
.logo-text {
  font-weight: 700;
  font-size: 17px;
  color: #fff;
  letter-spacing: 0.5px;
  white-space: nowrap;
}
.sider-menu {
  background: transparent !important;
}
.sider-menu :deep(.n-menu .n-menu-item-content::before) {
  left: 10px;
  right: 10px;
  border-radius: 8px;
}

/* 顶部栏 */
.main {
  background: #f3f5f9 !important;
}
.header {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(15, 23, 42, 0.06);
  position: relative;
  z-index: 1;
}
.user-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 5px 12px 5px 5px;
  border-radius: 999px;
  cursor: pointer;
  transition: background 0.2s ease;
  font-size: 14px;
}
.user-chip:hover {
  background: #f1f3f7;
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
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
}

/* 内容区 */
.content {
  background: transparent !important;
}

/* 路由切换淡入 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.18s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
