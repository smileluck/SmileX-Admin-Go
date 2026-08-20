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
      return (code: string) => state.permissions.some((p) => p.code === code)
    },
  },
  actions: {
    async login(username: string, password: string) {
      const { data: resp } = await apiLogin({ username, password })
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
    async logout() {
      try {
        await apiLogout()
      } catch { /* 忽略登出接口异常 */ }
      this.$reset()
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    },
  },
})
