// Package model GORM 持久化对象（PO）与领域实体的转换。
// PO 只在 data 层出现，biz 层不可见。
package model

import (
	"time"

	"github.com/smilex/smilex-admin-gin/internal/biz/permission"
	"github.com/smilex/smilex-admin-gin/internal/biz/role"
	"github.com/smilex/smilex-admin-gin/internal/biz/user"
	"gorm.io/gorm"
)

// UserPO 用户表
type UserPO struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"size:64;uniqueIndex"`
	Password  string `gorm:"size:128"`
	Nickname  string `gorm:"size:64"`
	Email     string `gorm:"size:128"`
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserPO) TableName() string { return "users" }

// RolePO 角色表
type RolePO struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:64"`
	Code      string `gorm:"size:64;uniqueIndex"`
	Remark    string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (RolePO) TableName() string { return "roles" }

// PermissionPO 权限表
type PermissionPO struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:64"`
	Code      string `gorm:"size:64;uniqueIndex"`
	Type      string `gorm:"size:16"` // api | menu
	Method    string `gorm:"size:16"`
	Path      string `gorm:"size:255"`
	ParentID  uint
	Icon      string `gorm:"size:64"`
	Sort      int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (PermissionPO) TableName() string { return "permissions" }

// UserRolePO 用户-角色关联
type UserRolePO struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

func (UserRolePO) TableName() string { return "user_roles" }

// RolePermissionPO 角色-权限关联
type RolePermissionPO struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

func (RolePermissionPO) TableName() string { return "role_permissions" }

// ---- 转换器 ----

func UserToPO(u *user.User) *UserPO {
	return &UserPO{
		ID: u.ID, Username: u.Username, Password: string(u.Password),
		Nickname: u.Nickname, Email: u.Email, Status: int(u.Status),
	}
}

func UserFromPO(p *UserPO) *user.User {
	return &user.User{
		ID: p.ID, Username: p.Username, Password: user.Password(p.Password),
		Nickname: p.Nickname, Email: p.Email, Status: user.Status(p.Status),
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func RoleToPO(r *role.Role) *RolePO {
	return &RolePO{ID: r.ID, Name: r.Name, Code: r.Code, Remark: r.Remark}
}

func RoleFromPO(p *RolePO) *role.Role {
	return &role.Role{ID: p.ID, Name: p.Name, Code: p.Code, Remark: p.Remark, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

func PermissionToPO(m *permission.Permission) *PermissionPO {
	return &PermissionPO{
		ID: m.ID, Name: m.Name, Code: m.Code, Type: string(m.Type),
		Method: m.Method, Path: m.Path, ParentID: m.ParentID,
		Icon: m.Icon, Sort: m.Sort,
	}
}

func PermissionFromPO(p *PermissionPO) *permission.Permission {
	return &permission.Permission{
		ID: p.ID, Name: p.Name, Code: p.Code, Type: permission.Type(p.Type),
		Method: p.Method, Path: p.Path, ParentID: p.ParentID,
		Icon: p.Icon, Sort: p.Sort,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}
