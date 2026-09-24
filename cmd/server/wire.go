//go:build wireinject

// wire 依赖注入装配（与 Kratos 相同的 DI 方式，切换时 Provider 图可直接复用）
package main

import (
	"github.com/google/wire"
	base64Captcha "github.com/mojocn/base64Captcha"
	bizagent "github.com/smilex/smilex-admin-gin/internal/biz/agent"
	bizappuser "github.com/smilex/smilex-admin-gin/internal/biz/appuser"
	"github.com/smilex/smilex-admin-gin/internal/biz/auth"
	bizblacklist "github.com/smilex/smilex-admin-gin/internal/biz/blacklist"
	bizcaptcha "github.com/smilex/smilex-admin-gin/internal/biz/captcha"
	bizdash "github.com/smilex/smilex-admin-gin/internal/biz/dashboard"
	bizdict "github.com/smilex/smilex-admin-gin/internal/biz/dict"
	bizexport "github.com/smilex/smilex-admin-gin/internal/biz/export"
	bizfile "github.com/smilex/smilex-admin-gin/internal/biz/file"
	bizjob "github.com/smilex/smilex-admin-gin/internal/biz/job"
	bizlog "github.com/smilex/smilex-admin-gin/internal/biz/log"
	bizmerchant "github.com/smilex/smilex-admin-gin/internal/biz/merchant"
	bizmonitor "github.com/smilex/smilex-admin-gin/internal/biz/monitor"
	biznotice "github.com/smilex/smilex-admin-gin/internal/biz/notice"
	biznotify "github.com/smilex/smilex-admin-gin/internal/biz/notify"
	bizperm "github.com/smilex/smilex-admin-gin/internal/biz/permission"
	bizrole "github.com/smilex/smilex-admin-gin/internal/biz/role"
	bizsession "github.com/smilex/smilex-admin-gin/internal/biz/session"
	bizsys "github.com/smilex/smilex-admin-gin/internal/biz/sysconfig"
	biztenant "github.com/smilex/smilex-admin-gin/internal/biz/tenant"
	bizuser "github.com/smilex/smilex-admin-gin/internal/biz/user"
	"github.com/smilex/smilex-admin-gin/internal/data"
	dataagent "github.com/smilex/smilex-admin-gin/internal/data/agent"
	dataappuser "github.com/smilex/smilex-admin-gin/internal/data/appuser"
	datablacklist "github.com/smilex/smilex-admin-gin/internal/data/blacklist"
	datacaptcha "github.com/smilex/smilex-admin-gin/internal/data/captcha"
	datadash "github.com/smilex/smilex-admin-gin/internal/data/dashboard"
	datadict "github.com/smilex/smilex-admin-gin/internal/data/dict"
	dataexport "github.com/smilex/smilex-admin-gin/internal/data/export"
	datafile "github.com/smilex/smilex-admin-gin/internal/data/file"
	datajob "github.com/smilex/smilex-admin-gin/internal/data/job"
	datalog "github.com/smilex/smilex-admin-gin/internal/data/log"
	datamerchant "github.com/smilex/smilex-admin-gin/internal/data/merchant"
	datamonitor "github.com/smilex/smilex-admin-gin/internal/data/monitor"
	datanotice "github.com/smilex/smilex-admin-gin/internal/data/notice"
	datanotify "github.com/smilex/smilex-admin-gin/internal/data/notify"
	dataperm "github.com/smilex/smilex-admin-gin/internal/data/permission"
	datarole "github.com/smilex/smilex-admin-gin/internal/data/role"
	datasession "github.com/smilex/smilex-admin-gin/internal/data/session"
	datasys "github.com/smilex/smilex-admin-gin/internal/data/sysconfig"
	datatenant "github.com/smilex/smilex-admin-gin/internal/data/tenant"
	datauser "github.com/smilex/smilex-admin-gin/internal/data/user"
	"github.com/smilex/smilex-admin-gin/internal/server"
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
	notifysvc "github.com/smilex/smilex-admin-gin/internal/service/notify"
	permsvc "github.com/smilex/smilex-admin-gin/internal/service/permission"
	rolesvc "github.com/smilex/smilex-admin-gin/internal/service/role"
	sessionsvc "github.com/smilex/smilex-admin-gin/internal/service/session"
	syssvc "github.com/smilex/smilex-admin-gin/internal/service/sysconfig"
	tenantsvc "github.com/smilex/smilex-admin-gin/internal/service/tenant"
	usersvc "github.com/smilex/smilex-admin-gin/internal/service/user"
)

