import request from './request'
import type { CaptchaInfo, MenuNode, PageResult, Permission, R, Role, TokenPair, UserInfo } from './types'

// ---- 认证 ----
export const getCaptcha = () => request.get<R<CaptchaInfo>>('/auth/captcha')
export const login = (data: { username: string; password: string; captcha_id: string; captcha_code: string }) =>
  request.post<R<TokenPair>>('/auth/login', data)
export const refreshToken = (refresh_token: string) =>
  request.post<R<TokenPair>>('/auth/refresh', { refresh_token })
export const logout = () => request.post<R<null>>('/auth/logout')
export const getProfile = () =>
  request.get<R<{ user: UserInfo; permissions: Permission[] }>>('/auth/profile')
export const getMenus = () => request.get<R<MenuNode[]>>('/menus')

// ---- 用户 ----
export const listUsers = (params: { page: number; page_size: number; username?: string; status?: number }) =>
  request.get<R<PageResult<UserInfo>>>('/users', { params })
export const createUser = (data: Partial<UserInfo> & { password: string }) =>
  request.post<R<UserInfo>>('/users', data)
export const updateUser = (id: number, data: Partial<UserInfo>) =>
  request.put<R<null>>(`/users/${id}`, data)
export const deleteUser = (id: number) => request.delete<R<null>>(`/users/${id}`)
export const setUserRoles = (id: number, role_ids: number[]) =>
  request.put<R<null>>(`/users/${id}/roles`, { role_ids })
export const resetPassword = (id: number, password: string) =>
  request.put<R<null>>(`/users/${id}/password`, { password })

// ---- 角色 ----
export const listRoles = (params: { page: number; page_size: number; name?: string }) =>
  request.get<R<PageResult<Role>>>('/roles', { params })
export const createRole = (data: Partial<Role>) => request.post<R<Role>>('/roles', data)
export const getRole = (id: number) => request.get<R<Role>>(`/roles/${id}`)
export const updateRole = (id: number, data: Partial<Role>) =>
  request.put<R<null>>(`/roles/${id}`, data)
export const deleteRole = (id: number) => request.delete<R<null>>(`/roles/${id}`)
export const setRolePermissions = (id: number, permission_ids: number[]) =>
  request.put<R<null>>(`/roles/${id}/permissions`, { permission_ids })

// ---- 权限 / 菜单 ----
export const listPermissions = (params: { page: number; page_size: number; type?: string }) =>
  request.get<R<PageResult<Permission>>>('/permissions', { params })
// 全量权限点（page_size=0 不分页）：菜单管理树 / 角色分配权限树需整表构建，分页会静默截断
export const listAllPermissions = (params?: { type?: string }) =>
  request.get<R<PageResult<Permission>>>('/permissions', { params: { page_size: 0, ...params } })
export const createPermission = (data: Partial<Permission>) =>
  request.post<R<Permission>>('/permissions', data)
export const updatePermission = (id: number, data: Partial<Permission>) =>
  request.put<R<null>>(`/permissions/${id}`, data)
export const deletePermission = (id: number) => request.delete<R<null>>(`/permissions/${id}`)
