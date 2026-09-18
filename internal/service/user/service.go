// Package user 用户应用服务
package user

import (
	"context"

	bizuser "github.com/smilex/smilex-admin-gin/internal/biz/user"
	"github.com/smilex/smilex-admin-gin/pkg/security"
)

type Service struct {
	uc *bizuser.Usecase
}

func NewService(uc *bizuser.Usecase) *Service { return &Service{uc: uc} }

// CreateRequest 创建用户入参
type CreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Nickname string `json:"nickname" binding:"max=20"`
	Phone    string `json:"phone" binding:"omitempty,max=32"`
	Email    string `json:"email" binding:"omitempty,max=128,email"`
	RoleIDs  []uint `json:"role_ids"`
}

// UpdateRequest 更新用户入参
type UpdateRequest struct {
	Nickname string `json:"nickname" binding:"max=20"`
	Phone    string `json:"phone" binding:"omitempty,max=32"`
	Email    string `json:"email" binding:"omitempty,max=128,email"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

type SetRolesRequest struct {
	RoleIDs []uint `json:"role_ids"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// ListVO 列表视图（隐藏密码）
type ListVO struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Status    int    `json:"status"`
	RoleIDs   []uint `json:"role_ids"`
	CreatedAt string `json:"created_at"`
}

// toVO 用户实体转视图：手机号/邮箱默认脱敏，reveal 为 true 时输出明文
// （须由 handler 层按 user:viewSensitive 校验后传入）
func toVO(u *bizuser.User, reveal bool) *ListVO {
	phone, email := security.MaskPhone(u.Phone), security.MaskEmail(u.Email)
	if reveal {
		phone, email = u.Phone, u.Email
	}
	return &ListVO{
		ID: u.ID, Username: u.Username, Nickname: u.Nickname, Phone: phone, Email: email,
		Status: int(u.Status), RoleIDs: u.RoleIDs, CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*ListVO, error) {
	u, err := s.uc.Create(ctx, req.Username, req.Password, req.Nickname, req.Phone, req.Email, req.RoleIDs)
	if err != nil {
		return nil, err
	}
	return toVO(u, false), nil
}

func (s *Service) Update(ctx context.Context, id uint, req UpdateRequest) error {
	var st *bizuser.Status
	if req.Status != nil {
		sv := bizuser.Status(*req.Status)
		st = &sv
	}
	return s.uc.Update(ctx, id, req.Nickname, req.Phone, req.Email, st)
}

func (s *Service) Delete(ctx context.Context, id uint) error { return s.uc.Delete(ctx, id) }

func (s *Service) Get(ctx context.Context, id uint, reveal bool) (*ListVO, error) {
	u, err := s.uc.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toVO(u, reveal), nil
}

func (s *Service) List(ctx context.Context, q bizuser.Query, page, pageSize int, reveal bool) ([]*ListVO, interface{}, error) {
	users, pg, err := s.uc.List(ctx, q, page, pageSize)
	if err != nil {
		return nil, nil, err
	}
	out := make([]*ListVO, 0, len(users))
	for _, u := range users {
		out = append(out, toVO(u, reveal))
	}
	return out, pg, nil
}

func (s *Service) SetRoles(ctx context.Context, id uint, req SetRolesRequest) error {
	return s.uc.SetRoles(ctx, id, req.RoleIDs)
}

func (s *Service) ResetPassword(ctx context.Context, id uint, req ResetPasswordRequest) error {
	return s.uc.ResetPassword(ctx, id, req.Password)
}
