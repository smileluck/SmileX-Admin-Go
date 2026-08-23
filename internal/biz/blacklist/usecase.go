package blacklist

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/smilex/smilex-admin-gin/pkg/pagination"
)

// Usecase IP 黑名单领域用例
type Usecase struct {
	repo Repo

	mu         sync.RWMutex
	cache      map[string]struct{}
	expireAt   time.Time
	cacheValid bool
}

func NewUsecase(repo Repo) *Usecase {
	return &Usecase{repo: repo}
}

// Create 新增 IP 黑名单
func (uc *Usecase) Create(ctx context.Context, b *IPBlacklist) error {
	b.IP = strings.TrimSpace(b.IP)
	parsed := net.ParseIP(b.IP)
	if parsed == nil {
		return ErrInvalidIP
	}
	b.IP = parsed.String()
	if b.ExpireAt != nil && !b.ExpireAt.After(time.Now()) {
		return ErrInvalidExpire
	}
	if err := uc.repo.Create(ctx, b); err != nil {
		return err
	}
	uc.invalidate()
	return nil
}

// Delete 解封指定 ID
func (uc *Usecase) Delete(ctx context.Context, id uint) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}
	uc.invalidate()
	return nil
}

// List 分页查询
func (uc *Usecase) List(ctx context.Context, q Query, page, pageSize int) ([]*IPBlacklist, pagination.Page, error) {
	list, total, err := uc.repo.List(ctx, q, page, pageSize)
	if err != nil {
		return nil, pagination.Page{}, err
	}
	return list, pagination.Page{Page: page, PageSize: pageSize, Total: total}, nil
}

// IsBlocked 判定 IP 是否在当前黑名单中（带 30s 内存缓存）
func (uc *Usecase) IsBlocked(ip string) bool {
	uc.mu.RLock()
	if uc.cacheValid && time.Now().Before(uc.expireAt) {
		_, blocked := uc.cache[ip]
		uc.mu.RUnlock()
		return blocked
	}
	uc.mu.RUnlock()

	uc.mu.Lock()
	defer uc.mu.Unlock()
	// 双重检查
	if uc.cacheValid && time.Now().Before(uc.expireAt) {
		_, blocked := uc.cache[ip]
		return blocked
	}

	ips, err := uc.repo.ListActive(context.Background())
	if err != nil {
		// 查库失败时保守放行，避免误伤；失败状态会在下次请求重试
		uc.cache = map[string]struct{}{}
		uc.cacheValid = true
		uc.expireAt = time.Now().Add(30 * time.Second)
		return false
	}
	m := make(map[string]struct{}, len(ips))
	for _, v := range ips {
		m[v] = struct{}{}
	}
	uc.cache = m
	uc.cacheValid = true
	uc.expireAt = time.Now().Add(30 * time.Second)
	_, blocked := m[ip]
	return blocked
}

// invalidate 主动失效缓存（增删后立即生效）
func (uc *Usecase) invalidate() {
	uc.mu.Lock()
	uc.cacheValid = false
	uc.cache = nil
	uc.mu.Unlock()
}
