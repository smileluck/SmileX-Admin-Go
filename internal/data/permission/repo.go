// Package permission 权限仓储 GORM 实现
package permission

import (
	"context"
	"errors"

	bizperm "github.com/smilex/smilex-admin-gin/internal/biz/permission"
	"github.com/smilex/smilex-admin-gin/internal/data"
	"github.com/smilex/smilex-admin-gin/internal/data/model"
	"gorm.io/gorm"
)

type repo struct {
	data *data.Data
}

// NewRepo 创建权限仓储
func NewRepo(d *data.Data) bizperm.Repo { return &repo{data: d} }

func mapErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return bizperm.ErrPermissionNotFound
	}
	return err
}

func (r *repo) Create(ctx context.Context, m *bizperm.Permission) error {
	return r.data.DB.WithContext(ctx).Create(model.PermissionToPO(m)).Error
}

func (r *repo) Update(ctx context.Context, m *bizperm.Permission) error {
	return r.data.DB.WithContext(ctx).Model(&model.PermissionPO{}).Where("id = ?", m.ID).
		Updates(map[string]interface{}{"name": m.Name, "method": m.Method, "path": m.Path, "icon": m.Icon, "sort": m.Sort}).Error
}

func (r *repo) Delete(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.PermissionPO{}, id).Error; err != nil {
			return err
		}
		return tx.Where("permission_id = ?", id).Delete(&model.RolePermissionPO{}).Error
	})
}

func (r *repo) FindByID(ctx context.Context, id uint) (*bizperm.Permission, error) {
	var po model.PermissionPO
	if err := r.data.DB.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, mapErr(err)
	}
	return model.PermissionFromPO(&po), nil
}

func (r *repo) List(ctx context.Context, q bizperm.Query, page, pageSize int) ([]*bizperm.Permission, int64, error) {
	tx := r.data.DB.WithContext(ctx).Model(&model.PermissionPO{})
	if q.Type != "" {
		tx = tx.Where("type = ?", q.Type)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var pos []model.PermissionPO
	if err := tx.Offset((page - 1) * pageSize).Limit(pageSize).Order("id").Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*bizperm.Permission, 0, len(pos))
	for i := range pos {
		out = append(out, model.PermissionFromPO(&pos[i]))
	}
	return out, total, nil
}

// FindByUserID users -> user_roles -> role_permissions -> permissions 联查
func (r *repo) FindByUserID(ctx context.Context, userID uint) ([]*bizperm.Permission, error) {
	var pos []model.PermissionPO
	err := r.data.DB.WithContext(ctx).
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN user_roles ur ON ur.role_id = rp.role_id").
		Where("ur.user_id = ?", userID).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}
	out := make([]*bizperm.Permission, 0, len(pos))
	for i := range pos {
		out = append(out, model.PermissionFromPO(&pos[i]))
	}
	return out, nil
}
