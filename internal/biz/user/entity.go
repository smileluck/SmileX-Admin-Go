// Package user 用户限界上下文 —— 领域层
package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Status 用户状态值对象
type Status int

const (
	StatusDisabled Status = 0
	StatusEnabled  Status = 1
)

// Password 密码值对象（bcrypt 封装）
type Password string

// NewPassword 从明文创建密码
func NewPassword(plain string) (Password, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return Password(b), nil
}

// Verify 校验明文密码
func (p Password) Verify(plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p), []byte(plain)) == nil
}

// User 用户聚合根（纯 Go，不依赖任何框架）
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  Password  `json:"-"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Status    Status    `json:"status"`
	RoleIDs   []uint    `json:"role_ids"` // 聚合内的角色关联
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Enabled 是否启用
func (u *User) Enabled() bool { return u.Status == StatusEnabled }

// SetPassword 设置密码
func (u *User) SetPassword(plain string) error {
	p, err := NewPassword(plain)
	if err != nil {
		return err
	}
	u.Password = p
	return nil
}

// CheckPassword 校验密码
func (u *User) CheckPassword(plain string) bool { return u.Password.Verify(plain) }
