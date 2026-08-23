// Package blacklist 黑名单仓储 GORM 实现
package blacklist

import (
	"context"
	"time"

	bizblacklist "github.com/smilex/smilex-admin-gin/internal/biz/blacklist"
	"github.com/smilex/smilex-admin-gin/internal/data"
	"github.com/smilex/smilex-admin-gin/internal/data/model"
	"github.com/smilex/smilex-admin-gin/pkg/security"
	"gorm.io/gorm"
)

// Repo 黑名单仓储
type Repo struct {
	data *data.Data
}

func NewRepo(d *data.Data) *Repo {
	return &Repo{data: d}
}

func (r *Repo) Create(ctx context.Context, b *bizblacklist.IPBlacklist) error {
	po := model.IPBlacklistToPO(b)
	po.CreatedAt = time.Now()
	po.UpdatedAt = po.CreatedAt
	return r.data.DB.WithContext(ctx).Create(po).Error
}

func (r *Repo) Delete(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Delete(&model.IPBlacklistPO{}, id).Error
}

func (r *Repo) List(ctx context.Context, q bizblacklist.Query, page, pageSize int) ([]*bizblacklist.IPBlacklist, int64, error) {
	tx := r.data.DB.WithContext(ctx).Model(&model.IPBlacklistPO{})
	if q.IP != "" {
		tx = tx.Where("ip LIKE ? ESCAPE '/'", security.EscapeLike(q.IP)+"%")
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var pos []model.IPBlacklistPO
	tx = tx.Order("id DESC")
	if pageSize > 0 {
		tx = tx.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	if err := tx.Find(&pos).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*bizblacklist.IPBlacklist, 0, len(pos))
	for i := range pos {
		out = append(out, model.IPBlacklistFromPO(&pos[i]))
	}
	return out, total, nil
}

func (r *Repo) ListActive(ctx context.Context) ([]string, error) {
	var ips []string
	now := time.Now()
	err := r.data.DB.WithContext(ctx).Model(&model.IPBlacklistPO{}).
		Where("(expire_at IS NULL OR expire_at > ?) AND deleted_at IS NULL", now).
		Pluck("ip", &ips).Error
	return ips, err
}

// listPage 占位兼容，如需统一可替换（本文件暂用内联）
var _ = gorm.DB{}
