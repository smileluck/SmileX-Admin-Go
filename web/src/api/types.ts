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

export interface CaptchaInfo {
  captcha_id: string
  // PNG 的 dataURL（data:image/png;base64,...），可直接作 <img src>
  captcha_image: string
  // false 表示服务端已停用验证码（本地调试），前端隐藏验证码表单
  enabled?: boolean
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
  type: 'menu' | 'button' // menu 菜单 | button 按钮权限点（method/path 非空时同时参与后端 RBAC 校验）
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
  remark: string
  permission_ids?: number[] | null
}

export interface PageResult<T> {
  list: T[]
  page: { page: number; page_size: number; total: number }
}
