// Package data 基础设施层：多数据库工厂、AutoMigrate 与种子数据。
package data

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/smilex/smilex-admin-gin/internal/biz/permission"
	"github.com/smilex/smilex-admin-gin/internal/biz/user"
	"github.com/smilex/smilex-admin-gin/internal/conf"
	"github.com/smilex/smilex-admin-gin/internal/data/model"
	"github.com/smilex/smilex-admin-gin/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"go.uber.org/zap"
)

// Data 持久化入口（对应 Kratos 的 Data struct）
type Data struct {
	DB *gorm.DB
}

// NewData 按 config.db.driver 创建对应数据库连接
func NewData(c *conf.Bootstrap) (*Data, func(), error) {
	var (
		db  *gorm.DB
		err error
	)
	gormCfg := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	}
	switch c.DB.Driver {
	case "mysql":
		ensureMySQLDatabase(c.DB.MySQL)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			c.DB.MySQL.User, c.DB.MySQL.Password, c.DB.MySQL.Host, c.DB.MySQL.Port, c.DB.MySQL.DBName, c.DB.MySQL.Charset)
		db, err = gorm.Open(mysql.Open(dsn), gormCfg)
	case "postgres":
		ensurePostgresDatabase(c.DB.Postgres)
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.DB.Postgres.Host, c.DB.Postgres.Port, c.DB.Postgres.User, c.DB.Postgres.Password, c.DB.Postgres.DBName, c.DB.Postgres.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), gormCfg)
	case "sqlite":
		_ = os.MkdirAll(filepath.Dir(c.DB.SQLite.Path), 0o755)
		db, err = gorm.Open(sqlite.Open(c.DB.SQLite.Path), gormCfg)
	default:
		return nil, nil, fmt.Errorf("unsupported db driver: %s", c.DB.Driver)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", c.DB.Driver, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	maxOpen, maxIdle := 20, 10
	switch c.DB.Driver {
	case "mysql":
		maxOpen, maxIdle = c.DB.MySQL.MaxOpenConns, c.DB.MySQL.MaxIdleConns
	case "postgres":
		maxOpen, maxIdle = c.DB.Postgres.MaxOpenConns, c.DB.Postgres.MaxIdleConns
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)

	d := &Data{DB: db}
	if c.DB.AutoMigrate {
		if err := d.migrateAndSeed(); err != nil {
			return nil, nil, err
		}
	}
	cleanup := func() { _ = sqlDB.Close() }
	return d, cleanup, nil
}

// ensureMySQLDatabase 库不存在时自动创建（连接 information_schema 建库后再正常连接）
func ensureMySQLDatabase(c conf.MySQL) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	defer db.Close()
	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET %s", c.DBName, c.Charset)); err != nil {
		logger.Warn("auto create mysql database failed", zap.Error(err))
	}
}

// ensurePostgresDatabase 库不存在时自动创建（先连 postgres 库执行 CREATE DATABASE）
func ensurePostgresDatabase(c conf.Postgres) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	defer sqlDB.Close()
	if _, err := sqlDB.Exec(fmt.Sprintf(`SELECT 'CREATE DATABASE %s' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')`, c.DBName, c.DBName)); err != nil {
		return
	}
	if _, err := sqlDB.Exec("CREATE DATABASE " + c.DBName); err != nil {
		logger.Warn("auto create postgres database failed", zap.Error(err))
	}
}

// migrateAndSeed 自动建表 + 种子数据（超管 admin/123456）
func (d *Data) migrateAndSeed() error {
	if err := d.DB.AutoMigrate(
		&model.UserPO{}, &model.RolePO{}, &model.PermissionPO{},
		&model.UserRolePO{}, &model.RolePermissionPO{},
	); err != nil {
		return err
	}

	var userCount int64
	if err := d.DB.Model(&model.UserPO{}).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return nil
	}

	// 超管角色 + 通配权限
	rolePO := model.RolePO{ID: 1, Name: "超级管理员", Code: "super_admin", Remark: "拥有全部权限"}
	permPO := model.PermissionPO{ID: 1, Name: "全部权限", Code: "all", Type: string(permission.TypeAPI), Method: "*", Path: "*"}

	// 菜单种子数据（type=menu）
	menus := []model.PermissionPO{
		{ID: 100, Name: "首页", Code: "menu:dashboard", Type: "menu", Path: "/dashboard", Icon: "HomeOutline", Sort: 1},
		{ID: 110, Name: "系统管理", Code: "menu:system", Type: "menu", Path: "/system", Icon: "SettingsOutline", Sort: 2},
		{ID: 111, Name: "用户管理", Code: "menu:user", Type: "menu", Path: "/system/users", ParentID: 110, Icon: "PersonOutline", Sort: 1},
		{ID: 112, Name: "角色管理", Code: "menu:role", Type: "menu", Path: "/system/roles", ParentID: 110, Icon: "IdCardOutline", Sort: 2},
		{ID: 113, Name: "权限管理", Code: "menu:permission", Type: "menu", Path: "/system/permissions", ParentID: 110, Icon: "LockClosedOutline", Sort: 3},
		{ID: 114, Name: "菜单管理", Code: "menu:menu", Type: "menu", Path: "/system/menus", ParentID: 110, Icon: "MenuOutline", Sort: 4},
	}
	adminPwd, err := user.NewPassword("123456")
	if err != nil {
		return err
	}
	adminPO := model.UserPO{ID: 1, Username: "admin", Password: string(adminPwd), Nickname: "超级管理员", Email: "admin@smilex.local", Status: int(user.StatusEnabled)}

	return d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&permPO).Error; err != nil {
			return err
		}
		if err := tx.Create(&rolePO).Error; err != nil {
			return err
		}
		if err := tx.Create(&adminPO).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.UserRolePO{UserID: 1, RoleID: 1}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.RolePermissionPO{RoleID: 1, PermissionID: 1}).Error; err != nil {
			return err
		}
		for i := range menus {
			if err := tx.Create(&menus[i]).Error; err != nil {
				return err
			}
			if err := tx.Create(&model.RolePermissionPO{RoleID: 1, PermissionID: menus[i].ID}).Error; err != nil {
				return err
			}
		}
		logger.Info("seeded super admin: admin/123456 (please change password before production)")
		return nil
	})
}
