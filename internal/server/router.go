// 路由注册：全部 API 路由按模块分组挂载（handler 见 handler_*.go）。
package server

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	bizauth "github.com/smilex/smilex-admin-gin/internal/biz/auth"
	bizblacklist "github.com/smilex/smilex-admin-gin/internal/biz/blacklist"
	bizlog "github.com/smilex/smilex-admin-gin/internal/biz/log"
	bizperm "github.com/smilex/smilex-admin-gin/internal/biz/permission"
	bizsession "github.com/smilex/smilex-admin-gin/internal/biz/session"
	"github.com/smilex/smilex-admin-gin/internal/server/middleware"
	authsvc "github.com/smilex/smilex-admin-gin/internal/service/auth"
	merchantsvc "github.com/smilex/smilex-admin-gin/internal/service/merchant"
	"github.com/smilex/smilex-admin-gin/pkg/i18n"
	"github.com/smilex/smilex-admin-gin/pkg/logger"
	"github.com/smilex/smilex-admin-gin/pkg/response"
	"go.uber.org/zap"
)

func (s *HTTPServer) registerRoutes() {
	// 持久化 IP 黑名单在 JWT 之前拦截全部 /api/ 请求（静态前端资源不经过此组）
	v1 := s.engine.Group("/api/v1", middleware.IPBlacklist(s.blacklist.Checker()))

	// ---- 公开接口 ----
	// 登录限流（后台与应用用户登录共用）：IP 固定窗口计数，沿用 bl:rl: 前缀，
	// 黑名单提前解封时联动清零；参数见 biz/blacklist 常量
	loginRateLimit := middleware.NewRateLimit(s.rdb, middleware.RateLimitConfig{
		KeyPrefix: "bl:rl:", Max: bizblacklist.LoginRateMax, Window: bizblacklist.LoginRateWindow, MessageKey: "security.login_frequent",
	})
	authg := v1.Group("/auth")
	{
		// 图形验证码：无需鉴权，登录页拉取；按 IP 限流防接口刷量（占用 Redis 存储与带宽）
		authg.GET("/captcha", middleware.NewRateLimit(s.rdb, middleware.RateLimitConfig{
			KeyPrefix: "rl:captcha:", Max: 30, Window: time.Minute, MessageKey: "security.rate_limited",
		}), func(c *gin.Context) {
			vo, err := s.auth.GenerateCaptcha()
			if err != nil {
				response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
				return
			}
			response.OK(c, vo)
		})
		// 登录接口：IP 临时封禁（连续失败拉黑）→ 频率限制 → 登录，防口令爆破
		authg.POST("/login", middleware.LoginIPGuard(s.log, s.blacklist.LoginGuard()), loginRateLimit, func(c *gin.Context) {
			var req authsvc.LoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
				return
			}
			// 登录环境注入（建立会话用；UA 截断防止超长存储）
			req.IP = c.ClientIP()
			req.UserAgent = truncate(c.GetHeader("User-Agent"), 255)
			tp, err := s.auth.Login(c.Request.Context(), req)
			// 登录尝试（成功/失败）均落登录日志（异步，不影响登录响应）
			loginStatus := bizlog.LoginStatusSuccess
			loginMsg := ""
			if err != nil {
				loginStatus = bizlog.LoginStatusFail
				loginMsg = err.Error()
			}
			s.log.RecordLogin(c.Request.Context(), &bizlog.LoginLog{
				Username: req.Username, IP: req.IP, UserAgent: req.UserAgent,
				Device: bizsession.NormalizeDevice(req.DeviceType),
				Status: loginStatus, Msg: loginMsg,
			})
			if err != nil {
				// 验证码错误属入参问题，返回 400 便于前端区分提示并自动刷新
				if errors.Is(err, bizauth.ErrCaptcha) {
					response.FailI18n(c, http.StatusBadRequest, response.CodeErr, err)
					return
				}
				response.FailI18n(c, http.StatusUnauthorized, response.CodeUnauthorized, err)
				return
			}
			response.OK(c, tp)
		})
		authg.POST("/refresh", func(c *gin.Context) {
			var req authsvc.RefreshRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
				return
			}
			tp, err := s.auth.Refresh(c.Request.Context(), req)
			if err != nil {
				response.FailI18n(c, http.StatusUnauthorized, response.CodeUnauthorized, err)
				return
			}
			response.OK(c, tp)
		})
	}

	// ---- 应用用户认证（与后台账号体系 typ 隔离；无验证码、无服务端会话） ----
	// 登录接口挂与后台登录同款的 IP 临时封禁 + 频率限制防护，防口令爆破
	appauthg := v1.Group("/app-auth")
	{
		appauthg.POST("/login", middleware.LoginIPGuard(s.log, s.blacklist.LoginGuard()), loginRateLimit, s.appLogin)
		appauthg.POST("/refresh", s.appRefresh)
	}

	// 应用用户自身数据接口：仅 AppJWT 认证（查库校验启用状态），不做 RBAC
	appAuth := v1.Group("/app-auth", middleware.AppJWT(s.appIssuer, s.appuserUC))
	{
		appAuth.GET("/profile", s.appProfile)
		appAuth.PUT("/password", s.appChangePassword)
	}

	// ---- 自身数据接口：仅 JWT 认证，不做 RBAC ----
	// profile 是本人信息、menus 是已按角色过滤的本人菜单树、logout 无服务端状态，均无越权面；
	// 若纳入默认拒绝的 RBAC，仅绑定了菜单/按钮权限的普通用户登录后即 403 白屏。
	// OpLog 自动审计写请求（登出/改资料/改密码）。
	basic := v1.Group("", middleware.JWT(s.auth), middleware.OpLog(s.log))
	{
		// 登出：吊销当前会话，token 立即失效
		basic.POST("/auth/logout", func(c *gin.Context) {
			if sub := middleware.Subject(c); sub != nil {
				if err := s.auth.Logout(c.Request.Context(), sub.SessionID); err != nil {
					logger.Warn("logout revoke session failed", zap.Error(err))
				}
			}
			response.OK(c, nil)
		})
		basic.GET("/auth/profile", func(c *gin.Context) {
			sub := middleware.Subject(c)
			vo, err := s.auth.Profile(c.Request.Context(), sub.UserID)
			if err != nil {
				response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
				return
			}
			response.OK(c, vo)
		})
		// 本人更新昵称/邮箱
		basic.PUT("/auth/profile", func(c *gin.Context) {
			sub := middleware.Subject(c)
			var req authsvc.UpdateProfileRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
				return
			}
			vo, err := s.auth.UpdateProfile(c.Request.Context(), sub.UserID, req)
			if err != nil {
				response.FailI18n(c, http.StatusBadRequest, response.CodeErr, err)
				return
			}
			response.OK(c, vo)
		})
		// 本人修改密码（校验旧密码；成功后吊销其他端会话，当前端保持登录）
		basic.PUT("/auth/password", func(c *gin.Context) {
			sub := middleware.Subject(c)
			var req authsvc.ChangePasswordRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, i18n.T(c.Request.Context(), "common.invalid_params"))
				return
			}
			sid := ""
			if sub != nil {
				sid = sub.SessionID
			}
			if err := s.auth.ChangePassword(c.Request.Context(), sub.UserID, req, sid); err != nil {
				if errors.Is(err, bizauth.ErrInvalidCredentials) {
					response.BadRequest(c, i18n.T(c.Request.Context(), "profile.wrong_old_password"))
					return
				}
				response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
				return
			}
			response.OK(c, nil)
		})
		// 当前用户可见菜单树
		basic.GET("/menus", func(c *gin.Context) {
			sub := middleware.Subject(c)
			tree, err := s.perm.UserMenuTree(c.Request.Context(), sub.UserID)
			if err != nil {
				response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
				return
			}
			response.OK(c, tree)
		})
		// 菜单搜索（顶栏命令面板）：当前用户可见菜单内按关键词模糊匹配
		basic.GET("/menus/search", func(c *gin.Context) {
			sub := middleware.Subject(c)
			kw := strings.TrimSpace(c.Query("kw"))
			hits, err := s.perm.SearchUserMenus(c.Request.Context(), sub.UserID, kw)
			if err != nil {
				response.FailI18n(c, http.StatusInternalServerError, response.CodeErr, err)
				return
			}
			if hits == nil {
				hits = []*bizperm.MenuHit{}
			}
			response.OK(c, hits)
		})

		// 异步导出：记录归属当前用户，列表/下载/删除均强制按 sub.UserID 过滤（biz 层校验），
		// 与 profile/menus 同属自身数据接口，故仅 JWT 不走 RBAC；导出入口（POST */export）在 protected 组按按钮权限点控制
		exports := basic.Group("/exports")
		{
			// recent=1 返回近期 5 条（任务浮层轮询）；否则分页返回本人记录
			exports.GET("", s.listExports)
			// 下载：云存储 302 到预签名 URL，本地存储后端代理流式输出
			exports.GET("/:id/download", s.downloadExport)
			exports.DELETE("/:id", s.deleteExport)
		}
	}

	// ---- 受保护接口：JWT -> 操作日志（RBAC 拒绝的尝试也记录）-> RBAC ----
	protected := v1.Group("", middleware.JWT(s.auth), middleware.OpLog(s.log), middleware.RBAC(s.auth, s.rbacCache))

	users := protected.Group("/users")
	{
		users.GET("", s.listUsers)
		users.POST("", s.createUser)
		users.GET("/:id", s.getUser)
		users.PUT("/:id", s.updateUser)
		users.DELETE("/:id", s.deleteUser)
		users.PUT("/:id/roles", s.setUserRoles)
		users.PUT("/:id/password", s.resetUserPassword)
		// 导出用户列表（查询条件透传 query，与列表页一致）
		users.POST("/export", func(c *gin.Context) { s.submitExport(c, "user") })
		// 踢某用户全部端下线（挂在 users 资源下，避免与 /online-users/:sid 路由冲突）
		users.DELETE("/:id/sessions", s.kickUserSessions)
	}

	roles := protected.Group("/roles")
	{
		roles.GET("", s.listRoles)
		roles.POST("", s.createRole)
		roles.GET("/:id", s.getRole)
		roles.PUT("/:id", s.updateRole)
		roles.DELETE("/:id", s.deleteRole)
		roles.PUT("/:id/permissions", s.setRolePermissions)
	}

	perms := protected.Group("/permissions")
	{
		perms.GET("", s.listPerms)
		perms.POST("", s.createPerm)
		perms.GET("/:id", s.getPerm)
		perms.PUT("/:id", s.updatePerm)
		perms.DELETE("/:id", s.deletePerm)
	}

	onlines := protected.Group("/online-users")
	{
		onlines.GET("", s.listOnlineUsers)
		onlines.DELETE("/:sid", s.kickOnlineSession)
	}

	loginLogs := protected.Group("/login-logs")
	{
		loginLogs.GET("", s.listLoginLogs)
		loginLogs.DELETE("", s.clearLoginLogs)
		loginLogs.POST("/export", func(c *gin.Context) { s.submitExport(c, "login_log") })
	}

	opLogs := protected.Group("/operation-logs")
	{
		opLogs.GET("", s.listOperationLogs)
		opLogs.DELETE("", s.clearOperationLogs)
		opLogs.POST("/export", func(c *gin.Context) { s.submitExport(c, "op_log") })
	}

	files := protected.Group("/files")
	{
		files.GET("", s.listFiles)
		files.POST("", s.uploadFile)
		// 下载/预览：云存储 302 到预签名 URL，本地存储后端代理流式输出
		files.GET("/:id/raw", s.downloadFile)
		files.DELETE("/:id", s.deleteFile)
	}

	ipBlacklist := protected.Group("/ip-blacklist")
	{
		ipBlacklist.GET("", s.listIPBlacklist)
		ipBlacklist.POST("", s.createIPBlacklist)
		// 解封（软删留痕）
		ipBlacklist.DELETE("/:id", s.deleteIPBlacklist)
	}

	merchants := protected.Group("/merchants")
	{
		merchants.GET("", s.listMerchants)
		merchants.POST("", s.createMerchant)
		merchants.GET("/:id", s.getMerchant)
		merchants.PUT("/:id", s.updateMerchant)
		merchants.DELETE("/:id", s.deleteMerchant)
		// 重置密钥（旧 secret 立即失效，新明文仅此一次返回）
		merchants.PUT("/:id/secret", s.resetMerchantSecret)
		merchants.PUT("/:id/status", s.setMerchantStatus)
	}

	mapiLogs := protected.Group("/merchant-api-logs")
	{
		mapiLogs.GET("", s.listMerchantAPILogs)
	}

	tenants := protected.Group("/tenants")
	{
		tenants.GET("", s.listTenants)
		tenants.POST("", s.createTenant)
		tenants.GET("/:id", s.getTenant)
		tenants.PUT("/:id", s.updateTenant)
		tenants.DELETE("/:id", s.deleteTenant)
		tenants.PUT("/:id/status", s.setTenantStatus)
	}

	appUsers := protected.Group("/app-users")
	{
		appUsers.GET("", s.listAppUsers)
		appUsers.POST("", s.createAppUser)
		appUsers.GET("/:id", s.getAppUser)
		appUsers.PUT("/:id", s.updateAppUser)
		appUsers.DELETE("/:id", s.deleteAppUser)
		// 重置密码（新密码由管理员指定，旧密码立即失效）
		appUsers.PUT("/:id/password", s.resetAppUserPassword)
	}

	// 服务器状态监控（只读快照；CPU%/网卡速率由后台采样器固定 3s 窗口差值计算）
	monitors := protected.Group("/monitor")
	{
		monitors.GET("", s.getServerStatus)
	}

	// 智能体（LLM 配置底座）：供应商 -> 模型 -> Agent
	agentProviders := protected.Group("/agent/providers")
	{
		agentProviders.GET("", s.listAgentProviders)
		agentProviders.POST("", s.createAgentProvider)
		agentProviders.GET("/:id", s.getAgentProvider)
		agentProviders.PUT("/:id", s.updateAgentProvider)
		agentProviders.DELETE("/:id", s.deleteAgentProvider)
		// 连通性测试：真实调用一次上游（model_id 缺省取首个启用模型）
		agentProviders.POST("/:id/test", s.testAgentProvider)
		// 拉取上游 /models 列表（录入辅助）
		agentProviders.GET("/:id/remote-models", s.listAgentRemoteModels)
	}

	agentModels := protected.Group("/agent/models")
	{
		agentModels.GET("", s.listAgentModels)
		agentModels.POST("", s.createAgentModel)
		agentModels.PUT("/:id", s.updateAgentModel)
		agentModels.DELETE("/:id", s.deleteAgentModel)
		agentModels.POST("/:id/test", s.testAgentModel)
	}

	// 可绑定工具清单（Agent 表单多选）
	agentTools := protected.Group("/agent/tools")
	{
		agentTools.GET("", s.listAgentTools)
	}

	// 用量统计（读聚合）
	usage := protected.Group("/agent/usage")
	{
		usage.GET("", s.getAgentUsage)
	}

	// 会话（本人数据，biz 层强制 user_id 过滤）
	convs := protected.Group("/agent/conversations")
	{
		convs.GET("", s.listAgentConversations)
		convs.POST("", s.createAgentConversation)
		convs.GET("/:id/messages", s.listAgentConversationMessages)
		convs.PUT("/:id", s.renameAgentConversation)
		convs.DELETE("/:id", s.deleteAgentConversation)
	}

	agents := protected.Group("/agents")
	{
		agents.GET("", s.listAgents)
		agents.POST("", s.createAgent)
		agents.GET("/:id", s.getAgent)
		agents.PUT("/:id", s.updateAgent)
		agents.DELETE("/:id", s.deleteAgent)
		// 调试对话（Playground，SSE 流式；无状态，不落库）；
		// 按用户限流：每次调用都产生真实上游 token 费用，防误用/滥用刷量
		agents.POST("/:id/chat", middleware.NewRateLimit(s.rdb, middleware.RateLimitConfig{
			KeyPrefix: "rl:agent-chat:", Max: 20, Window: time.Minute, ByUser: true, MessageKey: "security.rate_limited",
		}), s.chatAgent)
	}

	// ---- 开放 API：IP 黑名单 → 商户 HMAC 验签（时间戳偏差 + nonce 防重放） ----
	open := s.engine.Group("/open-api/v1",
		middleware.IPBlacklist(s.blacklist.Checker()),
		middleware.OpenAPIAuth(s.merchantUC, s.rdb, s.cfg.OpenAPI),
	)
	{
		// 探活/凭证自检：返回当前商户脱敏信息
		open.GET("/ping", func(c *gin.Context) {
			response.OK(c, merchantsvc.ToVO(middleware.MerchantFromContext(c)))
		})
	}
}
