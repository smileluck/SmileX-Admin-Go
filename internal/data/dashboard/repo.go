// Package dashboard 仪表盘聚合仓储 GORM 实现。
// 日期聚合用 DATE()：MySQL / PostgreSQL / SQLite 三方言均原生支持。
package dashboard

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/smilex/smilex-admin-gin/internal/biz/dashboard"
	"github.com/smilex/smilex-admin-gin/internal/data"
	"github.com/smilex/smilex-admin-gin/internal/data/model"
)

type repo struct {
	data *data.Data
}

func NewRepo(d *data.Data) dashboard.Repo { return &repo{data: d} }

func (r *repo) Counts() (users, roles int64, err error) {
	db := r.data.DB
	if err = db.Model(&model.UserPO{}).Count(&users).Error; err != nil {
		return
	}
	err = db.Model(&model.RolePO{}).Count(&roles).Error
	return
}

func (r *repo) TodayLogins() (int64, error) {
	var n int64
	start := time.Now().Truncate(24 * time.Hour)
	err := r.data.DB.Model(&model.LoginLogPO{}).
		Where("created_at >= ?", start).Count(&n).Error
	return n, err
}

// LoginTrend 按日聚合登录（总量与成功量）
func (r *repo) LoginTrend(days int) ([]dashboard.DailyPoint, error) {
	since := time.Now().AddDate(0, 0, -(days - 1)).Format("2006-01-02")
	var rows []struct {
		Date    string `gorm:"column:d"`
		Total   int64  `gorm:"column:total"`
		Success int64  `gorm:"column:success"`
	}
	err := r.data.DB.Model(&model.LoginLogPO{}).
		Select("DATE(created_at) AS d, COUNT(*) AS total, SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) AS success").
		Where("created_at >= ?", since).
		Group("d").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]dashboard.DailyPoint, 0, len(rows))
	for _, w := range rows {
		out = append(out, dashboard.DailyPoint{Date: w.Date, Total: w.Total, Success: w.Success})
	}
	return out, nil
}

// OpTrend 按日聚合操作日志量
func (r *repo) OpTrend(days int) ([]dashboard.DailyPoint, error) {
	since := time.Now().AddDate(0, 0, -(days - 1)).Format("2006-01-02")
	var rows []struct {
		Date  string `gorm:"column:d"`
		Total int64  `gorm:"column:total"`
	}
	err := r.data.DB.Model(&model.OperationLogPO{}).
		Select("DATE(created_at) AS d, COUNT(*) AS total").
		Where("created_at >= ?", since).
		Group("d").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]dashboard.DailyPoint, 0, len(rows))
	for _, w := range rows {
		out = append(out, dashboard.DailyPoint{Date: w.Date, Total: w.Total})
	}
	return out, nil
}

func (r *repo) RecentLogins(n int) ([]dashboard.LoginItem, error) {
	var pos []model.LoginLogPO
	if err := r.data.DB.WithContext(context.Background()).
		Order("id DESC").Limit(n).Find(&pos).Error; err != nil {
		return nil, err
	}
	out := make([]dashboard.LoginItem, 0, len(pos))
	for _, p := range pos {
		status := dashboard.LoginFail
		if p.Status == 1 {
			status = dashboard.LoginSuccess
		}
		out = append(out, dashboard.LoginItem{
			Username: p.Username, IP: p.IP,
			Status: status, CreatedAt: p.CreatedAt,
		})
	}
	return out, nil
}

var _ dashboard.Repo = (*repo)(nil)
var _ = gorm.ErrRecordNotFound // 保留 gorm 引用（扫描目标为匿名结构体）
