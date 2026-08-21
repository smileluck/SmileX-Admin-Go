// Package server 传输层：Gin HTTP Server 与路由注册。
// 切换 Kratos 时本层是唯一需要替换的层（由 proto 生成的 HTTP/gRPC server 代替）。
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	bizperm "github.com/smilex/smilex-admin-gin/internal/biz/permission"
	"github.com/smilex/smilex-admin-gin/internal/biz/role"
	"github.com/smilex/smilex-admin-gin/internal/biz/user"
	"github.com/smilex/smilex-admin-gin/internal/conf"
	"github.com/smilex/smilex-admin-gin/internal/server/middleware"
	authsvc "github.com/smilex/smilex-admin-gin/internal/service/auth"
	permsvc "github.com/smilex/smilex-admin-gin/internal/service/permission"
	rolesvc "github.com/smilex/smilex-admin-gin/internal/service/role"
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
	engine  *gin.Engine
	srv     *http.Server
}

// NewHTTPServer 构造并注册路由
func NewHTTPServer(cfg *conf.Bootstrap, auth *authsvc.Service, user *usersvc.Service,
	role *rolesvc.Service, perm *permsvc.Service) *HTTPServer {
	gin.SetMode(cfg.Server.Mode)
	e := gin.New()
	e.Use(gin.Recovery(), middleware.CORS())

	s := &HTTPServer{cfg: cfg, auth: auth, user: user, role: role, perm: perm, engine: e}
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
		authg.POST("/login", func(c *gin.Context) {
			var req authsvc.LoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, err.Error())
				return
			}
			tp, err := s.auth.Login(c.Request.Context(), req)
			if err != nil {
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

	// ---- 受保护接口：JWT -> RBAC ----
	protected := v1.Group("", middleware.JWT(s.auth), middleware.RBAC(s.auth))
	{
		protected.POST("/auth/logout", func(c *gin.Context) { response.OK(c, nil) })
		protected.GET("/auth/profile", func(c *gin.Context) {
			sub := middleware.Subject(c)
			vo, err := s.auth.Profile(c.Request.Context(), sub.UserID)
			if err != nil {
				response.ServerError(c, err.Error())
				return
			}
			response.OK(c, vo)
		})
		// 当前用户可见菜单树
		protected.GET("/menus", func(c *gin.Context) {
			sub := middleware.Subject(c)
			tree, err := s.perm.UserMenuTree(c.Request.Context(), sub.UserID)
			if err != nil {
				response.ServerError(c, err.Error())
				return
			}
			response.OK(c, tree)
		})
	}

	users := protected.Group("/users")
	{
		users.GET("", s.listUsers)
		users.POST("", s.createUser)
		users.GET("/:id", s.getUser)
		users.PUT("/:id", s.updateUser)
		users.DELETE("/:id", s.deleteUser)
		users.PUT("/:id/roles", s.setUserRoles)
		users.PUT("/:id/password", s.resetUserPassword)
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
	List  interface{} `json:"list"`
	Page  interface{} `json:"page"`
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
	if err := s.user.Update(c.Request.Context(), id, req); err != nil {
		response.BadRequest(c, err.Error())
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
		response.BadRequest(c, err.Error())
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
	if err := s.user.SetRoles(c.Request.Context(), id, req); err != nil {
		response.BadRequest(c, err.Error())
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
	if err := s.user.ResetPassword(c.Request.Context(), id, req); err != nil {
		response.BadRequest(c, err.Error())
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
		response.BadRequest(c, err.Error())
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
		response.BadRequest(c, err.Error())
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
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, nil)
}

// ---- 权限 ----

func (s *HTTPServer) listPerms(c *gin.Context) {
	page, size := pageParams(c)
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
