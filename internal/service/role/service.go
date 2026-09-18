// Package role 角色应用服务
package role

import (
	"context"

	bizrole "github.com/smilex/smilex-admin-gin/internal/biz/role"
)

type Service struct {
	uc *bizrole.Usecase
}

func NewService(uc *bizrole.Usecase) *Service { return &Service{uc: uc} }

type CreateRequest struct {
	Name   string `json:"name" binding:"required,max=20"`
	Remark string `json:"remark" binding:"max=200"`
}

type UpdateRequest struct {
	Name   string `json:"name" binding:"omitempty,max=20"` // 空=保持原名
	Remark string `json:"remark" binding:"max=200"`
}

type SetPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*bizrole.Role, error) {
	return s.uc.Create(ctx, req.Name, req.Remark)
}

func (s *Service) Update(ctx context.Context, id uint, req UpdateRequest) error {
	return s.uc.Update(ctx, id, req.Name, req.Remark)
}

func (s *Service) Delete(ctx context.Context, id uint) error { return s.uc.Delete(ctx, id) }

func (s *Service) Get(ctx context.Context, id uint) (*bizrole.Role, error) { return s.uc.Get(ctx, id) }

func (s *Service) List(ctx context.Context, q bizrole.Query, page, pageSize int) ([]*bizrole.Role, interface{}, error) {
	return s.uc.List(ctx, q, page, pageSize)
}

// SetPermissions operatorID 为当前操作者：只能分配自己拥有的权限（超管除外，见 biz 层）
func (s *Service) SetPermissions(ctx context.Context, operatorID, id uint, req SetPermissionsRequest) error {
	return s.uc.SetPermissions(ctx, operatorID, id, req.PermissionIDs)
}
