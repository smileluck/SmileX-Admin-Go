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
  // 个人中心 profile 接口额外返回的角色名列表
  role_names?: string[]
}

export interface Permission {
  id: number
  name: string
  code: string
  type: 'dir' | 'menu' | 'button' // dir 目录分组 | menu 菜单页面 | button 按钮权限点（method/path 非空时同时参与后端 RBAC 校验）
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
  type: 'dir' | 'menu' // dir 目录分组（无路由）| menu 菜单页面
  path: string
  icon: string
  sort: number
  children: MenuNode[] | null
}

// 菜单搜索命中项（顶栏命令面板；parents 为父级链提示，不含自身）
// dir 标记目录（含子菜单、无路由，选中时软提示选择具体菜单）
export interface MenuHit {
  name: string
  path: string
  icon: string
  parents: string
  dir?: boolean
  depth?: number
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

// 在线会话（一行 = 一个「用户 × 设备端」会话）
export interface OnlineSession {
  sid: string
  user_id: number
  username: string
  nickname: string
  device: 'web' | 'app' // web 网页端 | app 移动端
  ip: string
  user_agent: string
  login_at: string
  last_active_at: string
  is_current: boolean // 是否当前登录者自己的会话
}

// 登录日志（一次登录尝试 = 一条）
export interface LoginLogInfo {
  id: number
  username: string // 尝试登录的用户名（可能不存在）
  ip: string
  user_agent: string
  device: string // web | app
  status: number // 1 成功 0 失败
  msg: string // 失败原因（成功为空）
  created_at: string
}

// 操作日志（一条写请求审计）
export interface OperationLogInfo {
  id: number
  user_id: number
  username: string
  method: string // POST / PUT / DELETE / PATCH
  path: string // 实际请求路径
  route: string // 路由模板
  action: string // 中文动作名
  params: string // 参数摘要（敏感字段已脱敏）
  ip: string
  user_agent: string
  status_code: number
  latency_ms: number
  created_at: string
}

// 日志列表响应：列表 + 分页 + 保留天数（页面展示保留说明用）
export interface LogPageResult<T> {
  list: T[]
  page: { page: number; page_size: number; total: number }
  retention_days: number
}

// 文件元数据（后端 files 表；object_key 不下发）
export interface FileInfo {
  id: number
  driver: string // 落库时的存储后端：local | oss | cos | tos | minio
  name: string
  ext: string
  size: number
  content_type: string
  uploader_id: number
  uploader_name: string
  created_at: string
}
