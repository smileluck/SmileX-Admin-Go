import request from './request'
import type { AgentConversation, AgentConversationMessage, AgentInfo, AgentModel, AgentProvider, AgentTestResult, AppUser, BlacklistItem, CaptchaInfo, ExportRecord, FileInfo, LoginLogInfo, MenuHit, MenuNode, Merchant, MerchantAPILog, NotifyChannel, NotifyRecord, NotifyRule, NotifyTestResult, OnlineSession, OperationLogInfo, PageResult, Permission, R, Role, ServerStatus, Tenant, TokenPair, DashboardStats, DictItem, MonitorHistoryPoint, DictType, JobHandler, JobInfo, JobLog, NoticeInfo, NoticeRoleOption, NoticeUserOption, SysConfig, UserInfo, UsageStats, LogPageResult } from './types'

// ---- 认证 ----
export const getCaptcha = () => request.get<R<CaptchaInfo>>('/auth/captcha')
export const login = (data: { username: string; password: string; captcha_id: string; captcha_code: string; device_type?: string }) =>
  request.post<R<TokenPair>>('/auth/login', data)
export const refreshToken = (refresh_token: string) =>
  request.post<R<TokenPair>>('/auth/refresh', { refresh_token })
export const logout = () => request.post<R<null>>('/auth/logout')
export const getProfile = () =>
  request.get<R<{ user: UserInfo; permissions: Permission[] }>>('/auth/profile')
export const getMenus = () => request.get<R<MenuNode[]>>('/menus')
// 菜单搜索（顶栏命令面板）：后端在当前用户可见菜单内模糊匹配，空关键词返回空
export const searchMenus = (kw: string) =>
  request.get<R<MenuHit[]>>('/menus/search', { params: { kw } })
// 本人更新昵称/邮箱，返回更新后的完整个人信息
export const updateProfile = (data: { nickname?: string; email?: string }) =>
  request.put<R<{ user: UserInfo; permissions: Permission[] }>>('/auth/profile', data)
// 本人修改密码（校验旧密码）；成功后需重新登录
export const changePassword = (data: { old_password: string; new_password: string }) =>
  request.put<R<null>>('/auth/password', data)

// ---- 用户 ----
export const listUsers = (params: { page: number; page_size: number; username?: string; status?: number; reveal?: number }) =>
  request.get<R<PageResult<UserInfo>>>('/users', { params })
export const createUser = (data: Partial<UserInfo> & { password: string }) =>
  request.post<R<UserInfo>>('/users', data)
export const updateUser = (id: number, data: Partial<UserInfo>) =>
  request.put<R<null>>(`/users/${id}`, data)
export const deleteUser = (id: number) => request.delete<R<null>>(`/users/${id}`)
export const getUser = (id: number) => request.get<R<UserInfo>>(`/users/${id}`)
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

// ---- 在线用户 ----
export const listOnlineUsers = (params: { page: number; page_size: number; username?: string; device?: string }) =>
  request.get<R<PageResult<OnlineSession>>>('/online-users', { params })
// 踢单个会话下线（某端）
export const kickOnlineSession = (sid: string) => request.delete<R<null>>(`/online-users/${sid}`)
// 踢某用户全部端下线
export const kickUserSessions = (userId: number) => request.delete<R<null>>(`/users/${userId}/sessions`)

// ---- 日志 ----
// start/end 为 unix 秒级时间戳
export const listLoginLogs = (params: { page: number; page_size: number; username?: string; ip?: string; status?: number; start?: number; end?: number }) =>
  request.get<R<LogPageResult<LoginLogInfo>>>('/login-logs', { params })
export const clearLoginLogs = () => request.delete<R<{ deleted: number }>>('/login-logs')
export const listOperationLogs = (params: { page: number; page_size: number; username?: string; method?: string; kw?: string; start?: number; end?: number }) =>
  request.get<R<LogPageResult<OperationLogInfo>>>('/operation-logs', { params })
export const clearOperationLogs = () => request.delete<R<{ deleted: number }>>('/operation-logs')

// ---- 商户（商户接入）----
export const listMerchants = (params: { page: number; page_size: number; name?: string; code?: string; app_key?: string; status?: number; reveal?: number }) =>
  request.get<R<PageResult<Merchant>>>('/merchants', { params })
