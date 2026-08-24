// Package middleware HTTP 中间件：JWT 认证、RBAC 鉴权、CORS。
package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smilex/smilex-admin-gin/internal/biz/auth"
	authsvc "github.com/smilex/smilex-admin-gin/internal/service/auth"
	"github.com/smilex/smilex-admin-gin/pkg/cache"
	"github.com/smilex/smilex-admin-gin/pkg/i18n"
	"github.com/smilex/smilex-admin-gin/pkg/response"
)

// I18n 按 Accept-Language 头识别请求语言并注入 context（须在全局链最前注册，
// 保证后续中间件与 handler 都能从 context 取到语言）
func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := i18n.Detect(c.GetHeader("Accept-Language"))
		c.Set("locale", string(l))
		c.Request = c.Request.WithContext(i18n.WithLocale(c.Request.Context(), l))
		c.Next()
	}
}

// CORS 跨域
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

const ctxSubjectKey = "auth.subject"

// JWT 认证：Bearer token -> Subject 存入 context。
// 会话校验：token 携带 sid，仅当对应会话存活（未被吊销/顶替/过期）时放行；
// Redis 故障时 fail-closed，会话一律视为无效。
func JWT(authSvc *authsvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		token, ok := strings.CutPrefix(h, "Bearer ")
		if !ok || token == "" {
			response.Unauthorized(c, "missing bearer token")
			c.Abort()
			return
		}
		s, err := authSvc.ParseSubject(token)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}
		if s.SessionID == "" || !authSvc.ValidateSession(c.Request.Context(), s.SessionID) {
			response.Unauthorized(c, "session revoked")
			c.Abort()
			return
		}
		c.Set(ctxSubjectKey, s)
		c.Next()
	}
}

// Subject 从 context 取认证主体（JWT 中间件之后可用）
func Subject(c *gin.Context) *auth.Subject {
	if v, ok := c.Get(ctxSubjectKey); ok {
		if s, ok := v.(*auth.Subject); ok {
			return s
		}
	}
	return nil
}

// RBAC 接口鉴权（二级缓存：L1 进程内存 + L2 Redis，减少查库；一致性由短 TTL 兜底）
func RBAC(authSvc *authsvc.Service, cache *cache.TwoLevel) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := Subject(c)
		if s == nil {
			response.Unauthorized(c, "unauthenticated")
			c.Abort()
			return
		}
		key := fmt.Sprintf("%d|%s|%s", s.UserID, c.Request.Method, c.Request.URL.Path)
		if v, ok := cache.Get(c.Request.Context(), key); ok {
			if v != "1" {
				response.Forbidden(c, "permission denied")
				c.Abort()
			}
			return
		}
		allow := authSvc.Authorize(c.Request.Context(), s.UserID, c.Request.Method, c.Request.URL.Path)
		val := "0"
		if allow {
			val = "1"
		}
		cache.Set(c.Request.Context(), key, val)
		if !allow {
			response.Forbidden(c, "permission denied")
			c.Abort()
		}
	}
}
