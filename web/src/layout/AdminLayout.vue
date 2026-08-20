<template>
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="220" :collapsed="collapsed"
      show-trigger @collapse="collapsed = true" @expand="collapsed = false">
      <div class="logo">
        <span v-if="!collapsed">SmileX Admin</span>
        <span v-else>S</span>
      </div>
      <n-menu :collapsed="collapsed" :collapsed-width="64" :options="menuOptions" :value="activeKey" />
    </n-layout-sider>

    <n-layout>
      <n-layout-header bordered class="header">
        <div class="header-left">
          <n-breadcrumb>
            <n-breadcrumb-item>{{ route.meta?.title || '首页' }}</n-breadcrumb-item>
          </n-breadcrumb>
        </div>
        <n-dropdown :options="userOptions" @select="onUserAction">
          <n-button quaternary>
            {{ userStore.user?.nickname || userStore.user?.username }}
          </n-button>
        </n-dropdown>
      </n-layout-header>
      <n-layout-content content-style="padding: 16px;" :native-scrollbar="false">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NDropdown, NButton, NBreadcrumb, NBreadcrumbItem } from 'naive-ui'
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
.logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 18px;
  color: var(--n-item-text-color);
}
.header {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
}
</style>
