// Package model GORM 持久化对象（PO）与领域实体的转换。
// PO 只在 data 层出现，biz 层不可见。
package model

import (
	"time"

	"github.com/smilex/smilex-admin-gin/internal/biz/blacklist"
	"github.com/smilex/smilex-admin-gin/internal/biz/export"
	"github.com/smilex/smilex-admin-gin/internal/biz/file"
	"github.com/smilex/smilex-admin-gin/internal/biz/log"
	"github.com/smilex/smilex-admin-gin/internal/biz/permission"
	"github.com/smilex/smilex-admin-gin/internal/biz/role"
	"github.com/smilex/smilex-admin-gin/internal/biz/user"
	"gorm.io/gorm"
)

// UserPO 用户表
type UserPO struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"size:64;uniqueIndex"`
	Password  string `gorm:"size:128"`
	Nickname  string `gorm:"size:64"`
	Email     string `gorm:"size:128"`
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserPO) TableName() string { return "users" }

// RolePO 角色表
type RolePO struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:64"`
	Remark    string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (RolePO) TableName() string { return "roles" }

// PermissionPO 权限表
type PermissionPO struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:64"`
	Code      string `gorm:"size:64;uniqueIndex"`
	Type      string `gorm:"size:16"` // dir | menu | button（api 已废弃，启动迁移自动转为 button；有菜单子级的 menu 迁移为 dir）
	Method    string `gorm:"size:16"`
	Path      string `gorm:"size:255"`
	ParentID  uint
	Icon      string `gorm:"size:512"`
	Sort      int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (PermissionPO) TableName() string { return "permissions" }

// UserRolePO 用户-角色关联
type UserRolePO struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

func (UserRolePO) TableName() string { return "user_roles" }

// RolePermissionPO 角色-权限关联
type RolePermissionPO struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

func (RolePermissionPO) TableName() string { return "role_permissions" }

// LoginLogPO 登录日志表（追加型流水：无软删，清空/保留期清理均为物理删除）
type LoginLogPO struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"size:64;index"` // 尝试登录的用户名（可能不存在）
	IP        string    `gorm:"size:64;index"`
	UserAgent string    `gorm:"size:255"`
	Device    string    `gorm:"size:16"`  // web / app
	Status    int       `gorm:"index"`    // 1 成功 0 失败
	Msg       string    `gorm:"size:255"` // 失败原因（成功为空）
	CreatedAt time.Time `gorm:"index"`    // 登录时间
}

func (LoginLogPO) TableName() string { return "login_logs" }

// OperationLogPO 操作日志表（写请求审计流水：无软删，清空/保留期清理均为物理删除）
type OperationLogPO struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"index"`  // 操作人（JWT 校验失败被拒时为 0）
	Username   string    `gorm:"size:64;index"` // 操作人用户名快照
	Method     string    `gorm:"size:8;index"`
	Path       string    `gorm:"size:255"` // 实际请求路径（含资源 ID 与 query）
	Route      string    `gorm:"size:128"` // 路由模板（如 /api/v1/users/:id）
	Action     string    `gorm:"size:64"`  // 中文动作名（如「新增用户」）
	Params     string    `gorm:"type:text"` // 请求参数摘要（敏感字段脱敏、超长截断）
	IP         string    `gorm:"size:64"`
	UserAgent  string    `gorm:"size:255"`
	StatusCode int       // 响应状态码
	LatencyMs  int       // 耗时（毫秒）
	CreatedAt  time.Time `gorm:"index"`
}

func (OperationLogPO) TableName() string { return "operation_logs" }

