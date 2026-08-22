package auth

import (
	"context"
	"errors"

	"github.com/smilex/smilex-admin-gin/internal/biz/permission"
	"github.com/smilex/smilex-admin-gin/internal/biz/user"
)

var (
	// ErrInvalidCredentials 用户名或密码错误 / 用户被禁用
	ErrInvalidCredentials = errors.New("invalid username or password")
	// ErrDisabledAccount 账号被禁用
	ErrDisabledAccount = errors.New("account disabled")
	// ErrCaptcha 验证码缺失/过期/错误
	ErrCaptcha = errors.New("captcha code invalid")
)

// UserStore 认证上下文依赖的用户存取接口（只依赖所需最小方法，而非整个 user.Repo）
type UserStore interface {
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	FindByID(ctx context.Context, id uint) (*user.User, error)
	// FindByIDWithPassword 携带密码哈希返回（本人改密校验旧密码用）
	FindByIDWithPassword(ctx context.Context, id uint) (*user.User, error)
	// Update 本人更新基础资料（仓储层仅写昵称/邮箱/状态）
	Update(ctx context.Context, u *user.User) error
	// UpdatePassword 单独落库密码哈希
	UpdatePassword(ctx context.Context, id uint, passwordHash string) error
}

// RoleNameReader 角色名读取接口（个人中心展示角色名）
type RoleNameReader interface {
	FindNamesByIDs(ctx context.Context, ids []uint) ([]string, error)
}

// PermissionReader 权限读取接口（跨上下文走接口）
type PermissionReader interface {
	FindByUserID(ctx context.Context, userID uint) ([]*permission.Permission, error)
}

// CaptchaVerifier 登录验证码校验接口（由 captcha 上下文实现，跨上下文走接口）
type CaptchaVerifier interface {
	Verify(id, answer string) bool
}

// Usecase 认证领域用例
type Usecase struct {
	users   UserStore
	roles   RoleNameReader
	perms   PermissionReader
	tokens  TokenIssuer
	captcha CaptchaVerifier
}

func NewUsecase(users UserStore, roles RoleNameReader, perms PermissionReader, tokens TokenIssuer, captcha CaptchaVerifier) *Usecase {
	return &Usecase{users: users, roles: roles, perms: perms, tokens: tokens, captcha: captcha}
}

// Login 登录签发令牌（先校验一次性图形验证码）
func (uc *Usecase) Login(ctx context.Context, username, password, captchaID, captchaCode string) (*TokenPair, error) {
	if !uc.captcha.Verify(captchaID, captchaCode) {
		return nil, ErrCaptcha
	}
	u, err := uc.users.FindByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if !u.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}
	if !u.Enabled() {
		return nil, ErrDisabledAccount
	}
	access, expiresAt, err := uc.tokens.IssueAccessToken(Subject{UserID: u.ID, Username: u.Username})
	if err != nil {
		return nil, err
	}
	refresh, err := uc.tokens.IssueRefreshToken(Subject{UserID: u.ID, Username: u.Username})
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh, ExpiresAt: expiresAt}, nil
}

// Refresh 刷新令牌
func (uc *Usecase) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	s, err := uc.tokens.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	u, err := uc.users.FindByID(ctx, s.UserID)
	if err != nil || !u.Enabled() {
		return nil, ErrInvalidCredentials
	}
	access, expiresAt, err := uc.tokens.IssueAccessToken(Subject{UserID: u.ID, Username: u.Username})
	if err != nil {
		return nil, err
	}
	refresh, err := uc.tokens.IssueRefreshToken(Subject{UserID: u.ID, Username: u.Username})
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh, ExpiresAt: expiresAt}, nil
}

// Profile 当前用户信息（含角色名与权限）
func (uc *Usecase) Profile(ctx context.Context, userID uint) (*Profile, error) {
	u, err := uc.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	ps, err := uc.perms.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	// 角色名仅用于展示，读取失败时降级为空列表，不阻断个人中心
	roleNames, err := uc.roles.FindNamesByIDs(ctx, u.RoleIDs)
	if err != nil {
		roleNames = nil
	}
	return &Profile{User: u, RoleNames: roleNames, Permissions: ps}, nil
}

// UpdateProfile 本人更新基础资料（昵称/邮箱；不经过管理员接口的超管保护）
func (uc *Usecase) UpdateProfile(ctx context.Context, userID uint, nickname, email string) (*user.User, error) {
	u, err := uc.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if nickname != "" {
		u.Nickname = nickname
	}
	if email != "" {
		u.Email = email
	}
	if err := uc.users.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// ChangePassword 本人修改密码：先校验旧密码，再落库新哈希
func (uc *Usecase) ChangePassword(ctx context.Context, userID uint, oldPlain, newPlain string) error {
	u, err := uc.users.FindByIDWithPassword(ctx, userID)
	if err != nil {
		return err
	}
	if !u.CheckPassword(oldPlain) {
		return ErrInvalidCredentials
	}
	if err := u.SetPassword(newPlain); err != nil {
		return err
	}
	return uc.users.UpdatePassword(ctx, userID, string(u.Password))
}

// ParseSubject 解析 access token
func (uc *Usecase) ParseSubject(token string) (*Subject, error) {
	return uc.tokens.ParseAccessToken(token)
}

// Authorize RBAC 接口鉴权
func (uc *Usecase) Authorize(ctx context.Context, userID uint, method, path string) bool {
	ps, err := uc.perms.FindByUserID(ctx, userID)
	if err != nil {
		return false
	}
	for _, p := range ps {
		if p.Match(method, path) {
			return true
		}
	}
	return false
}
