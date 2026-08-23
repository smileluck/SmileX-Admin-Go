// Package blacklist 黑名单限界上下文 —— 领域层
package blacklist

import "time"

// 哨兵错误
var (
	ErrInvalidIP   = errorString("invalid IP address")
	ErrIPExists    = errorString("IP already exists in blacklist")
	ErrNotFound    = errorString("blacklist entry not found")
	ErrSelfBan     = errorString("cannot ban your own IP")
	ErrInvalidExpire = errorString("expire time must be in the future")
)

type errorString string

func (e errorString) Error() string { return string(e) }

// IPBlacklist IP 黑名单实体
type IPBlacklist struct {
	ID          uint
	IP          string
	Reason      string
	ExpireAt    *time.Time // nil 表示永久
	CreatorID   uint
	CreatorName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Query 列表查询条件（零值为不限）
type Query struct {
	IP string // IP 前缀模糊
}
