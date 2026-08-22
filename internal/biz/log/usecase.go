package log

import (
	"context"
	"time"

	"github.com/smilex/smilex-admin-gin/internal/conf"
	"github.com/smilex/smilex-admin-gin/pkg/pagination"
)

// Usecase 日志领域用例（记录异步、查询/清空同步）
type Usecase struct {
	repo          Repo
	retentionDays int // 保留天数（0 = 永久），保留期自动清理与前端说明共用
}

func NewUsecase(repo Repo, c *conf.Bootstrap) *Usecase {
	return &Usecase{repo: repo, retentionDays: c.Log.RetentionDays}
}

// RecordLogin 记录一次登录尝试（异步落库，不阻塞登录响应）
func (uc *Usecase) RecordLogin(l *LoginLog) {
	if l.CreatedAt.IsZero() {
		l.CreatedAt = time.Now()
	}
	uc.repo.CreateLogin(l)
}

// RecordOperation 记录一条写请求审计（异步落库）
func (uc *Usecase) RecordOperation(o *OperationLog) {
	if o.CreatedAt.IsZero() {
		o.CreatedAt = time.Now()
	}
	uc.repo.CreateOperation(o)
}

// ListLoginLogs 分页查询登录日志
func (uc *Usecase) ListLoginLogs(ctx context.Context, q LoginLogQuery, page, pageSize int) ([]*LoginLog, pagination.Page, error) {
	logs, total, err := uc.repo.ListLoginLogs(ctx, q, page, pageSize)
	if err != nil {
		return nil, pagination.Page{}, err
	}
	return logs, pagination.Page{Page: page, PageSize: pageSize, Total: total}, nil
}

// ListOperationLogs 分页查询操作日志
func (uc *Usecase) ListOperationLogs(ctx context.Context, q OperationLogQuery, page, pageSize int) ([]*OperationLog, pagination.Page, error) {
	logs, total, err := uc.repo.ListOperationLogs(ctx, q, page, pageSize)
	if err != nil {
		return nil, pagination.Page{}, err
	}
	return logs, pagination.Page{Page: page, PageSize: pageSize, Total: total}, nil
}

// ClearLoginLogs 清空登录日志（与保留期清理共用按时间删除路径），返回删除行数
func (uc *Usecase) ClearLoginLogs(ctx context.Context) (int64, error) {
	return uc.repo.DeleteLoginBefore(ctx, time.Now())
}

// ClearOperationLogs 清空操作日志，返回删除行数
func (uc *Usecase) ClearOperationLogs(ctx context.Context) (int64, error) {
	return uc.repo.DeleteOperationBefore(ctx, time.Now())
}

// RetentionDays 日志保留天数（0 = 永久保留）
func (uc *Usecase) RetentionDays() int {
	return uc.retentionDays
}
