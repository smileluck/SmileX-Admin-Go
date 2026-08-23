//go:build wireinject

// wire 依赖注入装配（与 Kratos 相同的 DI 方式，切换时 Provider 图可直接复用）
package main

import (
	"github.com/google/wire"
	"github.com/smilex/smilex-admin-gin/internal/biz/auth"
	bizcaptcha "github.com/smilex/smilex-admin-gin/internal/biz/captcha"
	bizexport "github.com/smilex/smilex-admin-gin/internal/biz/export"
	bizfile "github.com/smilex/smilex-admin-gin/internal/biz/file"
	bizlog "github.com/smilex/smilex-admin-gin/internal/biz/log"
	bizperm "github.com/smilex/smilex-admin-gin/internal/biz/permission"
	bizrole "github.com/smilex/smilex-admin-gin/internal/biz/role"
	bizsession "github.com/smilex/smilex-admin-gin/internal/biz/session"
	bizuser "github.com/smilex/smilex-admin-gin/internal/biz/user"
	"github.com/smilex/smilex-admin-gin/internal/data"
	dataexport "github.com/smilex/smilex-admin-gin/internal/data/export"
	datafile "github.com/smilex/smilex-admin-gin/internal/data/file"
	datalog "github.com/smilex/smilex-admin-gin/internal/data/log"
	dataperm "github.com/smilex/smilex-admin-gin/internal/data/permission"
	datarole "github.com/smilex/smilex-admin-gin/internal/data/role"
	datasession "github.com/smilex/smilex-admin-gin/internal/data/session"
	datauser "github.com/smilex/smilex-admin-gin/internal/data/user"
	"github.com/smilex/smilex-admin-gin/internal/server"
	authsvc "github.com/smilex/smilex-admin-gin/internal/service/auth"
	exportsvc "github.com/smilex/smilex-admin-gin/internal/service/export"
	filesvc "github.com/smilex/smilex-admin-gin/internal/service/file"
	logsvc "github.com/smilex/smilex-admin-gin/internal/service/log"
	permsvc "github.com/smilex/smilex-admin-gin/internal/service/permission"
	rolesvc "github.com/smilex/smilex-admin-gin/internal/service/role"
	sessionsvc "github.com/smilex/smilex-admin-gin/internal/service/session"
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
)

var dataRepoSet = wire.NewSet(
	data.NewData,
	data.NewRedisClient,
	data.NewJWTIssuer,
	datauser.NewRepo,
	datarole.NewRepo,
	dataperm.NewRepo,
	datasession.NewRepo,
	datalog.NewRepo,
	datafile.NewRepo,
	datafile.NewStorageManager,
	dataexport.NewRepo,
	dataexport.NewWorker,
	// 跨上下文最小依赖接口绑定
	wire.Bind(new(auth.UserStore), new(bizuser.Repo)),
	wire.Bind(new(auth.RoleNameReader), new(bizrole.Repo)),
	wire.Bind(new(auth.PermissionReader), new(bizperm.Repo)),
	wire.Bind(new(bizlog.Repo), new(*datalog.Repo)),
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
	exportsvc.NewService,
)

var providerSet = wire.NewSet(bizSet, dataRepoSet, serviceSet, ProvideConfig, server.NewHTTPServer)

// wireApp 由 wire 生成
func wireApp() (*server.HTTPServer, func(), error) {
	wire.Build(providerSet)
	return nil, nil, nil
}
