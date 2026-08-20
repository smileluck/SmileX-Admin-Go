// Package role 角色限界上下文 —— 领域层
package role

import "time"

// Role 角色聚合根
type Role struct {
	ID            uint
	Name          string
	Code          string
	Remark        string
	PermissionIDs []uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
