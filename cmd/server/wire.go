//go:build wireinject

// wire 依赖注入装配（与 Kratos 相同的 DI 方式，切换时 Provider 图可直接复用）
package main

import (
	"github.com/google/wire"
	"github.com/smilex/smilex-admin-gin/internal/biz/auth"
	bizcaptcha "github.com/smilex/smilex-admin-gin/internal/biz/captcha"
	bizperm "github.com/smilex/smilex-admin-gin/internal/biz/permission"
	bizrole "github.com/smilex/smilex-admin-gin/internal/biz/role"
	bizuser "github.com/smilex/smilex-admin-gin/internal/biz/user"
	"github.com/smilex/smilex-admin-gin/internal/data"
	dataperm "github.com/smilex/smilex-admin-gin/internal/data/permission"
	datarole "github.com/smilex/smilex-admin-gin/internal/data/role"
	datauser "github.com/smilex/smilex-admin-gin/internal/data/user"
	"github.com/smilex/smilex-admin-gin/internal/server"
	authsvc "github.com/smilex/smilex-admin-gin/internal/service/auth"
	permsvc "github.com/smilex/smilex-admin-gin/internal/service/permission"
	rolesvc "github.com/smilex/smilex-admin-gin/internal/service/role"
	usersvc "github.com/smilex/smilex-admin-gin/internal/service/user"
)

var bizSet = wire.NewSet(
	bizuser.NewUsecase,
	bizrole.NewUsecase,
	bizperm.NewUsecase,
	bizcaptcha.NewUsecase,
	auth.NewUsecase,
	// 跨上下文最小依赖接口绑定（provider 与 bind 需同 set）
	wire.Bind(new(auth.CaptchaVerifier), new(*bizcaptcha.Usecase)),
)

var dataRepoSet = wire.NewSet(
	data.NewData,
	data.NewJWTIssuer,
	datauser.NewRepo,
	datarole.NewRepo,
	dataperm.NewRepo,
	// 跨上下文最小依赖接口绑定
	wire.Bind(new(auth.UserStore), new(bizuser.Repo)),
	wire.Bind(new(auth.RoleNameReader), new(bizrole.Repo)),
	wire.Bind(new(auth.PermissionReader), new(bizperm.Repo)),
)

var serviceSet = wire.NewSet(
	authsvc.NewService,
	usersvc.NewService,
	rolesvc.NewService,
	permsvc.NewService,
)

var providerSet = wire.NewSet(bizSet, dataRepoSet, serviceSet, ProvideConfig, server.NewHTTPServer)

// wireApp 由 wire 生成
func wireApp() (*server.HTTPServer, func(), error) {
	wire.Build(providerSet)
	return nil, nil, nil
}
