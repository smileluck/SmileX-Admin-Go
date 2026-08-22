import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '../stores/user'
import { setupDynamicRoutes } from './dynamic'

// 静态路由：登录页、主布局壳
// 注意：404 兜底不能静态注册——刷新深链接时动态菜单路由尚未注册，静态 catch-all 的
// redirect:'/' 会先把原始路径吞掉（to.path 变成 '/'），导致刷新回到首页。兜底在
// 动态路由注册完成后于 ./dynamic.ts 中补充。
const staticRoutes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: () => import('../views/login/Login.vue'), meta: { title: '登录' } },
  {
    path: '/',
    name: 'layout-root',
    component: () => import('../layout/AdminLayout.vue'),
    children: [], // 动态菜单路由运行时注入（菜单 path 为绝对路径），见 ./dynamic.ts
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes: staticRoutes,
})

// 路由守卫：未登录跳登录页；登录后按菜单动态注册路由（静态/动态分离见 ./dynamic.ts）
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
      const firstPath = await setupDynamicRoutes()
      // 重导航以匹配刚注册的动态路由；'/' 落到第一个菜单
      const target = to.path === '/' ? firstPath : to.path
      next({ path: target, query: to.query, replace: true })
    } catch {
      // token 失效（如 JWT 格式升级/改密后）：清态回登录页。
      // 401 拦截器可能已发起 /login 导航，这里幂等处理避免重复导航竞态导致白屏
      await userStore.logout()
      if (router.currentRoute.value.path === '/login') {
        next(false)
      } else {
        next('/login')
      }
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
