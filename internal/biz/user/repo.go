package user

import (
	"context"
	"errors"
)

// ErrUserNotFound 用户不存在
var ErrUserNotFound = errors.New("用户不存在")

// ErrDuplicateUsername 用户名重复
var ErrDuplicateUsername = errors.New("用户名已存在，请更换")

// Query 用户列表查询条件
type Query struct {
	Username string
	Status   *int
}

// Repo 用户仓储接口（由 data 层实现，依赖倒置）
type Repo interface {
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context, q Query, page, pageSize int) ([]*User, int64, error)
	SetRoles(ctx context.Context, userID uint, roleIDs []uint) error
}
