// 安全关键中间件测试：CORS 白名单、JWT 认证（会话吊销/typ 隔离）、RBAC 鉴权（缓存）、通用限流。
// Redis 依赖 miniredis 内存实现；令牌签发用真实 data 层 JWT 签发器保证签名/typ 行为一致。
package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	bizauth "github.com/smilex/smilex-admin-gin/internal/biz/auth"
	"github.com/smilex/smilex-admin-gin/internal/biz/permission"
	"github.com/smilex/smilex-admin-gin/internal/biz/session"
	"github.com/smilex/smilex-admin-gin/internal/conf"
	"github.com/smilex/smilex-admin-gin/internal/data"
	authsvc "github.com/smilex/smilex-admin-gin/internal/service/auth"
	"github.com/smilex/smilex-admin-gin/pkg/cache"
)

const testSecret = "test-secret-0123456789abcdef0123456789" // 32 字节，过强度校验

// ---- 领域接口测试替身 ----

// fakeSessions 会话管理替身：Validate 按 map 判定存活
type fakeSessions struct{ valid map[string]bool }

func (f *fakeSessions) Create(ctx context.Context, userID uint, username, nickname, device, ip, userAgent string) (*session.Session, error) {
	return nil, nil
}
func (f *fakeSessions) Validate(ctx context.Context, sid string) bool { return f.valid[sid] }
func (f *fakeSessions) Extend(ctx context.Context, sid string) error  { return nil }
func (f *fakeSessions) Revoke(ctx context.Context, sid string) error  { return nil }
func (f *fakeSessions) RevokeByUser(ctx context.Context, userID uint) (int, error) {
	return 0, nil
}
func (f *fakeSessions) RevokeByUserExcept(ctx context.Context, userID uint, keepSid string) (int, error) {
	return 0, nil
}

// fakePerms 权限读取替身：记录回源次数（验证缓存命中）
type fakePerms struct {
	perms []*permission.Permission
	calls int
}

func (f *fakePerms) FindByUserID(ctx context.Context, userID uint) ([]*permission.Permission, error) {
	f.calls++
	return f.perms, nil
}

// newTestAuthSvc 构造挂载替身的 auth 应用服务
func newTestAuthSvc(perms *fakePerms, sess *fakeSessions) *authsvc.Service {
	issuer := data.NewJWTIssuer(&conf.Bootstrap{JWT: conf.JWT{
		Secret: testSecret, Issuer: "test", ExpireHours: 2, RefreshHours: 168,
	}})
	uc := bizauth.NewUsecase(nil, nil, perms, issuer, nil, sess)
	return authsvc.NewService(uc, nil)
}

// newTestRedis 启动 miniredis 并返回客户端
func newTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })
	return rdb, mr
}

// do 执行一次请求并返回响应
func do(r http.Handler, method, target, bearer, origin string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, nil)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	r.ServeHTTP(w, req)
	return w
}

func okHandler(c *gin.Context) { c.Status(http.StatusOK) }

// ---- CORS ----

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name      string
		origins   []string
		reqOrigin string
		wantACAO  string // 期望 Access-Control-Allow-Origin；"" 表示不应下发
		wantCred  string
	}{
		{"未配置白名单不下发任何跨域头", nil, "https://evil.com", "", ""},
		{"白名单命中回显来源并允许凭证", []string{"https://a.com"}, "https://a.com", "https://a.com", "true"},
		{"白名单未命中不下发跨域头", []string{"https://a.com"}, "https://evil.com", "", ""},
		{"显式通配保持历史行为", []string{"*"}, "https://any.com", "*", ""},
		{"大小写不敏感匹配", []string{"https://a.com"}, "https://A.COM", "https://a.com", "true"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := gin.New()
			e.Use(CORS(tc.origins))
			e.GET("/ping", okHandler)
			w := do(e, "GET", "/ping", "", tc.reqOrigin)
			if got := w.Header().Get("Access-Control-Allow-Origin"); got != tc.wantACAO {
				t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, tc.wantACAO)
			}
			if got := w.Header().Get("Access-Control-Allow-Credentials"); got != tc.wantCred {
				t.Fatalf("Access-Control-Allow-Credentials = %q, want %q", got, tc.wantCred)
			}
		})
	}
}

func TestCORSPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	e := gin.New()
	e.Use(CORS([]string{"https://a.com"}))
	e.GET("/ping", okHandler)
	w := do(e, "OPTIONS", "/ping", "", "https://a.com")
	if w.Code != http.StatusNoContent {
		t.Fatalf("preflight status = %d, want 204", w.Code)
	}
}

// ---- JWT ----

func TestJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := data.NewJWTIssuer(&conf.Bootstrap{JWT: conf.JWT{
		Secret: testSecret, Issuer: "test", ExpireHours: 2, RefreshHours: 168,
	}})
	sub := bizauth.Subject{UserID: 7, Username: "alice", SessionID: "sid-ok"}
	access, _, err := issuer.IssueAccessToken(sub)
	if err != nil {
		t.Fatal(err)
	}
	refresh, err := issuer.IssueRefreshToken(sub)
	if err != nil {
		t.Fatal(err)
	}
	// 同强度异密钥签发器：构造签名不匹配的 token
	other := data.NewJWTIssuer(&conf.Bootstrap{JWT: conf.JWT{
		Secret: strings.Repeat("x", 32), Issuer: "test", ExpireHours: 2, RefreshHours: 168,
	}})
	forged, _, err := other.IssueAccessToken(sub)
	if err != nil {
		t.Fatal(err)
	}
	// 过期 token：手工签发同密钥、过期时间为 1 小时前的 access token
	claims := jwt.MapClaims{
		"uid": 7, "username": "alice", "sid": "sid-ok", "typ": "access",
		"iss": "test", "exp": time.Now().Add(-time.Hour).Unix(),
	}
	expired, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}

	svc := newTestAuthSvc(&fakePerms{}, &fakeSessions{valid: map[string]bool{"sid-ok": true}})
	e := gin.New()
	e.Use(JWT(svc))
	e.GET("/ping", func(c *gin.Context) {
		s := Subject(c)
		if s == nil || s.UserID != 7 || s.SessionID != "sid-ok" {
			t.Errorf("subject not propagated: %+v", s)
		}
		c.Status(http.StatusOK)
	})

	if w := do(e, "GET", "/ping", access, ""); w.Code != http.StatusOK {
		t.Fatalf("valid token: status = %d, want 200 (body %s)", w.Code, w.Body)
	}
	for _, tc := range []struct {
		name   string
		bearer string
	}{
		{"缺少 Bearer", ""},
		{"乱码 token", "not-a-jwt"},
		{"签名不匹配", forged},
		{"refresh 不能当 access 用", refresh},
		{"过期 token", expired},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if w := do(e, "GET", "/ping", tc.bearer, ""); w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401", w.Code)
			}
		})
	}
}

func TestJWTSessionRevoked(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := data.NewJWTIssuer(&conf.Bootstrap{JWT: conf.JWT{
		Secret: testSecret, Issuer: "test", ExpireHours: 2, RefreshHours: 168,
	}})
	access, _, err := issuer.IssueAccessToken(bizauth.Subject{UserID: 7, Username: "alice", SessionID: "sid-gone"})
	if err != nil {
		t.Fatal(err)
	}
	// 会话已被吊销（Validate=false）：token 未过期也必须拒绝（登出/踢人即时生效的根基）
	svc := newTestAuthSvc(&fakePerms{}, &fakeSessions{valid: map[string]bool{}})
	e := gin.New()
	e.Use(JWT(svc))
	e.GET("/ping", okHandler)
	if w := do(e, "GET", "/ping", access, ""); w.Code != http.StatusUnauthorized {
		t.Fatalf("revoked session: status = %d, want 401", w.Code)
	}
}

// ---- RBAC ----

func TestRBAC(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := data.NewJWTIssuer(&conf.Bootstrap{JWT: conf.JWT{
		Secret: testSecret, Issuer: "test", ExpireHours: 2, RefreshHours: 168,
	}})
	perms := &fakePerms{perms: []*permission.Permission{
		{ID: 1, Type: permission.TypeButton, Method: "GET", Path: "/api/v1/users"},
		{ID: 2, Type: permission.TypeButton, Method: "*", Path: "/api/v1/exports/*"},
	}}
	svc := newTestAuthSvc(perms, &fakeSessions{valid: map[string]bool{"sid-ok": true}})
	rdb, _ := newTestRedis(t)
	twoLevel := cache.NewTwoLevel(rdb, "rbac:", 30*time.Second, 60*time.Second, true)

	access, _, err := issuer.IssueAccessToken(bizauth.Subject{UserID: 7, Username: "alice", SessionID: "sid-ok"})
	if err != nil {
		t.Fatal(err)
	}

	e := gin.New()
	e.Use(func(c *gin.Context) {
		c.Set(ctxSubjectKey, &bizauth.Subject{UserID: 7, Username: "alice", SessionID: "sid-ok"})
		c.Next()
	})
	e.Use(RBAC(svc, twoLevel))
	e.GET("/api/v1/users", okHandler)
	e.GET("/api/v1/exports/abc", okHandler)
	e.GET("/api/v1/roles", okHandler)

	// 未认证（无 Subject）：RBAC 直接 401 而非 403
	e2 := gin.New()
	e2.Use(RBAC(svc, twoLevel))
	e2.GET("/api/v1/users", okHandler)
	if w := do(e2, "GET", "/api/v1/users", "", ""); w.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated: status = %d, want 401", w.Code)
	}

	// 精确命中
	if w := do(e, "GET", "/api/v1/users", access, ""); w.Code != http.StatusOK {
		t.Fatalf("exact match: status = %d, want 200", w.Code)
	}
	// glob 前缀通配命中
	if w := do(e, "GET", "/api/v1/exports/abc", access, ""); w.Code != http.StatusOK {
		t.Fatalf("glob match: status = %d, want 200", w.Code)
	}
	// 未授权路径：默认拒绝
	if w := do(e, "GET", "/api/v1/roles", access, ""); w.Code != http.StatusForbidden {
		t.Fatalf("default deny: status = %d, want 403", w.Code)
	}
	// 缓存命中：三条不同路由各自回源一次（users/exports/roles），重复请求不再回源
	if perms.calls != 3 {
		t.Fatalf("perm loader calls = %d, want 3 (one per unique route key)", perms.calls)
	}
	do(e, "GET", "/api/v1/users", access, "")
	if perms.calls != 3 {
		t.Fatalf("perm loader calls after repeat = %d, want 3 (cache hit)", perms.calls)
	}
}