var bizSet = wire.NewSet(
	bizuser.NewUsecase,
	bizrole.NewUsecase,
	bizperm.NewUsecase,
	bizcaptcha.NewUsecase,
	bizsession.NewUsecase,
	bizlog.NewUsecase,
	bizfile.NewUsecase,
	bizblacklist.NewUsecase,
	bizmerchant.NewUsecase,
	biztenant.NewUsecase,
	bizappuser.NewUsecase,
	bizdict.NewUsecase,
	bizdash.NewUsecase,
	bizsys.NewUsecase,
	biznotice.NewUsecase,
	bizjob.NewUsecase,
	bizmonitor.NewUsecase,
	bizagent.NewUsecase,
	biznotify.NewUsecase,
	bizexport.NewUsecase,
	bizexport.NewRegistry,
	bizexport.NewUserExporter,
	bizexport.NewLoginLogExporter,
	bizexport.NewOpLogExporter,
	auth.NewUsecase,
	// 跨上下文最小依赖接口绑定（provider 与 bind 需同 set）
	wire.Bind(new(auth.CaptchaVerifier), new(*bizcaptcha.Usecase)),
	wire.Bind(new(auth.SessionManager), new(*bizsession.Usecase)),
	wire.Bind(new(bizuser.SessionRevoker), new(*bizsession.Usecase)),
	wire.Bind(new(bizexport.PermissionChecker), new(*auth.Usecase)),
	wire.Bind(new(bizagent.ServerStatusReader), new(*bizmonitor.Usecase)),
	wire.Bind(new(biznotify.SnapshotReader), new(*bizmonitor.Usecase)),
	wire.Bind(new(bizjob.NotifyCleaner), new(*biznotify.Usecase)),
)

var dataRepoSet = wire.NewSet(
	data.NewData,
	data.NewRedisClient,
	data.NewJWTIssuer,
	data.NewAppTokenIssuer,
	datauser.NewRepo,
	datarole.NewRepo,
	dataperm.NewRepo,
	datasession.NewRepo,
	datalog.NewRepo,
	datafile.NewRepo,
	datafile.NewStorageManager,
	datablacklist.NewRepo,
	datamerchant.NewRepo,
	datamerchant.NewAPILogRepo,
	datatenant.NewRepo,
	dataappuser.NewRepo,
	dataagent.NewRepo,
	datadict.NewRepo,
	datadash.NewRepo,
	datamonitor.NewSnapshotRepo,
	wire.Bind(new(bizmonitor.SnapshotRepo), new(*datamonitor.SnapshotRepo)),
	datasys.NewRepo,
	datanotice.NewRepo,
	datanotify.NewRepo,
	datajob.NewRepo,
	datacaptcha.NewStore,
	dataexport.NewRepo,
	dataexport.NewWorker,
	// 跨上下文最小依赖接口绑定
	wire.Bind(new(base64Captcha.Store), new(*datacaptcha.Store)),
	wire.Bind(new(bizjob.LogCleaner), new(*datalog.Repo)),
	wire.Bind(new(bizjob.ExportCleaner), new(*dataexport.Worker)),
	wire.Bind(new(bizmerchant.APILogRepo), new(*datamerchant.APILogRepo)),
	wire.Bind(new(bizjob.MerchantLogCleaner), new(*datamerchant.APILogRepo)),
	wire.Bind(new(bizagent.Repo), new(*dataagent.Repo)),
	wire.Bind(new(bizjob.UsageCleaner), new(*dataagent.Repo)),
	wire.Bind(new(auth.UserStore), new(bizuser.Repo)),
	wire.Bind(new(auth.RoleNameReader), new(bizrole.Repo)),
	wire.Bind(new(auth.PermissionReader), new(bizperm.Repo)),
	wire.Bind(new(bizrole.PermissionReader), new(bizperm.Repo)),
	wire.Bind(new(bizlog.Repo), new(*datalog.Repo)),
	wire.Bind(new(bizblacklist.Repo), new(*datablacklist.Repo)),
	wire.Bind(new(bizblacklist.LoginProtector), new(*datablacklist.Repo)),
	wire.Bind(new(bizexport.Enqueuer), new(*dataexport.Worker)),
)

var serviceSet = wire.NewSet(
	authsvc.NewService,
	usersvc.NewService,
	rolesvc.NewService,
	permsvc.NewService,
	sessionsvc.NewService,
	logsvc.NewService,
	filesvc.NewService,
	blacklistsvc.NewService,
	exportsvc.NewService,
	merchantsvc.NewService,
	tenantsvc.NewService,
	appusersvc.NewService,
	monitorsvc.NewService,
	agentsvc.NewService,
	dictsvc.NewService,
	dashsvc.NewService,
	syssvc.NewService,
	noticesvc.NewService,
	notifysvc.NewService,
	jobsvc.NewService,
)

var providerSet = wire.NewSet(bizSet, dataRepoSet, serviceSet, ProvideConfig, server.NewHTTPServer)

// wireApp 由 wire 生成
func wireApp() (*server.HTTPServer, func(), error) {
	wire.Build(providerSet)
	return nil, nil, nil
}
