package permission

import (
	"context"
	"errors"
	"strings"

	"github.com/smilex/smilex-admin-gin/pkg/pagination"
)

// Usecase 权限领域用例
type Usecase struct {
	repo Repo
}

func NewUsecase(repo Repo) *Usecase { return &Usecase{repo: repo} }

func (uc *Usecase) Create(ctx context.Context, name, code string, t Type, method, path string, parentID uint, icon string, sort int) (*Permission, error) {
	if code != "" {
		if _, err := uc.repo.FindByCode(ctx, code); err == nil {
			return nil, ErrDuplicateCode
		}
	}
	p := &Permission{Name: name, Code: code, Type: t, Method: method, Path: path, ParentID: parentID, Icon: icon, Sort: sort}
	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Update 更新权限（空字符串字段不更新；icon 为指针，传空字符串可清空图标；parentID 为指针，传 0 可挪到顶级）
func (uc *Usecase) Update(ctx context.Context, id uint, name, method, path string, icon *string, sort int, parentID *uint) error {
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
	if parentID != nil {
		if err := uc.validateParent(ctx, id, *parentID); err != nil {
			return err
		}
		p.ParentID = *parentID
	}
	return uc.repo.Update(ctx, p)
}

// validateParent 校验新父级合法性：不能是自身，也不能位于自身的后代链上（防环）
func (uc *Usecase) validateParent(ctx context.Context, id, parentID uint) error {
	if parentID == 0 {
		return nil
	}
	if parentID == id {
		return errors.New("父级不能是自身")
	}
	cur := parentID
	for cur != 0 {
		if cur == id {
			return errors.New("父级不能是自身的子级")
		}
		p, err := uc.repo.FindByID(ctx, cur)
		if err != nil {
			return err
		}
		cur = p.ParentID
	}
	return nil
}

// Delete 删除权限：超管通配权限（id=1）禁止删除；存在子级时须先删子级
func (uc *Usecase) Delete(ctx context.Context, id uint) error {
	if id == 1 {
		return errors.New("超管通配权限禁止删除")
	}
	n, err := uc.repo.CountByParentID(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrHasChildren
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *Usecase) Get(ctx context.Context, id uint) (*Permission, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *Usecase) List(ctx context.Context, q Query, page, pageSize int) ([]*Permission, pagination.Page, error) {
	ps, total, err := uc.repo.List(ctx, q, page, pageSize)
	if pageSize <= 0 { // 全量返回时分页元信息按单页全量填充
		page, pageSize = 1, int(total)
	}
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

// MenuHit 菜单搜索命中项（顶栏命令面板用；Parents 为父级链提示，不含自身）
// Dir 标记目录（含子菜单，无对应路由，前端选中时软提示选择具体菜单）
type MenuHit struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Icon    string `json:"icon"`
	Parents string `json:"parents"`
	Dir     bool   `json:"dir"`
}

// SearchUserMenus 当前用户可见菜单中按关键词模糊搜索（名称/路径，不区分大小写），上限 limit 条。
// 空关键词返回空。命中目录（含子菜单）时，该目录本身与其全部子孙一并纳入结果，
// 便于输入"系统"这类目录词时看到目录下的全部菜单。
func (uc *Usecase) SearchUserMenus(ctx context.Context, userID uint, kw string, limit int) ([]*MenuHit, error) {
	if kw == "" {
		return nil, nil
	}
	tree, err := uc.UserMenuTree(ctx, userID)
	if err != nil {
		return nil, err
	}
	lkwl := strings.ToLower(kw)
	var out []*MenuHit
	// forceIn：祖先已命中，整棵子树纳入
	var walk func(nodes []*MenuNode, parents string, forceIn bool)
	walk = func(nodes []*MenuNode, parents string, forceIn bool) {
		for _, n := range nodes {
			if limit > 0 && len(out) >= limit {
				return
			}
			selfHit := strings.Contains(strings.ToLower(n.Name), lkwl) || strings.Contains(strings.ToLower(n.Path), lkwl)
			if forceIn || selfHit {
				out = append(out, &MenuHit{
					Name: n.Name, Path: n.Path, Icon: n.Icon,
					Parents: parents, Dir: len(n.Children) > 0,
				})
			}
			if len(n.Children) > 0 {
				trail := n.Name
				if parents != "" {
					trail = parents + " / " + trail
				}
				walk(n.Children, trail, forceIn || selfHit)
			}
		}
	}
	walk(tree, "", false)
	return out, nil
}
