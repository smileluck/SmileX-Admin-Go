package blacklist

import "context"

// Repo 黑名单仓储接口（由 data 层实现，依赖倒置）
type Repo interface {
	// Create 新增黑名单记录
	Create(ctx context.Context, b *IPBlacklist) error
	// Delete 按 ID 删除（解封）
	Delete(ctx context.Context, id uint) error
	// List 分页查询（page_size=0 全量）
	List(ctx context.Context, q Query, page, pageSize int) ([]*IPBlacklist, int64, error)
	// ListActive 返回当前生效的 IP 列表（未过期且未被软删），供中间件拦截使用
	ListActive(ctx context.Context) ([]string, error)
}
