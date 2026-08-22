// Package server 传输层：Gin HTTP Server 与路由注册。
// 切换 Kratos 时本层是唯一需要替换的层（由 proto 生成的 HTTP/gRPC server 代替）。
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	bizauth "github.com/smilex/smilex-admin-gin/internal/biz/auth"
	bizperm "github.com/smilex/smilex-admin-gin/internal/biz/permission"
	"github.com/smilex/smilex-admin-gin/internal/biz/role"
	bizsession "github.com/smilex/smilex-admin-gin/internal/biz/session"
	"github.com/smilex/smilex-admin-gin/internal/biz/user"
	"github.com/smilex/smilex-admin-gin/internal/conf"
	"github.com/smilex/smilex-admin-gin/internal/server/middleware"
	authsvc "github.com/smilex/smilex-admin-gin/internal/service/auth"
	permsvc "github.com/smilex/smilex-admin-gin/internal/service/permission"
	rolesvc "github.com/smilex/smilex-admin-gin/internal/service/role"
	sessionsvc "github.com/smilex/smilex-admin-gin/internal/service/session"
	usersvc "github.com/smilex/smilex-admin-gin/internal/service/user"
	"github.com/smilex/smilex-admin-gin/pkg/logger"
	"github.com/smilex/smilex-admin-gin/pkg/response"
	"go.uber.org/zap"
)

// HTTPServer 聚合全部应用服务
type HTTPServer struct {
	cfg     *conf.Bootstrap
	auth    *authsvc.Service
	user    *usersvc.Service
	role    *rolesvc.Service
	perm    *permsvc.Service
	session *sessionsvc.Service
	engine  *gin.Engine
	srv     *http.Server
}

// NewHTTPServer 构造并注册路由
func NewHTTPServer(cfg *conf.Bootstrap, auth *authsvc.Service, user *usersvc.Service,
	role *rolesvc.Service, perm *permsvc.Service, session *sessionsvc.Service) *HTTPServer {
	gin.SetMode(cfg.Server.Mode)
	e := gin.New()
	e.Use(gin.Recovery(),
		middleware.SecurityHeaders(),
		middleware.CORS(),
		middleware.XSSFilter(),
		middleware.SQLInjectionGuard(),
	)

	s := &HTTPServer{cfg: cfg, auth: auth, user: user, role: role, perm: perm, session: session, engine: e}
	s.registerRoutes()
	s.registerStatic()

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: e,
	}
	return s
}

func (s *HTTPServer) registerRoutes() {
	v1 := s.engine.Group("/api/v1")

	// ---- 公开接口 ----
	authg := v1.Group("/auth")
	{
		// 图形验证码：无需鉴权，登录页拉取
		authg.GET("/captcha", func(c *gin.Context) {
			vo, err := s.auth.GenerateCaptcha()
			if err != nil {
				response.ServerError(c, err.Error())
				return
			}
			response.OK(c, vo)
		})
		// 登录接口按 IP 限流，防口令爆破
		authg.POST("/login", middleware.LoginRateLimit(), func(c *gin.Context) {
			var req authsvc.LoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, err.Error())
				return
			}
			// 登录环境注入（建立会话用；UA 截断防止超长存储）
			req.IP = c.ClientIP()
			req.UserAgent = truncate(c.GetHeader("User-Agent"), 255)
			tp, err := s.auth.Login(c.Request.Context(), req)
			if err != nil {
				// 验证码错误属入参问题，返回 400 便于前端区分提示并自动刷新
				if errors.Is(err, bizauth.ErrCaptcha) {
					response.BadRequest(c, err.Error())
					return
				}
				response.Unauthorized(c, err.Error())
				return
			}
			response.OK(c, tp)
		})
		authg.POST("/refresh", func(c *gin.Context) {
			var req authsvc.RefreshRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, err.Error())
				return
			}
			tp, err := s.auth.Refresh(c.Request.Context(), req)
			if err != nil {
				response.Unauthorized(c, err.Error())
				return
			}
			response.OK(c, tp)
		})
	}

	// ---- 自身数据接口：仅 JWT 认证，不做 RBAC ----
	// profile 是本人信息、menus 是已按角色过滤的本人菜单树、logout 无服务端状态，均无越权面；
	// 若纳入默认拒绝的 RBAC，仅绑定了菜单/按钮权限的普通用户登录后即 403 白屏。
	basic := v1.Group("", middleware.JWT(s.auth))
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
				response.ServerError(c, err.Error())
				return
			}
			response.OK(c, vo)
		})
		// 本人更新昵称/邮箱
		basic.PUT("/auth/profile", func(c *gin.Context) {
			sub := middleware.Subject(c)
			var req authsvc.UpdateProfileRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, err.Error())
				return
			}
			vo, err := s.auth.UpdateProfile(c.Request.Context(), sub.UserID, req)
			if err != nil {
				response.BadRequest(c, err.Error())
				return
			}
			response.OK(c, vo)
		})
		// 本人修改密码（校验旧密码；成功后吊销其他端会话，当前端保持登录）
		basic.PUT("/auth/password", func(c *gin.Context) {
			sub := middleware.Subject(c)
			var req authsvc.ChangePasswordRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, err.Error())
				return
			}
			sid := ""
			if sub != nil {
				sid = sub.SessionID
			}
			if err := s.auth.ChangePassword(c.Request.Context(), sub.UserID, req, sid); err != nil {
				if errors.Is(err, bizauth.ErrInvalidCredentials) {
					response.BadRequest(c, "原密码不正确")
					return
				}
				response.ServerError(c, err.Error())
				return
			}
			response.OK(c, nil)
		})
		// 当前用户可见菜单树
		basic.GET("/menus", func(c *gin.Context) {
			sub := middleware.Subject(c)
			tree, err := s.perm.UserMenuTree(c.Request.Context(), sub.UserID)
			if err != nil {
				response.ServerError(c, err.Error())
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
				response.ServerError(c, err.Error())
				return
			}
			if hits == nil {
				hits = []*bizperm.MenuHit{}
			}
			response.OK(c, hits)
		})
	}

	// ---- 受保护接口：JWT -> RBAC ----
	protected := v1.Group("", middleware.JWT(s.auth), middleware.RBAC(s.auth))

	users := protected.Group("/users")
	{
		users.GET("", s.listUsers)
		users.POST("", s.createUser)
		users.GET("/:id", s.getUser)
		users.PUT("/:id", s.updateUser)
		users.DELETE("/:id", s.deleteUser)
		users.PUT("/:id/roles", s.setUserRoles)
		users.PUT("/:id/password", s.resetUserPassword)
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
}