// 创建成功时 app_secret 仅此一次明文返回，前端需弹窗展示并提示保存
export const createMerchant = (data: { name: string; code: string; contact_name?: string; contact_phone?: string; contact_email?: string; remark?: string }) =>
  request.post<R<Merchant & { app_secret: string }>>('/merchants', data)
export const getMerchant = (id: number) => request.get<R<Merchant>>(`/merchants/${id}`)
// 联系人字段留空表示不修改（列表/详情返回的均为脱敏值，不可原样回传）
export const updateMerchant = (id: number, data: { name: string; contact_name?: string; contact_phone?: string; contact_email?: string; remark?: string }) =>
  request.put<R<null>>(`/merchants/${id}`, data)
export const deleteMerchant = (id: number) => request.delete<R<null>>(`/merchants/${id}`)
// 重置密钥：新 app_secret 仅此一次明文返回
export const resetMerchantSecret = (id: number) =>
  request.put<R<{ app_key: string; app_secret: string }>>(`/merchants/${id}/secret`)
export const setMerchantStatus = (id: number, status: number) =>
  request.put<R<null>>(`/merchants/${id}/status`, { status })
// start/end 为 unix 秒级时间戳
export const listMerchantAPILogs = (params: { page: number; page_size: number; app_key?: string; path?: string; status_code?: number; start?: number; end?: number; reveal?: number }) =>
  request.get<R<PageResult<MerchantAPILog>>>('/merchant-api-logs', { params })

// ---- 租户 ----
export const listTenants = (params: { page: number; page_size: number; name?: string; code?: string; status?: number }) =>
  request.get<R<PageResult<Tenant>>>('/tenants', { params })
export const createTenant = (data: { name: string; code: string; contact_name?: string; contact_phone?: string; remark?: string }) =>
  request.post<R<Tenant>>('/tenants', data)
export const getTenant = (id: number) => request.get<R<Tenant>>(`/tenants/${id}`)
// code 创建后不可修改，更新入参不含 code
export const updateTenant = (id: number, data: { name: string; contact_name?: string; contact_phone?: string; remark?: string }) =>
  request.put<R<null>>(`/tenants/${id}`, data)
export const deleteTenant = (id: number) => request.delete<R<null>>(`/tenants/${id}`)
export const setTenantStatus = (id: number, status: number) =>
  request.put<R<null>>(`/tenants/${id}/status`, { status })

// ---- 应用用户 ----
// kw 模糊匹配用户名/昵称，phone 精确匹配，tenant_id 按租户筛选
export const listAppUsers = (params: { page: number; page_size: number; kw?: string; phone?: string; status?: number; tenant_id?: number }) =>
  request.get<R<PageResult<AppUser>>>('/app-users', { params })
export const createAppUser = (data: { username: string; password: string; nickname?: string; phone?: string; email?: string; tenant_ids?: number[] }) =>
  request.post<R<AppUser>>('/app-users', data)
export const getAppUser = (id: number) => request.get<R<AppUser>>(`/app-users/${id}`)
// username 创建后不可修改；tenant_ids 全量替换；status 可选（状态切换也走此接口）
export const updateAppUser = (id: number, data: { nickname?: string; phone?: string; email?: string; status?: number; tenant_ids?: number[] }) =>
  request.put<R<null>>(`/app-users/${id}`, data)
export const deleteAppUser = (id: number) => request.delete<R<null>>(`/app-users/${id}`)
export const setAppUserStatus = (id: number, status: number) =>
  request.put<R<null>>(`/app-users/${id}`, { status })
// 重置密码（新密码由管理员指定，旧密码立即失效）
export const resetAppUserPassword = (id: number, password: string) =>
  request.put<R<null>>(`/app-users/${id}/password`, { password })

// ---- 服务器状态监控 ----
// 只读快照：CPU%/网卡速率为后端 3s 窗口差值，页面轮询读最新值
export const getServerStatus = () => request.get<R<ServerStatus>>('/monitor')

