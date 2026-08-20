import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '../stores/user'

// 菜单 code -> 页面组件 映射（新增页面在此登记）
const viewModules: Record<string, () => Promise<any>> = {
  'menu:dashboard': () => import('../views/dashboard/Dashboard.vue'),
  'menu:user': () => import('../views/system/Users.vue'),
  'menu:role': () => import('../views/system/Roles.vue'),
  'menu:permission': () => import('../views/system/Permissions.vue'),
  'menu:menu': () => import('../views/system/Menus.vue'),
}

// 将后端菜单树转换为路由
export function menuToRoutes(menus: any[], parentPath = ''): RouteRecordRaw[] {
  const routes: RouteRecordRaw[] = []
  for (const m of menus ?? []) {
    const comp = viewModules[m.code]
    const route: RouteRecordRaw = {
      path: m.path,
      name: m.code,
      component: comp,
      meta: { title: m.name, icon: m.icon },
      children: [],
    }
    if (m.children?.length) {
      route.children = menuToRoutes(m.children, m.path)
    }
    routes.push(route)
  }
  return routes
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: () => import('../views/login/Login.vue'), meta: { title: '登录' } },
    {
      path: '/',
      name: 'layout-root',
      component: () => import('../layout/AdminLayout.vue'),
      children: [], // 动态菜单路由运行时注入（菜单 path 为绝对路径）
    },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

// 路由守卫：未登录跳登录页；登录后按菜单动态注册路由
router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore()
  if (to.path === '/login') {
    next()
    return
  }
  if (!userStore.accessToken) {
    next('/login')
    return
  }
  if (!userStore.routesLoaded) {
    try {
      await userStore.loadUserContext()
      const dynamic = menuToRoutes(userStore.menus)
      for (const r of dynamic) {
        // 菜单 path 为绝对路径，作为 layout-root 的子路由注册
        router.addRoute('layout-root', r)
      }
      next({ ...to, replace: true })
    } catch {
      userStore.logout()
      next('/login')
    }
    return
  }
  next()
})

router.afterEach((to) => {
  const title = (to.meta?.title as string) || ''
  document.title = title ? `${title} - SmileX Admin` : 'SmileX Admin'
})

export default router