// Start 启动 HTTP 服务（阻塞）
func (s *HTTPServer) Start() error {
	return s.srv.ListenAndServe()
}

// registerStatic 托管前端 SPA 产物（web/dist），存在时启用；SPA history 路由 fallback 到 index.html
func (s *HTTPServer) registerStatic() {
	dir := s.cfg.Server.StaticDir
	if dir == "" {
		dir = "web/dist"
	}
	if index, err := filepath.Abs(filepath.Join(dir, "index.html")); err == nil {
		if _, err := os.Stat(index); err == nil {
			// assets 文件名带 hash，可长期缓存；index.html 禁止缓存避免发版后白屏
			s.engine.Static("/assets", filepath.Join(dir, "assets"))
			s.engine.NoRoute(func(c *gin.Context) {
				if strings.HasPrefix(c.Request.URL.Path, "/api/") {
					response.NotFound(c, "not found")
					return
				}
				c.Header("Cache-Control", "no-cache")
				c.File(index)
			})
			logger.Info("serving static frontend", zap.String("dir", dir))
		}
	}
}

// Stop 优雅关停
func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// ---- 用户 ----

type listResult struct {
	List interface{} `json:"list"`
	Page interface{} `json:"page"`
}

func (s *HTTPServer) listUsers(c *gin.Context) {
	page, size := pageParams(c)
	q := user.Query{Username: c.Query("username")}
	if v := c.Query("status"); v != "" {
		if st, err := strconv.Atoi(v); err == nil {
			q.Status = &st
		}
	}
	users, pg, err := s.user.List(c.Request.Context(), q, page, size)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, listResult{List: users, Page: pg})
}

func (s *HTTPServer) createUser(c *gin.Context) {
	var req usersvc.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	vo, err := s.user.Create(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, vo)
}

func (s *HTTPServer) getUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	vo, err := s.user.Get(c.Request.Context(), id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, vo)
}

func (s *HTTPServer) updateUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req usersvc.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := s.user.Update(user.WithOperator(c.Request.Context(), middleware.Subject(c).UserID), id, req); err != nil {
		s.userErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) deleteUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := s.user.Delete(c.Request.Context(), id); err != nil {
		s.userErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) setUserRoles(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req usersvc.SetRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := s.user.SetRoles(user.WithOperator(c.Request.Context(), middleware.Subject(c).UserID), id, req); err != nil {
		s.userErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) resetUserPassword(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req usersvc.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := s.user.ResetPassword(user.WithOperator(c.Request.Context(), middleware.Subject(c).UserID), id, req); err != nil {
		s.userErr(c, err)
		return
	}
	response.OK(c, nil)
}

// ---- 角色 ----

func (s *HTTPServer) listRoles(c *gin.Context) {
	page, size := pageParams(c)
	roles, pg, err := s.role.List(c.Request.Context(), role.Query{Name: c.Query("name")}, page, size)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, listResult{List: roles, Page: pg})
}