// ---- IP 黑名单 ----
export const listBlacklist = (params: { page: number; page_size: number; ip?: string }) =>
  request.get<R<PageResult<BlacklistItem>>>('/ip-blacklist', { params })
export const createBlacklist = (data: { ip: string; reason?: string; expire_at?: number | null }) =>
  request.post<R<BlacklistItem>>('/ip-blacklist', data)
export const deleteBlacklist = (id: number) => request.delete<R<null>>(`/ip-blacklist/${id}`)

// ---- 文件 ----
export const listFiles = (params: { page: number; page_size: number; name?: string }) =>
  request.get<R<PageResult<FileInfo>>>('/files', { params })
// 大文件上传：覆写默认 15s 超时
export const uploadFile = (file: File) => {
  const fd = new FormData()
  fd.append('file', file)
  return request.post<R<FileInfo>>('/files', fd, { timeout: 0 })
}
export const deleteFile = (id: number) => request.delete<R<null>>(`/files/${id}`)
// 下载/预览均走鉴权接口：云存储会 302 到预签名 URL（axios 自动跟随），统一按 blob 取回
export const getFileBlob = (id: number, download = false) =>
  request.get<Blob>(`/files/${id}/raw`, { responseType: 'blob', timeout: 0, params: download ? { download: 1 } : {} })

// ---- 异步导出 ----
// 提交导出任务：params 为当前列表过滤条件（剔除 page/page_size 与空值）；429 表示队列满
export const createExport = (biz: 'users' | 'login-logs' | 'operation-logs', params: Record<string, any>) => {
  const query: Record<string, any> = {}
  for (const [k, v] of Object.entries(params)) {
    if (k === 'page' || k === 'page_size') continue
    if (v === undefined || v === null || v === '') continue
    query[k] = v
  }
  return request.post<R<ExportRecord>>(`/${biz}/export`, null, { params: query })
}
// 近期 5 条导出记录（顶栏悬浮框）
export const listRecentExports = () => request.get<R<ExportRecord[]>>('/exports', { params: { recent: 1 } })
// 本人导出记录分页
export const listExportRecords = (params: { page: number; page_size: number }) =>
  request.get<R<PageResult<ExportRecord>>>('/exports', { params })
// 导出文件下载：大文件覆写默认 15s 超时；409 未完成 / 403 非本人
export const getExportBlob = (id: number) =>
  request.get<Blob>(`/exports/${id}/download`, { responseType: 'blob', timeout: 0 })
export const deleteExport = (id: number) => request.delete<R<null>>(`/exports/${id}`)

// ---- 智能体（LLM 配置底座）----
// 供应商
export const listAgentProviders = (params: { page: number; page_size: number; name?: string; code?: string; status?: number }) =>
  request.get<R<PageResult<AgentProvider>>>('/agent/providers', { params })
export const createAgentProvider = (data: Partial<AgentProvider> & { api_key?: string }) =>
  request.post<R<AgentProvider>>('/agent/providers', data)
export const getAgentProvider = (id: number) => request.get<R<AgentProvider>>(`/agent/providers/${id}`)
export const updateAgentProvider = (id: number, data: Partial<AgentProvider> & { api_key?: string }) =>
  request.put<R<null>>(`/agent/providers/${id}`, data)
export const deleteAgentProvider = (id: number) => request.delete<R<null>>(`/agent/providers/${id}`)
// 连通性测试：真实调用一次上游（model_id 缺省取该供应商首个启用模型）
export const testAgentProvider = (id: number, model_id?: number) =>
  request.post<R<AgentTestResult>>(`/agent/providers/${id}/test`, { model_id })
// 拉取上游 /models 列表（录入辅助）
export const listAgentRemoteModels = (id: number) =>
  request.get<R<string[]>>(`/agent/providers/${id}/remote-models`)
// 模型
export const listAgentModels = (params: { page: number; page_size: number; provider_id?: number; name?: string; status?: number }) =>
  request.get<R<PageResult<AgentModel>>>('/agent/models', { params })
export const createAgentModel = (data: Partial<AgentModel>) =>
  request.post<R<AgentModel>>('/agent/models', data)
export const updateAgentModel = (id: number, data: Partial<AgentModel>) =>
  request.put<R<null>>(`/agent/models/${id}`, data)
