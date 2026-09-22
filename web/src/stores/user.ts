import { defineStore } from 'pinia'
import { getMenus, getProfile, login as apiLogin, logout as apiLogout, refreshToken } from '../api'
import type { MenuNode, Permission, UserInfo } from '../api/types'

export const useUserStore = defineStore('user', {
  state: () => ({
    accessToken: localStorage.getItem('access_token') || '',
    refreshTokenValue: localStorage.getItem('refresh_token') || '',
    user: null as UserInfo | null,
    permissions: [] as Permission[],
    menus: [] as MenuNode[],
    routesLoaded: false,
  }),
  getters: {
    codes: (s) => s.permissions.map((p) => p.code),
    has(state) {
      // 超管角色成员的 permissions 由后端返回全量，无需前端通配码
      return (code: string) => state.permissions.some((p) => p.code === code)
    },
    // 超管角色（id=1）成员：唯一能分配超管角色的人
    isSuper: (s) => !!s.user?.role_ids?.includes(1),
  },
  actions: {
    async login(username: string, password: string, captchaId: string, captchaCode: string) {
      const { data: resp } = await apiLogin({
        username,
        password,
        captcha_id: captchaId,
        captcha_code: captchaCode,
        device_type: 'web', // 设备端标识：同端互斥（新 web 登录顶掉旧 web 会话）
      })
      this.accessToken = resp.data.access_token
      this.refreshTokenValue = resp.data.refresh_token
      localStorage.setItem('access_token', this.accessToken)
      localStorage.setItem('refresh_token', this.refreshTokenValue)
    },
    // 静默刷新；成功返回新 token
    async refresh(): Promise<string | null> {
      if (!this.refreshTokenValue) return null
      try {
        const { data: resp } = await refreshToken(this.refreshTokenValue)
        this.accessToken = resp.data.access_token
        this.refreshTokenValue = resp.data.refresh_token
        localStorage.setItem('access_token', this.accessToken)
        localStorage.setItem('refresh_token', this.refreshTokenValue)
        return this.accessToken
      } catch {
        return null
      }
    },
    // 拉取用户信息 + 菜单（登录后 / 刷新页面后调用）
    async loadUserContext() {
      const [profileRes, menusRes] = await Promise.all([getProfile(), getMenus()])
      this.user = profileRes.data.data.user
      this.permissions = profileRes.data.data.permissions
      this.menus = menusRes.data.data
      this.routesLoaded = true
    },
    // 轻量清态：只清内存与 localStorage，不发网络请求。
    // 供 401 拦截器与路由守卫使用——过期场景下再发 logout 请求只会引发二次 401 竞态。
    // 注意先清 localStorage 再 $reset：$reset 会重跑 state() 工厂读取 localStorage，
    // 顺序颠倒会把旧 token 读回内存
    clearAuth() {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      this.$reset()
    },
    async logout() {
      try {
        await apiLogout()
      } catch { /* 忽略登出接口异常 */ }
      this.clearAuth()
    },
  },
})