// FilePO 文件元数据表（对象本体在 driver 对应的存储后端；driver 落库保证后端升级后旧文件仍可访问）
type FilePO struct {
	ID           uint   `gorm:"primaryKey"`
	Driver       string `gorm:"size:16;index"`   // local | oss | cos | tos | minio
	ObjectKey    string `gorm:"size:512;uniqueIndex"` // 服务端生成的对象 key
	Name         string `gorm:"size:255"`        // 原始文件名
	Ext          string `gorm:"size:16;index"`   // 扩展名（小写，不含点）
	Size         int64
	ContentType  string `gorm:"size:128"`
	UploaderID   uint   `gorm:"index"`
	UploaderName string `gorm:"size:64"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (FilePO) TableName() string { return "files" }

// IPBlacklistPO IP 黑名单表（管理员手工维护的持久化封禁；软删即解封留痕）
type IPBlacklistPO struct {
	ID          uint       `gorm:"primaryKey"`
	IP          string     `gorm:"size:64;uniqueIndex"` // 仅单个 IP（不支持 CIDR），归一化后存储
	Reason      string     `gorm:"size:255"`            // 封禁原因（选填）
	ExpireAt    *time.Time // 过期时间（NULL 为永久封禁，到期惰性放行）
	CreatorID   uint       // 操作人
	CreatorName string     `gorm:"size:64"` // 操作人用户名快照
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (IPBlacklistPO) TableName() string { return "ip_blacklist" }

// ExportRecordPO 异步导出任务记录表（产物本体在 driver 对应的存储后端；
// 追加型流水：无软删，保留期清理与手动删除均为物理删除）
type ExportRecordPO struct {
	ID         uint       `gorm:"primaryKey"`
	UserID     uint       `gorm:"index"` // 任务归属用户
	Biz        string     `gorm:"size:32"` // 业务类型（user / login_log / op_log）
	Name       string     `gorm:"size:255"` // 展示名（兼作下载文件名）
	Params     string     `gorm:"type:text"` // 查询条件快照（JSON）
	Driver     string     `gorm:"size:16"` // 产物落库时的存储后端
	ObjectKey  string     `gorm:"size:512"` // 产物对象 key
	Size       int64      // 产物字节数（含 BOM）
	Rows       int        // 已导出数据行数（不含表头）
	Status     string     `gorm:"size:16;index"` // pending | running | done | failed
	Truncated  bool       // 触及大小/行数上限被截断
	Error      string     `gorm:"size:512"` // 失败原因（成功为空）
	CreatedAt  time.Time  `gorm:"index"`
	FinishedAt *time.Time // 完成/失败时间（未结束为 NULL）
}

func (ExportRecordPO) TableName() string { return "export_records" }

// ---- 转换器 ----

func UserToPO(u *user.User) *UserPO {
	return &UserPO{
		ID: u.ID, Username: u.Username, Password: string(u.Password),
		Nickname: u.Nickname, Email: u.Email, Status: int(u.Status),
	}
}

func UserFromPO(p *UserPO) *user.User {
	return &user.User{
		ID: p.ID, Username: p.Username, Password: user.Password(p.Password),
		Nickname: p.Nickname, Email: p.Email, Status: user.Status(p.Status),
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func RoleToPO(r *role.Role) *RolePO {
	return &RolePO{ID: r.ID, Name: r.Name, Remark: r.Remark}
}

func RoleFromPO(p *RolePO) *role.Role {
	return &role.Role{ID: p.ID, Name: p.Name, Remark: p.Remark, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

func PermissionToPO(m *permission.Permission) *PermissionPO {
	return &PermissionPO{
		ID: m.ID, Name: m.Name, Code: m.Code, Type: string(m.Type),
		Method: m.Method, Path: m.Path, ParentID: m.ParentID,
		Icon: m.Icon, Sort: m.Sort,
	}
}

func PermissionFromPO(p *PermissionPO) *permission.Permission {
	return &permission.Permission{
		ID: p.ID, Name: p.Name, Code: p.Code, Type: permission.Type(p.Type),
		Method: p.Method, Path: p.Path, ParentID: p.ParentID,
		Icon: p.Icon, Sort: p.Sort,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func LoginLogToPO(l *log.LoginLog) *LoginLogPO {
	return &LoginLogPO{
		ID: l.ID, Username: l.Username, IP: l.IP, UserAgent: l.UserAgent,
		Device: l.Device, Status: l.Status, Msg: l.Msg, CreatedAt: l.CreatedAt,
	}
}

func LoginLogFromPO(p *LoginLogPO) *log.LoginLog {
	return &log.LoginLog{
		ID: p.ID, Username: p.Username, IP: p.IP, UserAgent: p.UserAgent,
		Device: p.Device, Status: p.Status, Msg: p.Msg, CreatedAt: p.CreatedAt,
	}
}

func OperationLogToPO(o *log.OperationLog) *OperationLogPO {
	return &OperationLogPO{
		ID: o.ID, UserID: o.UserID, Username: o.Username, Method: o.Method,
		Path: o.Path, Route: o.Route, Action: o.Action, Params: o.Params,
		IP: o.IP, UserAgent: o.UserAgent, StatusCode: o.StatusCode,
		LatencyMs: o.LatencyMs, CreatedAt: o.CreatedAt,
	}
}

func OperationLogFromPO(p *OperationLogPO) *log.OperationLog {
	return &log.OperationLog{
		ID: p.ID, UserID: p.UserID, Username: p.Username, Method: p.Method,
		Path: p.Path, Route: p.Route, Action: p.Action, Params: p.Params,
		IP: p.IP, UserAgent: p.UserAgent, StatusCode: p.StatusCode,
		LatencyMs: p.LatencyMs, CreatedAt: p.CreatedAt,
	}
}

func FileToPO(f *file.File) *FilePO {
	return &FilePO{
		ID: f.ID, Driver: f.Driver, ObjectKey: f.ObjectKey, Name: f.Name,
		Ext: f.Ext, Size: f.Size, ContentType: f.ContentType,
		UploaderID: f.UploaderID, UploaderName: f.UploaderName,
	}
}

func FileFromPO(p *FilePO) *file.File {
	return &file.File{
		ID: p.ID, Driver: p.Driver, ObjectKey: p.ObjectKey, Name: p.Name,
		Ext: p.Ext, Size: p.Size, ContentType: p.ContentType,
		UploaderID: p.UploaderID, UploaderName: p.UploaderName,
		CreatedAt: p.CreatedAt,
	}
}

func ExportRecordToPO(r *export.ExportRecord) *ExportRecordPO {
	return &ExportRecordPO{
		ID: r.ID, UserID: r.UserID, Biz: r.Biz, Name: r.Name, Params: r.Params,
		Driver: r.Driver, ObjectKey: r.ObjectKey, Size: r.Size, Rows: r.Rows,
		Status: r.Status, Truncated: r.Truncated, Error: r.Error,
		CreatedAt: r.CreatedAt, FinishedAt: r.FinishedAt,
	}
}

func ExportRecordFromPO(p *ExportRecordPO) *export.ExportRecord {
	return &export.ExportRecord{
		ID: p.ID, UserID: p.UserID, Biz: p.Biz, Name: p.Name, Params: p.Params,
		Driver: p.Driver, ObjectKey: p.ObjectKey, Size: p.Size, Rows: p.Rows,
		Status: p.Status, Truncated: p.Truncated, Error: p.Error,
		CreatedAt: p.CreatedAt, FinishedAt: p.FinishedAt,
	}
}

func IPBlacklistToPO(b *blacklist.IPBlacklist) *IPBlacklistPO {
	return &IPBlacklistPO{
		ID: b.ID, IP: b.IP, Reason: b.Reason, ExpireAt: b.ExpireAt,
		CreatorID: b.CreatorID, CreatorName: b.CreatorName,
	}
}

func IPBlacklistFromPO(p *IPBlacklistPO) *blacklist.IPBlacklist {
	return &blacklist.IPBlacklist{
		ID: p.ID, IP: p.IP, Reason: p.Reason, ExpireAt: p.ExpireAt,
		CreatorID: p.CreatorID, CreatorName: p.CreatorName,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}