export const deleteAgentModel = (id: number) => request.delete<R<null>>(`/agent/models/${id}`)
export const testAgentModel = (id: number) => request.post<R<AgentTestResult>>(`/agent/models/${id}/test`)
// Agent 配置
export const listAgents = (params: { page: number; page_size: number; name?: string; code?: string; status?: number }) =>
  request.get<R<PageResult<AgentInfo>>>('/agents', { params })
export const createAgent = (data: Partial<AgentInfo>) => request.post<R<AgentInfo>>('/agents', data)
export const getAgent = (id: number) => request.get<R<AgentInfo>>(`/agents/${id}`)
export const updateAgent = (id: number, data: Partial<AgentInfo>) =>
  request.put<R<null>>(`/agents/${id}`, data)
export const deleteAgent = (id: number) => request.delete<R<null>>(`/agents/${id}`)

// ---- 会话（本人数据） ----

export const listAgentConversations = (params: { page: number; page_size: number; agent_id?: number }) =>
  request.get<R<PageResult<AgentConversation>>>('/agent/conversations', { params })
export const createAgentConversation = (agent_id: number) =>
  request.post<R<AgentConversation>>('/agent/conversations', { agent_id })
export const renameAgentConversation = (id: number, title: string) =>
  request.put<R<AgentConversation>>(`/agent/conversations/${id}`, { title })
export const deleteAgentConversation = (id: number) =>
  request.delete<R<null>>(`/agent/conversations/${id}`)
export const listAgentConversationMessages = (id: number, params: { page: number; page_size: number }) =>
  request.get<R<PageResult<AgentConversationMessage>>>(`/agent/conversations/${id}/messages`, { params })

export const getAgentUsage = (days = 7) => request.get<R<UsageStats>>('/agent/usage', { params: { days } })

export const listAgentTools = () => request.get<R<string[]>>('/agent/tools')

// ---- 数据字典 ----

export const listDictTypes = (params: { page: number; page_size: number; name?: string; code?: string }) =>
  request.get<R<PageResult<DictType>>>('/dict-types', { params })
export const createDictType = (data: Partial<DictType>) => request.post<R<DictType>>('/dict-types', data)
export const updateDictType = (id: number, data: Partial<DictType>) => request.put<R<null>>(`/dict-types/${id}`, data)
export const deleteDictType = (id: number) => request.delete<R<null>>(`/dict-types/${id}`)
export const listDictItems = (typeID: number, params: { page: number; page_size: number }) =>
  request.get<R<PageResult<DictItem>>>(`/dict-types/${typeID}/items`, { params })
export const createDictItem = (typeID: number, data: Partial<DictItem>) => request.post<R<DictItem>>(`/dict-types/${typeID}/items`, data)
export const updateDictItem = (id: number, data: Partial<DictItem>) => request.put<R<null>>(`/dict-items/${id}`, data)
export const deleteDictItem = (id: number) => request.delete<R<null>>(`/dict-items/${id}`)
export const getDictItemsByCode = (code: string) => request.get<R<DictItem[]>>(`/dicts/${code}/items`)

// ---- 系统参数 ----

export const listSysConfigs = (keyword?: string) =>
  request.get<R<SysConfig[]>>('/sys-configs', { params: { keyword } })
export const createSysConfig = (data: Partial<SysConfig>) => request.post<R<SysConfig>>('/sys-configs', data)
export const updateSysConfig = (key: string, data: { value: string }) =>
  request.put<R<SysConfig>>(`/sys-configs/${key}`, data)
export const deleteSysConfig = (key: string) => request.delete<R<null>>(`/sys-configs/${key}`)

// ---- 通知公告 ----

export const listNotices = (params: { page: number; page_size: number; title?: string; level?: string; status?: string }) =>
  request.get<R<PageResult<NoticeInfo>>>('/notices', { params })
