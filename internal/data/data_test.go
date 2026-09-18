package data

import (
	"testing"

	"github.com/smilex/smilex-admin-gin/internal/conf"
	"github.com/smilex/smilex-admin-gin/internal/data/model"
	"github.com/smilex/smilex-admin-gin/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func init() { logger.L = zap.NewNop() } // 播种日志静音，避免测试依赖日志初始化

// newSeedTestDB 内存 sqlite 构建 Data（不触发 NewData 的连接池/清理逻辑，直接复用 migrateAndSeed）
func newSeedTestDB(t *testing.T) *Data {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close() // 释放共享内存库，避免用例间串数据
		}
	})
	return &Data{DB: db, cfg: &conf.Bootstrap{Seed: conf.Seed{AdminPassword: "seed-test-123456"}}}
}

func countBy(t *testing.T, db *gorm.DB, model_ interface{}, cond string, args ...interface{}) int64 {
	t.Helper()
	var n int64
	if err := db.Unscoped().Model(model_).Where(cond, args...).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	return n
}

// TestMigrateAndSeed_FreshAndIdempotent 全新库播种齐备；重复执行不产生重复行
func TestMigrateAndSeed_FreshAndIdempotent(t *testing.T) {
	d := newSeedTestDB(t)
	if err := d.migrateAndSeed(); err != nil {
		t.Fatal(err)
	}

	// 首次播种：超管/角色/通配权限/固定 ID 菜单 + ensure 补齐的清单菜单与按钮权限
	var admin model.UserPO
	if err := d.DB.Where("id = ? AND username = ?", 1, "admin").First(&admin).Error; err != nil {
		t.Fatalf("admin not seeded: %v", err)
	}
	if n := countBy(t, d.DB, &model.RolePO{}, "id = ?", 1); n != 1 {
		t.Fatalf("super role not seeded: %d", n)
	}
	for _, code := range []string{"all", "menu:dashboard", "menu:openapi", "tenant:status", "user:viewSensitive"} {
		if n := countBy(t, d.DB, &model.PermissionPO{}, "code = ?", code); n != 1 {
			t.Fatalf("permission %s: want 1, got %d", code, n)
		}
	}

	// 幂等：再跑一遍，权限/绑定/用户行数均不变
	var perms, binds, users int64
	_ = d.DB.Unscoped().Model(&model.PermissionPO{}).Count(&perms).Error
	_ = d.DB.Model(&model.RolePermissionPO{}).Count(&binds).Error
	_ = d.DB.Model(&model.UserPO{}).Count(&users).Error
	if err := d.migrateAndSeed(); err != nil {
		t.Fatal(err)
	}
	if n := countBy(t, d.DB, &model.PermissionPO{}, "1 = 1"); n != perms {
		t.Fatalf("permission rows duplicated after reseed: %d -> %d", perms, n)
	}
	if n := countBy(t, d.DB, &model.RolePermissionPO{}, "1 = 1"); n != binds {
		t.Fatalf("role bindings duplicated after reseed: %d -> %d", binds, n)
	}
	if n := countBy(t, d.DB, &model.UserPO{}, "1 = 1"); n != users {
		t.Fatalf("user rows changed after reseed: %d -> %d", users, n)
	}
}

// TestEnsureSystemMenus_SelfHeal 菜单漂移（改名/改排序/软删）会被启动对齐回规范定义
func TestEnsureSystemMenus_SelfHeal(t *testing.T) {
	d := newSeedTestDB(t)
	if err := d.migrateAndSeed(); err != nil {
		t.Fatal(err)
	}

	// 手工漂移：登录日志改名为"登录记录"、sort 改 99、icon 清空，文件管理软删
	if err := d.DB.Model(&model.PermissionPO{}).Where("code = ?", "menu:loginLog").
		Updates(map[string]interface{}{"name": "登录记录", "sort": 99, "icon": ""}).Error; err != nil {
		t.Fatal(err)
	}
	if err := d.DB.Delete(&model.PermissionPO{}, "code = ?", "menu:file").Error; err != nil {
		t.Fatal(err)
	}

	if err := d.ensureSystemMenus(); err != nil {
		t.Fatal(err)
	}

	var m model.PermissionPO
	if err := d.DB.Unscoped().Where("code = ?", "menu:loginLog").First(&m).Error; err != nil {
		t.Fatal(err)
	}
	if m.Name != "登录日志" || m.Sort != 1 || m.Icon != "LogInOutline" {
		t.Errorf("menu drift not healed: %+v", m)
	}
	var f model.PermissionPO
	if err := d.DB.Unscoped().Where("code = ?", "menu:file").First(&f).Error; err != nil {
		t.Fatal(err)
	}
	if f.DeletedAt.Valid {
		t.Errorf("soft-deleted system menu not restored: %+v", f)
	}
}
