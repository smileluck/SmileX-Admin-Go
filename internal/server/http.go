package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	bizappuser "github.com/smilex/smilex-admin-gin/internal/biz/appuser"
	bizmerchant "github.com/smilex/smilex-admin-gin/internal/biz/merchant"
	"github.com/smilex/smilex-admin-gin/internal/conf"
	"github.com/smilex/smilex-admin-gin/internal/server/middleware"
	agentsvc "github.com/smilex/smilex-admin-gin/internal/service/agent"
	appusersvc "github.com/smilex/smilex-admin-gin/internal/service/appuser"
	authsvc "github.com/smilex/smilex-admin-gin/internal/service/auth"
	blacklistsvc "github.com/smilex/smilex-admin-gin/internal/service/blacklist"
	dashsvc "github.com/smilex/smilex-admin-gin/internal/service/dashboard"
	dictsvc "github.com/smilex/smilex-admin-gin/internal/service/dict"
	exportsvc "github.com/smilex/smilex-admin-gin/internal/service/export"
	filesvc "github.com/smilex/smilex-admin-gin/internal/service/file"
	jobsvc "github.com/smilex/smilex-admin-gin/internal/service/job"
	logsvc "github.com/smilex/smilex-admin-gin/internal/service/log"
	merchantsvc "github.com/smilex/smilex-admin-gin/internal/service/merchant"
	monitorsvc "github.com/smilex/smilex-admin-gin/internal/service/monitor"
	noticesvc "github.com/smilex/smilex-admin-gin/internal/service/notice"
	permsvc "github.com/smilex/smilex-admin-gin/internal/service/permission"
	rolesvc "github.com/smilex/smilex-admin-gin/internal/service/role"
	sessionsvc "github.com/smilex/smilex-admin-gin/internal/service/session"
	syssvc "github.com/smilex/smilex-admin-gin/internal/service/sysconfig"
	tenantsvc "github.com/smilex/smilex-admin-gin/internal/service/tenant"
	usersvc "github.com/smilex/smilex-admin-gin/internal/service/user"
	"github.com/smilex/smilex-admin-gin/pkg/cache"
	"github.com/smilex/smilex-admin-gin/pkg/logger"
	"github.com/smilex/smilex-admin-gin/pkg/response"
	"go.uber.org/zap"
)

// HTTPServer 聚合全部应用服务
type HTTPServer struct {
	cfg        *conf.Bootstrap
	auth       *authsvc.Service
	user       *usersvc.Service
	role       *rolesvc.Service
	perm       *permsvc.Service
	session    *sessionsvc.Service
	log        *logsvc.Service
	file       *filesvc.Service
	export     *exportsvc.Service
	blacklist  *blacklistsvc.Service
	merchant   *merchantsvc.Service
	merchantUC *bizmerchant.Usecase // 开放 API 验签中间件直连领域用例
	tenant     *tenantsvc.Service
	appuser    *appusersvc.Service
	monitor    *monitorsvc.Service
	agent      *agentsvc.Service
	dict       *dictsvc.Service
	syscfg     *syssvc.Service
	notice     *noticesvc.Service
	job        *jobsvc.Service
	dashboard  *dashsvc.Service
	appuserUC  *bizappuser.Usecase    // AppJWT 中间件直连领域用例（校验用户启用状态）
	appIssuer  bizappuser.TokenIssuer // AppJWT 中间件解析 app-access token
	rbacCache  *cache.TwoLevel
	rdb        *redis.Client // nonce 防重放（开放 API 验签）
	engine     *gin.Engine
	srv        *http.Server
}

// NewHTTPServer 构造并注册路由
func NewHTTPServer(cfg *conf.Bootstrap, auth *authsvc.Service, user *usersvc.Service,
	role *rolesvc.Service, perm *permsvc.Service, session *sessionsvc.Service, log *logsvc.Service,
	file *filesvc.Service, export *exportsvc.Service, blacklist *blacklistsvc.Service,
	merchant *merchantsvc.Service, merchantUC *bizmerchant.Usecase,
	tenant *tenantsvc.Service, appuser *appusersvc.Service, appuserUC *bizappuser.Usecase,
	appIssuer bizappuser.TokenIssuer, monitor *monitorsvc.Service, agent *agentsvc.Service,
	dict *dictsvc.Service, syscfg *syssvc.Service, notice *noticesvc.Service,
	job *jobsvc.Service, dash *dashsvc.Service, rdb *redis.Client) *HTTPServer {
	gin.SetMode(cfg.Server.Mode)
	e := gin.New()
	// ClientIP 只信可信代理链上的 X-Forwarded-For：默认空表=不信任任何代理（取直连地址），
	// 防止伪造请求头绕过登录限流/封禁或污染审计 IP；反代部署在 configs 配置代理机地址
	if err := e.SetTrustedProxies(cfg.Server.TrustedProxies); err != nil {
		logger.Warn("set trusted proxies failed", zap.Strings("proxies", cfg.Server.TrustedProxies), zap.Error(err))
	}
	// multipart 表单内存上限保持较小值（超出部分落临时文件）；上传大小由 handler 显式校验
	e.MaxMultipartMemory = 8 << 20
	e.Use(gin.Recovery(),
		middleware.I18n(),
		middleware.SecurityHeaders(),
		middleware.CORS(cfg.Server.CORSOrigins),
		middleware.XSSFilter(),
		middleware.SQLInjectionGuard(),
	)

	// RBAC 权限判定缓存：L1 30s 进程内存 + L2 60s Redis（cache.l2Enabled 可关）
	rbacCache := cache.NewTwoLevel(rdb, "rbac:", 30*time.Second, 60*time.Second, cfg.Cache.L2Enabled)

	s := &HTTPServer{cfg: cfg, auth: auth, user: user, role: role, perm: perm, session: session, log: log, file: file, export: export, blacklist: blacklist, merchant: merchant, merchantUC: merchantUC, tenant: tenant, appuser: appuser, appuserUC: appuserUC, appIssuer: appIssuer, monitor: monitor, agent: agent, dict: dict, syscfg: syscfg, notice: notice, job: job, dashboard: dash, rbacCache: rbacCache, rdb: rdb, engine: e}
	// 内置系统参数幂等播种（不覆盖用户修改；失败不阻断启动）
	if err := syscfg.EnsureBuiltin(); err != nil {
		logger.Warn("ensure builtin sys-configs failed", zap.Error(err))
	}
	// 定时任务：播种内置清理任务并拉起调度器（失败不阻断启动）
	if err := job.EnsureSeededAndStart(); err != nil {
		logger.Warn("job scheduler start failed", zap.Error(err))
	}

	s.registerRoutes()
	s.registerStatic()

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: e,
	}
	return s
}

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
	s.job.Stop() // 停定时调度器（等待在跑任务完成）
	return s.srv.Shutdown(ctx)
}
