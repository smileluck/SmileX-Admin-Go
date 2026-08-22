// 安全类中间件：安全响应头、XSS 输入清洗、SQL 注入特征拦截、登录限流。
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smilex/smilex-admin-gin/pkg/response"
	"github.com/smilex/smilex-admin-gin/pkg/security"
)

// SecurityHeaders 安全响应头：防 MIME 嗅探、点击劫持与协议泄露。
// CSP 面向 SPA 托管（web/dist）：style unsafe-inline 兼容 naive-ui 的 cssinjs，
// img data: 兼容 base64 验证码、https: 兼容菜单网络图标。
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Content-Security-Policy",
			"default-src 'self'; img-src 'self' data: https:; style-src 'self' 'unsafe-inline'; "+
				"script-src 'self'; connect-src 'self'; font-src 'self' data:; "+
				"frame-ancestors 'none'; base-uri 'self'; form-action 'self'")
		c.Next()
	}
}

// XSSFilter 对 JSON body 的字符串值剥离 HTML（存储型 XSS 纵深防御层）。
// 仅处理 application/json 的写入类请求；解析失败原样放行，由后续 ShouldBindJSON 报错。
func XSSFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			if strings.HasPrefix(c.GetHeader("Content-Type"), "application/json") {
				if body, err := io.ReadAll(c.Request.Body); err == nil {
					c.Request.Body = io.NopCloser(bytes.NewReader(sanitizeJSON(body)))
				}
			}
		}
		c.Next()
	}
}

// sanitizeJSON 递归清洗 JSON 文档中的字符串值；非法 JSON 原样返回。
// 密码类字段（key 含 password）跳过 —— 合法密码允许包含尖括号。
func sanitizeJSON(body []byte) []byte {
	var payload interface{}
	if json.Unmarshal(body, &payload) != nil {
		return body
	}
	if _, ok := payload.(map[string]interface{}); !ok {
		return body
	}
	cleaned, err := json.Marshal(sanitizeValue(payload))
	if err != nil {
		return body
	}
	return cleaned
}

func sanitizeValue(v interface{}) interface{} {
	switch x := v.(type) {
	case string:
		return security.StripHTML(x)
	case map[string]interface{}:
		for k, val := range x {
			if strings.Contains(strings.ToLower(k), "password") {
				continue
			}
			x[k] = sanitizeValue(val)
		}
		return x
	case []interface{}:
		for i, val := range x {
			x[i] = sanitizeValue(val)
		}
		return x
	default:
		return v
	}
}

// SQLInjectionGuard 对 query 与路径参数做高危 SQL 注入特征检测（WAF 式前置拦截）。
// 数据层本就全参数化（GORM 占位符），这里是纵深防御的入口层。
func SQLInjectionGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, vs := range c.Request.URL.Query() {
			for _, v := range vs {
				if security.ContainsSQLInjection(v) {
					response.BadRequest(c, "请求参数包含非法字符")
					c.Abort()
					return
				}
			}
		}
		for _, p := range c.Params {
			if security.ContainsSQLInjection(p.Value) {
				response.BadRequest(c, "请求参数包含非法字符")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// loginLimiter 登录接口内存滑动窗口限流（单机版；多实例部署时应换共享存储实现）
type loginLimiter struct {
	mu      sync.Mutex
	hits    map[string][]time.Time
	window  time.Duration
	maxHits int
}

// LoginRateLimit 按 IP 限流登录尝试（默认 60s 内 5 次），超出返回 429，防口令爆破
func LoginRateLimit() gin.HandlerFunc {
	l := &loginLimiter{hits: map[string][]time.Time{}, window: time.Minute, maxHits: 5}
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		l.mu.Lock()
		hs := l.hits[ip]
		kept := hs[:0]
		for _, t := range hs {
			if now.Sub(t) < l.window {
				kept = append(kept, t)
			}
		}
		if len(kept) >= l.maxHits {
			l.hits[ip] = kept
			l.mu.Unlock()
			c.Header("Retry-After", "60")
			response.TooManyRequests(c, "登录尝试过于频繁，请稍后再试")
			c.Abort()
			return
		}
		l.hits[ip] = append(kept, now)
		// IP 数超阈值时清理窗口外的非活跃条目，防止长期运行缓慢增长
		if len(l.hits) > 10000 {
			for k, v := range l.hits {
				if len(v) == 0 || now.Sub(v[len(v)-1]) >= l.window {
					delete(l.hits, k)
				}
			}
		}
		l.mu.Unlock()
		c.Next()
	}
}
