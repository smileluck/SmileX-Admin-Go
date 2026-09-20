package role

import (
	"context"

	bizperm "github.com/smilex/smilex-admin-gin/internal/biz/permission"
	"github.com/smilex/smilex-admin-gin/pkg/pagination"
)

// PermissionReader 跨上下文最小依赖：读取某用户拥有的权限集合
type PermissionReader interface {
	FindByUserID(ctx context.Context, userID uint) ([]*bizperm.Permission, error)
}

// Usecase 角色领域用例
type Usecase struct {
	repo   Repo
	permRd PermissionReader
}

// superAdminRoleID 超管角色固定 ID：禁止修改和操作
const superAdminRoleID uint = 1

func NewUsecase(repo Repo, permRd PermissionReader) *Usecase {
	return &Usecase{repo: repo, permRd: permRd}
}

func (uc *Usecase) Create(ctx context.Context, name, remark string) (*Role, error) {
	r := &Role{Name: name, Remark: remark}
	if err := uc.repo.Create(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (uc *Usecase) Update(ctx context.Context, id uint, name, remark string) error {
	if id == superAdminRoleID {
		return ErrSuperRoleLocked
	}
	r, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if name != "" {
		r.Name = name
	}
	if remark != "" {
		r.Remark = remark
	}
	return uc.repo.Update(ctx, r)
}

func (uc *Usecase) Delete(ctx context.Context, id uint) error {
	if id == superAdminRoleID {
		return ErrSuperRoleLocked
	}
	if n, err := uc.repo.CountUsers(ctx, id); err != nil {
		return err
	} else if n > 0 {
		return ErrRoleHasUsers
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *Usecase) Get(ctx context.Context, id uint) (*Role, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *Usecase) List(ctx context.Context, q Query, page, pageSize int) ([]*Role, pagination.Page, error) {
	roles, total, err := uc.repo.List(ctx, q, page, pageSize)
	return roles, pagination.Page{Page: page, PageSize: pageSize, Total: total}, err
}

// SetPermissions 绑定权限：提交的权限必须是操作者本人权限的子集（防提权；
// 超管角色成员经 FindByUserID 返回全量权限，天然通过校验）
func (uc *Usecase) SetPermissions(ctx context.Context, operatorID, roleID uint, permissionIDs []uint) error {
	if roleID == superAdminRoleID {
		return ErrSuperRoleLocked
	}
	own, err := uc.permRd.FindByUserID(ctx, operatorID)
	if err != nil {
		return err
	}
	ownSet := make(map[uint]struct{}, len(own))
	for _, p := range own {
		ownSet[p.ID] = struct{}{}
	}
	for _, pid := range permissionIDs {
		if _, ok := ownSet[pid]; !ok {
			return ErrPermExceedsOperator
		}
	}
	return uc.repo.SetPermissions(ctx, roleID, permissionIDs)
}
