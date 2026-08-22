package user

import (
	"context"
	"errors"

	"github.com/smilex/smilex-admin-gin/pkg/eventbus"
	"github.com/smilex/smilex-admin-gin/pkg/pagination"
)

// 用户被创建事件（跨上下文通知，如发欢迎邮件、初始化数据等）
type CreatedEvent struct{ UserID uint }

func (CreatedEvent) Topic() string { return "user.created" }

// 超级管理员保护：admin（id=1）账号仅允许其本人操作，其余用户一律拒绝
var ErrSuperAdminProtected = errors.New("无权操作超级管理员账号")

// ErrDeleteSuperAdmin 超级管理员账号一律禁止删除（含其本人）
var ErrDeleteSuperAdmin = errors.New("超级管理员账号禁止删除")

const superAdminID uint = 1

type operatorKey struct{}

// WithOperator 将当前操作者 ID 注入 context（由传输层从认证信息中取出）
func WithOperator(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, operatorKey{}, userID)
}

func operatorFrom(ctx context.Context) uint {
	if v, ok := ctx.Value(operatorKey{}).(uint); ok {
		return v
	}
	return 0
}

// guardSuperAdmin 目标为超管且操作者不是超管本人时拒绝
func guardSuperAdmin(ctx context.Context, targetID uint) error {
	if targetID == superAdminID && operatorFrom(ctx) != superAdminID {
		return ErrSuperAdminProtected
	}
	return nil
}

// Usecase 用户领域用例
type Usecase struct {
	repo Repo
}

func NewUsecase(repo Repo) *Usecase { return &Usecase{repo: repo} }

func (uc *Usecase) Create(ctx context.Context, username, password, nickname, email string, roleIDs []uint) (*User, error) {
	if _, err := uc.repo.FindByUsername(ctx, username); err == nil {
		return nil, ErrDuplicateUsername
	}
	u := &User{Username: username, Nickname: nickname, Email: email, Status: StatusEnabled, RoleIDs: roleIDs}
	if err := u.SetPassword(password); err != nil {
		return nil, err
	}
	if err := uc.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	eventbus.Publish(CreatedEvent{UserID: u.ID})
	return u, nil
}

func (uc *Usecase) Update(ctx context.Context, id uint, nickname, email string, status *Status) error {
	if err := guardSuperAdmin(ctx, id); err != nil {
		return err
	}
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if nickname != "" {
		u.Nickname = nickname
	}
	if email != "" {
		u.Email = email
	}
	if status != nil {
		u.Status = *status
	}
	return uc.repo.Update(ctx, u)
}

func (uc *Usecase) Delete(ctx context.Context, id uint) error {
	if id == superAdminID {
		return ErrDeleteSuperAdmin
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *Usecase) Get(ctx context.Context, id uint) (*User, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *Usecase) List(ctx context.Context, q Query, page, pageSize int) ([]*User, pagination.Page, error) {
	users, total, err := uc.repo.List(ctx, q, page, pageSize)
	return users, pagination.Page{Page: page, PageSize: pageSize, Total: total}, err
}

// SetRoles 分配角色
func (uc *Usecase) SetRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	if err := guardSuperAdmin(ctx, userID); err != nil {
		return err
	}
	return uc.repo.SetRoles(ctx, userID, roleIDs)
}

// ResetPassword 管理员重置密码
func (uc *Usecase) ResetPassword(ctx context.Context, userID uint, newPlain string) error {
	if err := guardSuperAdmin(ctx, userID); err != nil {
		return err
	}
	u, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := u.SetPassword(newPlain); err != nil {
		return err
	}
	return uc.repo.UpdatePassword(ctx, userID, string(u.Password))
}