func (s *HTTPServer) createRole(c *gin.Context) {
	var req rolesvc.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	r, err := s.role.Create(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, r)
}

func (s *HTTPServer) getRole(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	r, err := s.role.Get(c.Request.Context(), id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, r)
}

func (s *HTTPServer) updateRole(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req rolesvc.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := s.role.Update(c.Request.Context(), id, req); err != nil {
		s.roleErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) deleteRole(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := s.role.Delete(c.Request.Context(), id); err != nil {
		s.roleErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) setRolePermissions(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req rolesvc.SetPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := s.role.SetPermissions(c.Request.Context(), id, req); err != nil {
		s.roleErr(c, err)
		return
	}
	response.OK(c, nil)
}

// ---- 权限 ----

func (s *HTTPServer) listPerms(c *gin.Context) {
	page, size := pageParams(c)
	// page_size=0：全量返回（菜单管理树/角色分配权限树需要整表构建，分页会静默截断）
	if c.Query("page_size") == "0" {
		size = 0
	}
	q := bizperm.Query{Type: c.Query("type")}
	ps, pg, err := s.perm.List(c.Request.Context(), q, page, size)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, listResult{List: ps, Page: pg})
}

func (s *HTTPServer) createPerm(c *gin.Context) {
	var req permsvc.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	p, err := s.perm.Create(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, p)
}

func (s *HTTPServer) getPerm(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	p, err := s.perm.Get(c.Request.Context(), id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, p)
}

func (s *HTTPServer) updatePerm(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var req permsvc.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := s.perm.Update(c.Request.Context(), id, req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) deletePerm(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := s.perm.Delete(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, nil)
}

// ---- helpers ----

// ---- 在线用户 ----

func (s *HTTPServer) listOnlineUsers(c *gin.Context) {
	page, size := pageParams(c)
	q := bizsession.Query{Username: c.Query("username"), Device: c.Query("device")}
	currentSid := ""
	if sub := middleware.Subject(c); sub != nil {
		currentSid = sub.SessionID
	}
	vos, pg, err := s.session.List(c.Request.Context(), q, page, size, currentSid)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, listResult{List: vos, Page: pg})
}

func (s *HTTPServer) kickOnlineSession(c *gin.Context) {
	sid := c.Param("sid")
	if sid == "" {
		response.BadRequest(c, "invalid sid")
		return
	}
	operatorID := uint(0)
	if sub := middleware.Subject(c); sub != nil {
		operatorID = sub.UserID
	}
	if err := s.session.Kick(c.Request.Context(), sid, operatorID); err != nil {
		s.sessionErr(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *HTTPServer) kickUserSessions(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	operatorID := uint(0)
	if sub := middleware.Subject(c); sub != nil {
		operatorID = sub.UserID
	}
	if _, err := s.session.KickUser(c.Request.Context(), id, operatorID); err != nil {
		s.sessionErr(c, err)
		return
	}
	response.OK(c, nil)
}

// sessionErr 会话操作错误映射：超管保护返回 403，会话不存在返回 404，其余返回 400
func (s *HTTPServer) sessionErr(c *gin.Context, err error) {
	if errors.Is(err, user.ErrSuperAdminProtected) {
		response.Forbidden(c, err.Error())
		return
	}
	if errors.Is(err, bizsession.ErrSessionNotFound) {
		response.NotFound(c, "会话不存在或已下线")
		return
	}
	response.BadRequest(c, err.Error())
}

// userErr 用户操作错误映射：超管保护类返回 403，其余返回 400
func (s *HTTPServer) userErr(c *gin.Context, err error) {
	if errors.Is(err, user.ErrSuperAdminProtected) || errors.Is(err, user.ErrDeleteSuperAdmin) {
		response.Forbidden(c, err.Error())
		return
	}
	response.BadRequest(c, err.Error())
}

// roleErr 角色操作错误映射：超管角色保护类返回 403，其余返回 400
func (s *HTTPServer) roleErr(c *gin.Context, err error) {
	if errors.Is(err, role.ErrSuperRoleLocked) {
		response.Forbidden(c, err.Error())
		return
	}
	response.BadRequest(c, err.Error())
}

func pageParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	return page, size
}

func idParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "invalid id")
		return 0, false
	}
	return uint(id), true
}

// truncate 按字节长度截断（UA 摘要存储用）
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
