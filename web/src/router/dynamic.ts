import { h } from 'vue'
import { RouterView, type RouteRecordRaw } from 'vue-router'
import router from './index'
import { useUserStore } from '../stores/user'

// 父级菜单（如"系统管理"）没有页面组件，用透传组件渲染子路由。
// 注意必须是组件对象：route.component 若是普通函数会被 vue-router 当作懒加载 loader（要求返回 Promise）。
const Passthrough = { name: 'RouterViewPassthrough', render: () => h(RouterView) }

// 菜单 code -> 页面组件 映射（新增页面在此登记）
const viewModules: Record<string, () => Promise<any>> = {
  'menu:dashboard': () => import('../views/dashboard/Dashboard.vue'),
  'menu:user': () => import('../views/system/Users.vue'),
  'menu:role': () => import('../views/system/Roles.vue'),
  'menu:menu': () => import('../views/system/Menus.vue'),
  'menu:online': () => import('../views/system/Online.vue'),
  'menu:about': () => import('../views/about/About.vue'),
}

// 将后端菜单树转换为路由
export function menuToRoutes(menus: any[]): RouteRecordRaw[] {
  const routes: RouteRecordRaw[] = []
  for (const m of menus ?? []) {
    const comp = m.children?.length && !viewModules[m.code] ? Passthrough : viewModules[m.code]
    const route: RouteRecordRaw = {
      path: m.path,
      name: m.code,
      component: comp,
      meta: { title: m.name, icon: m.icon },
      children: [],
    }
    if (m.children?.length) {
      route.children = menuToRoutes(m.children)
    }
    routes.push(route)
  }
  return routes
}

// 登录后按后端菜单动态注册路由（挂到 layout-root 下），返回首个菜单路径
export async function setupDynamicRoutes(): Promise<string> {
  const userStore = useUserStore()
  await userStore.loadUserContext()
  const dynamic = menuToRoutes(userStore.menus)
  for (const r of dynamic) {
    router.addRoute('layout-root', r)
  }
  // 隐藏路由：个人中心（不在菜单树中，不在侧栏/全局搜索显示）
  if (!router.hasRoute('profile')) {
    router.addRoute('layout-root', {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/profile/Profile.vue'),
      meta: { title: '个人中心' },
    })
  }
  // 菜单路由注册完成后再挂 404 兜底，避免刷新深链接时原始路径被吞（见 ./index.ts）
  if (!router.hasRoute('not-found')) {
    router.addRoute({ path: '/:pathMatch(.*)*', name: 'not-found', redirect: '/' })
  }
  userStore.routesLoaded = true
  return dynamic.length ? dynamic[0].path : '/'
}