// ---- 限流 ----

func TestRateLimitByIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rdb, mr := newTestRedis(t)
	e := gin.New()
	e.Use(NewRateLimit(rdb, RateLimitConfig{
		KeyPrefix: "rl:t:", Max: 2, Window: time.Minute, MessageKey: "security.rate_limited",
	}))
	e.GET("/ping", okHandler)

	for i := 1; i <= 2; i++ {
		if w := do(e, "GET", "/ping", "", ""); w.Code != http.StatusOK {
			t.Fatalf("request #%d: status = %d, want 200", i, w.Code)
		}
	}
	w := do(e, "GET", "/ping", "", "")
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("request #3: status = %d, want 429", w.Code)
	}
	if w.Header().Get("Retry-After") == "" {
		t.Fatal("missing Retry-After header")
	}
	// 窗口滑过（快进 TTL）后恢复放行
	mr.FastForward(time.Minute)
	if w := do(e, "GET", "/ping", "", ""); w.Code != http.StatusOK {
		t.Fatalf("after window reset: status = %d, want 200", w.Code)
	}
}

func TestRateLimitByUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rdb, _ := newTestRedis(t)
	e := gin.New()
	e.Use(func(c *gin.Context) {
		uid := c.Query("uid")
		if uid != "" {
			n, _ := strconv.ParseUint(uid, 10, 64)
			c.Set(ctxSubjectKey, &bizauth.Subject{UserID: uint(n), Username: "u" + uid, SessionID: "s"})
		}
		c.Next()
	})
	e.Use(NewRateLimit(rdb, RateLimitConfig{
		KeyPrefix: "rl:u:", Max: 1, Window: time.Minute, ByUser: true, MessageKey: "security.rate_limited",
	}))
	e.GET("/ping", okHandler)

	// 用户 1 首次放行、第二次 429；用户 2 不受用户 1 影响
	if w := do(e, "GET", "/ping?uid=1", "", ""); w.Code != http.StatusOK {
		t.Fatalf("user1 #1: status = %d, want 200", w.Code)
	}
	if w := do(e, "GET", "/ping?uid=1", "", ""); w.Code != http.StatusTooManyRequests {
		t.Fatalf("user1 #2: status = %d, want 429", w.Code)
	}
	if w := do(e, "GET", "/ping?uid=2", "", ""); w.Code != http.StatusOK {
		t.Fatalf("user2 #1: status = %d, want 200", w.Code)
	}
	// 按用户限流但请求未认证：直接 401（不落入按 IP 计数）
	if w := do(e, "GET", "/ping", "", ""); w.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated: status = %d, want 401", w.Code)
	}
}

func TestRateLimitRedisFailOpen(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rdb, mr := newTestRedis(t)
	e := gin.New()
	e.Use(NewRateLimit(rdb, RateLimitConfig{
		KeyPrefix: "rl:t:", Max: 1, Window: time.Minute, MessageKey: "security.rate_limited",
	}))
	e.GET("/ping", okHandler)

	// Redis 故障：限流是防御层，fail-open 放行（认证可用性由 JWT 会话校验 fail-closed 兜底）
	mr.SetError("redis down")
	for i := 0; i < 3; i++ {
		if w := do(e, "GET", "/ping", "", ""); w.Code != http.StatusOK {
			t.Fatalf("request #%d under redis failure: status = %d, want 200 (fail-open)", i+1, w.Code)
		}
	}
}
