export interface R<T = any> {
  code: number
  msg: string
  data: T
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_at: string
}

export interface UserInfo {
  id: number
  username: string
  nickname: string
  email: string
  status: number
  role_ids: number[] | null
  created_at: string
}

export interface Permission {
  id: number
  name: string
  code: string
  type: 'api' | 'menu'
  method: string
  path: string
  parent_id: number
  icon: string
  sort: number
}

export interface MenuNode {
  id: number
  name: string
  code: string
  path: string
  icon: string
  sort: number
  children: MenuNode[] | null
}

export interface Role {
  id: number
  name: string
  code: string
  remark: string
  permission_ids?: number[] | null
}

export interface PageResult<T> {
  list: T[]
  page: { page: number; page_size: number; total: number }
}
