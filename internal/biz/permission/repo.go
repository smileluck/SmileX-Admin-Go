package permission

import (
	"context"
	"errors"
)

// ErrPermissionNotFound 权限不存在
var ErrPermissionNotFound = errors.New("permission not found")

// Query 权限列表查询条件
type Query struct {
	Type string // api | menu，空为全部
}

// Repo 权限仓储接口
type Repo interface {
	Create(ctx context.Context, p *Permission) error
	Update(ctx context.Context, p *Permission) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*Permission, error)
	List(ctx context.Context, q Query, page, pageSize int) ([]*Permission, int64, error)
	// FindByUserID 查询用户经角色关联到的全部权限（data 层做 join）
	FindByUserID(ctx context.Context, userID uint) ([]*Permission, error)
}
