package role

import (
	"context"
	"errors"
)

// ErrRoleNotFound 角色不存在
var ErrRoleNotFound = errors.New("role not found")

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
	List(ctx context.Context, q Query, page, pageSize int) ([]*Role, int64, error)
	SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
}
