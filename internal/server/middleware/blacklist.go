// 登录 IP 黑名单中间件：连续登录失败拉黑该 IP，作为口令爆破的第二层防御
// （第一层为 LoginRateLimit 的 5 次/60s 频率限制；黑名单在限流之前生效）。
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	bizlog "github.com/smilex/smilex-admin-gin/internal/biz/log"
	"github.com/smilex/smilex-admin-gin/pkg/i18n"
	"github.com/smilex/smilex-admin-gin/pkg/response"
)

// 单机内存实现（与 loginLimiter 同风格；多实例部署时应换共享存储实现）
const (
	loginFailThreshold = 5                // 计数窗口内失败次数阈值
	loginFailWindow    = 10 * time.Minute // 失败计数滑动窗口
	loginBlockDuration = 10 * time.Minute // 拉黑时长
)

// failCounter 单 IP 失败计数（窗口内滑动）
type failCounter struct {
	count int
	first time.Time // 窗口内首次失败时间
}

type loginBlacklist struct {
	mu      sync.Mutex
	fails   map[string]*failCounter // ip -> 失败计数
	blocked map[string]time.Time    // ip -> 解封时间
}

// LoginLogRecorder 登录日志记录入口（由 log 应用服务实现，异步落库）
type LoginLogRecorder interface {
	RecordLogin(ctx context.Context, l *bizlog.LoginLog)
}

// LoginIPGuard 登录 IP 黑名单守卫：
//   - 请求前：已拉黑 IP 直接 403（提示剩余等待分钟），并记一条被拦截的登录日志；
//   - 请求后：登录成功（200）清空该 IP 失败计数；失败（401 = 密码错/账号禁用）
//     计数 +1，10 分钟窗口内满 5 次 → 拉黑 10 分钟；
//     400（验证码错误/参数问题）不计入，避免正常用户输错验证码被误伤。
func LoginIPGuard(rec LoginLogRecorder) gin.HandlerFunc {
	b := &loginBlacklist{fails: map[string]*failCounter{}, blocked: map[string]time.Time{}}
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		if until, ok := b.isBlocked(ip, now); ok {
			mins := int(time.Until(until).Minutes()) + 1
			c.Header("Retry-After", strconv.Itoa(mins*60))
			response.Forbidden(c, i18n.T(c.Request.Context(), "blacklist.ip_banned", mins))
			recordBlockedLogin(c, rec, ip)
			c.Abort()
			return
		}

		c.Next()

		switch c.Writer.Status() {
		case http.StatusOK:
			b.reset(ip)
		case http.StatusUnauthorized:
			b.recordFail(ip, now)
		}
	}
}

// isBlocked 是否处于拉黑期（惰性清理过期条目）
func (b *loginBlacklist) isBlocked(ip string, now time.Time) (time.Time, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	until, ok := b.blocked[ip]
	if !ok {
		return time.Time{}, false
	}
	if now.After(until) {
		delete(b.blocked, ip)
		return time.Time{}, false
	}
	return until, true
}

// reset 登录成功清空失败计数
func (b *loginBlacklist) reset(ip string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.fails, ip)
}

// recordFail 记一次失败；达到阈值时拉黑并返回 true
func (b *loginBlacklist) recordFail(ip string, now time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	fc := b.fails[ip]
	if fc == nil || now.Sub(fc.first) >= loginFailWindow {
		fc = &failCounter{first: now}
		b.fails[ip] = fc
	}
	fc.count++
	if fc.count < loginFailThreshold {
		return false
	}
	// 达到阈值：拉黑并清计数（解封后重新计数）
	delete(b.fails, ip)
	b.blocked[ip] = now.Add(loginBlockDuration)
	// IP 数超阈值时清理过期条目，防止长期运行缓慢增长
	if len(b.blocked) > 10000 {
		for k, until := range b.blocked {
			if now.After(until) {
				delete(b.blocked, k)
			}
		}
	}
	if len(b.fails) > 10000 {
		for k, fc := range b.fails {
			if now.Sub(fc.first) >= loginFailWindow {
				delete(b.fails, k)
			}
		}
	}
	return true
}

// recordBlockedLogin 被黑名单拦截的尝试也入登录日志（用户名/设备端尽力从 body 提取）
func recordBlockedLogin(c *gin.Context, rec LoginLogRecorder, ip string) {
	if rec == nil {
		return
	}
	l := &bizlog.LoginLog{
		IP:        ip,
		UserAgent: truncateStr(c.GetHeader("User-Agent"), 255),
		Status:    bizlog.LoginStatusFail,
		Msg:       "IP 已被临时封禁",
		CreatedAt: time.Now(),
	}
	if body := peekBody(c); len(body) > 0 {
		var m map[string]interface{}
		if json.Unmarshal(body, &m) == nil {
			if u, ok := m["username"].(string); ok {
				l.Username = truncateStr(u, 64)
			}
			if d, ok := m["deviceType"].(string); ok {
				l.Device = truncateStr(d, 16)
			}
		}
	}
	rec.RecordLogin(context.Background(), l)
}

// ---- 管理员手工维护的持久化 IP 黑名单 ----

// Checker 持久化黑名单判定接口（由 blacklist 领域用例实现，内部为 30s TTL 内存快照）
type Checker interface {
	IsBlocked(ip string) bool
}

// IPBlacklist 持久化 IP 黑名单中间件：挂在 /api/v1 组上、JWT 之前生效，
// 命中即 403 拦截全部 /api/ 请求（不拦截静态前端资源）；
// 与上方登录失败临时封禁 LoginIPGuard 独立共存。
func IPBlacklist(chk Checker) gin.HandlerFunc {
	return func(c *gin.Context) {
		if chk != nil && chk.IsBlocked(c.ClientIP()) {
			response.Forbidden(c, i18n.T(c.Request.Context(), "blacklist.ip_blocked"))
			c.Abort()
			return
		}
		c.Next()
	}
}
