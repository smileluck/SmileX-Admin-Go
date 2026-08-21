package role

import (
	"context"
	"errors"
)

// ErrRoleNotFound 角色不存在
var ErrRoleNotFound = errors.New("角色不存在")

// ErrRoleHasUsers 角色下仍有用户，须先移除
var ErrRoleHasUsers = errors.New("该角色下存在用户，请先移除用户与该角色的关联")

// ErrDuplicateRoleCode 角色编码已存在
var ErrDuplicateRoleCode = errors.New("角色编码已存在，请更换")

// Query 角色列表查询条件
type Query struct {
	Name string
}

// Repo 角色仓储接口
type Repo interface {
	Create(ctx context.Context, r *Role) error
	Update(ctx context.Context, r *Role) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*Role, error)
	// FindByCode 按编码查询（创建时查重用）
	FindByCode(ctx context.Context, code string) (*Role, error)
	// CountUsers 统计角色下的用户数量（删除保护用）
	CountUsers(ctx context.Context, roleID uint) (int64, error)
	List(ctx context.Context, q Query, page, pageSize int) ([]*Role, int64, error)
	SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
}
