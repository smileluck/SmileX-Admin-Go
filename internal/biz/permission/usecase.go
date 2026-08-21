package permission

import (
	"context"

	"github.com/smilex/smilex-admin-gin/pkg/pagination"
)

// Usecase 权限领域用例
type Usecase struct {
	repo Repo
}

func NewUsecase(repo Repo) *Usecase { return &Usecase{repo: repo} }

func (uc *Usecase) Create(ctx context.Context, name, code string, t Type, method, path string, parentID uint, icon string, sort int) (*Permission, error) {
	p := &Permission{Name: name, Code: code, Type: t, Method: method, Path: path, ParentID: parentID, Icon: icon, Sort: sort}
	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Update 更新权限（空字符串字段不更新；icon 为指针，传空字符串可清空图标）
func (uc *Usecase) Update(ctx context.Context, id uint, name, method, path string, icon *string, sort int) error {
	p, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if name != "" {
		p.Name = name
	}
	if method != "" {
		p.Method = method
	}
	if path != "" {
		p.Path = path
	}
	if icon != nil {
		p.Icon = *icon
	}
	p.Sort = sort
	return uc.repo.Update(ctx, p)
}

func (uc *Usecase) Delete(ctx context.Context, id uint) error { return uc.repo.Delete(ctx, id) }

func (uc *Usecase) Get(ctx context.Context, id uint) (*Permission, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *Usecase) List(ctx context.Context, q Query, page, pageSize int) ([]*Permission, pagination.Page, error) {
	ps, total, err := uc.repo.List(ctx, q, page, pageSize)
	return ps, pagination.Page{Page: page, PageSize: pageSize, Total: total}, err
}

// UserMenuTree 当前用户可见的菜单树
func (uc *Usecase) UserMenuTree(ctx context.Context, userID uint) ([]*MenuNode, error) {
	ps, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	menus := make([]*Permission, 0, len(ps))
	for _, p := range ps {
		if p.Type == TypeMenu {
			menus = append(menus, p)
		}
	}
	return BuildMenuTree(menus, 0), nil
}
