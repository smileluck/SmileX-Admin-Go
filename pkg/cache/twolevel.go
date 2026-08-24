// Package cache 二级缓存：L1 进程内存 + L2 Redis。
//
// 读路径 L1 → L2 → miss（由调用方回源并 Set 回填两级）；
// L2 可通过开关关闭（l2Enabled=false 时退化为纯 L1，多实例间不共享）。
// 定位为短 TTL 读缓存：不保证强一致，变更一致性由 TTL 兜底。
package cache

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/smilex/smilex-admin-gin/pkg/logger"
	"go.uber.org/zap"
)

// l1MaxEntries L1 条目数上限，超过时惰性清理过期条目
const l1MaxEntries = 4096

type l1Entry struct {
	val      string
	expireAt time.Time
}

// TwoLevel L1（内存）+ L2（Redis）二级缓存
type TwoLevel struct {
	rdb    *redis.Client // nil 表示 L2 关闭
	prefix string        // Redis key 前缀
	l1TTL  time.Duration
	l2TTL  time.Duration

	mu sync.RWMutex
	l1 map[string]l1Entry
}

// NewTwoLevel 构造二级缓存；l2Enabled=false 或 rdb=nil 时仅启用 L1
func NewTwoLevel(rdb *redis.Client, prefix string, l1TTL, l2TTL time.Duration, l2Enabled bool) *TwoLevel {
	if !l2Enabled {
		rdb = nil
	}
	return &TwoLevel{rdb: rdb, prefix: prefix, l1TTL: l1TTL, l2TTL: l2TTL, l1: map[string]l1Entry{}}
}

// Get 读取缓存；L2 命中时回填 L1
func (c *TwoLevel) Get(ctx context.Context, key string) (string, bool) {
	c.mu.RLock()
	e, ok := c.l1[key]
	c.mu.RUnlock()
	if ok {
		if time.Now().Before(e.expireAt) {
			return e.val, true
		}
		c.mu.Lock()
		delete(c.l1, key)
		c.mu.Unlock()
	}
	if c.rdb == nil {
		return "", false
	}
	val, err := c.rdb.Get(ctx, c.prefix+key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Warn("l2 cache get failed", zap.String("key", key), zap.Error(err))
		}
		return "", false
	}
	c.setL1(key, val)
	return val, true
}

// Set 写入两级缓存（L2 失败仅记日志，L1 仍生效）
func (c *TwoLevel) Set(ctx context.Context, key, val string) {
	c.setL1(key, val)
	if c.rdb == nil {
		return
	}
	if err := c.rdb.Set(ctx, c.prefix+key, val, c.l2TTL).Err(); err != nil {
		logger.Warn("l2 cache set failed", zap.String("key", key), zap.Error(err))
	}
}

// Del 双删（变更主动失效用）
func (c *TwoLevel) Del(ctx context.Context, key string) {
	c.mu.Lock()
	delete(c.l1, key)
	c.mu.Unlock()
	if c.rdb == nil {
		return
	}
	if err := c.rdb.Del(ctx, c.prefix+key).Err(); err != nil {
		logger.Warn("l2 cache del failed", zap.String("key", key), zap.Error(err))
	}
}

func (c *TwoLevel) setL1(key, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.l1) > l1MaxEntries {
		now := time.Now()
		for k, e := range c.l1 {
			if now.After(e.expireAt) {
				delete(c.l1, k)
			}
		}
	}
	c.l1[key] = l1Entry{val: val, expireAt: time.Now().Add(c.l1TTL)}
}