export const createNotice = (data: Partial<NoticeInfo>) => request.post<R<NoticeInfo>>('/notices', data)
export const getNotice = (id: number) => request.get<R<NoticeInfo>>(`/notices/${id}`)
export const updateNotice = (id: number, data: Partial<NoticeInfo>) => request.put<R<null>>(`/notices/${id}`, data)
export const deleteNotice = (id: number) => request.delete<R<null>>(`/notices/${id}`)
// 送达范围选项（发布表单；公告权限即可，无需角色/用户管理权限）
export const listNoticeRoleOptions = () => request.get<R<NoticeRoleOption[]>>('/notices/options/roles')
export const listNoticeUserOptions = (kw: string) =>
  request.get<R<NoticeUserOption[]>>('/notices/options/users', { params: { kw } })
// 消费端（basic）
export const listActiveNotices = () => request.get<R<NoticeInfo[]>>('/notices/active')
export const getUnreadNoticeCount = () => request.get<R<{ count: number }>>('/notices/unread-count')
export const markNoticeRead = (id: number) => request.post<R<null>>(`/notices/${id}/read`)

// ---- 告警通知：渠道 ----

export const listNotifyChannels = (params?: { kw?: string; type?: string; status?: number }) =>
  request.get<R<NotifyChannel[]>>('/notify/channels', { params })
export const createNotifyChannel = (data: Partial<NotifyChannel> & { smtp_password?: string; webhook_secret?: string }) =>
  request.post<R<NotifyChannel>>('/notify/channels', data)
export const getNotifyChannel = (id: number) => request.get<R<NotifyChannel>>(`/notify/channels/${id}`)
export const updateNotifyChannel = (id: number, data: Partial<NotifyChannel> & { smtp_password?: string; webhook_secret?: string }) =>
  request.put<R<null>>(`/notify/channels/${id}`, data)
export const deleteNotifyChannel = (id: number) => request.delete<R<null>>(`/notify/channels/${id}`)
export const testNotifyChannel = (id: number) => request.post<R<NotifyTestResult>>(`/notify/channels/${id}/test`)

// ---- 告警通知：规则 ----

export const listNotifyRules = (params?: { kw?: string; source?: string; status?: number }) =>
  request.get<R<NotifyRule[]>>('/notify/rules', { params })
export const createNotifyRule = (data: Partial<NotifyRule>) => request.post<R<NotifyRule>>('/notify/rules', data)
export const getNotifyRule = (id: number) => request.get<R<NotifyRule>>(`/notify/rules/${id}`)
export const updateNotifyRule = (id: number, data: Partial<NotifyRule>) => request.put<R<null>>(`/notify/rules/${id}`, data)
export const deleteNotifyRule = (id: number) => request.delete<R<null>>(`/notify/rules/${id}`)

// ---- 告警通知：发送记录 ----

export const listNotifyRecords = (params: { page: number; page_size: number; channel_id?: number; source?: string; status?: string }) =>
  request.get<R<PageResult<NotifyRecord>>>('/notify/records', { params })
export const clearNotifyRecords = () => request.delete<R<null>>('/notify/records')

// ---- 定时任务 ----

export const listJobs = (params: { page: number; page_size: number; name?: string }) =>
  request.get<R<PageResult<JobInfo>>>('/jobs', { params })
export const listJobHandlers = () => request.get<R<JobHandler[]>>('/jobs/handlers')
export const createJob = (data: Partial<JobInfo>) => request.post<R<JobInfo>>('/jobs', data)
export const getJob = (id: number) => request.get<R<JobInfo>>(`/jobs/${id}`)
export const updateJob = (id: number, data: Partial<JobInfo>) => request.put<R<null>>(`/jobs/${id}`, data)
export const setJobStatus = (id: number, status: number) => request.put<R<null>>(`/jobs/${id}/status`, { status })
export const deleteJob = (id: number) => request.delete<R<null>>(`/jobs/${id}`)
export const runJobOnce = (id: number) => request.post<R<null>>(`/jobs/${id}/run`)
export const listJobLogs = (id: number, params: { page: number; page_size: number }) =>
  request.get<R<PageResult<JobLog>>>(`/jobs/${id}/logs`, { params })

// ---- 仪表盘 ----

export const getDashboardStats = () => request.get<R<DashboardStats>>('/dashboard/stats')

// ---- 监控历史 ----

export const getMonitorHistory = (hours = 24) =>
  request.get<R<MonitorHistoryPoint[]>>('/monitor/history', { params: { hours } })
