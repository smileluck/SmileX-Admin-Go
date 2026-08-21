// Package auth 认证应用服务（薄用例编排，方法签名与 api/auth/v1 契约对应）
package auth

import (
	"context"
	"errors"

	bizauth "github.com/smilex/smilex-admin-gin/internal/biz/auth"
	bizcaptcha "github.com/smilex/smilex-admin-gin/internal/biz/captcha"
	"github.com/smilex/smilex-admin-gin/internal/biz/permission"
	"github.com/smilex/smilex-admin-gin/internal/biz/user"
)

type Service struct {
	uc      *bizauth.Usecase
	captcha *bizcaptcha.Usecase
}

func NewService(uc *bizauth.Usecase, captcha *bizcaptcha.Usecase) *Service {
	return &Service{uc: uc, captcha: captcha}
}

// LoginRequest 登录入参
type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

// RefreshRequest 刷新入参
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*bizauth.TokenPair, error) {
	return s.uc.Login(ctx, req.Username, req.Password, req.CaptchaID, req.CaptchaCode)
}

// CaptchaVO 图形验证码视图（image 为 PNG 的 base64，无 data: 前缀）
type CaptchaVO struct {
	CaptchaID    string `json:"captcha_id"`
	CaptchaImage string `json:"captcha_image"`
}

// GenerateCaptcha 生成一次性图形验证码
func (s *Service) GenerateCaptcha() (*CaptchaVO, error) {
	id, b64, err := s.captcha.Generate()
	if err != nil {
		return nil, err
	}
	return &CaptchaVO{CaptchaID: id, CaptchaImage: b64}, nil
}

func (s *Service) Refresh(ctx context.Context, req RefreshRequest) (*bizauth.TokenPair, error) {
	return s.uc.Refresh(ctx, req.RefreshToken)
}

// ProfileVO 个人信息视图
type ProfileVO struct {
	User        *user.User              `json:"user"`
	Permissions []*permission.Permission `json:"permissions"`
}

func (s *Service) Profile(ctx context.Context, userID uint) (*ProfileVO, error) {
	u, ps, err := s.uc.Profile(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}
	return &ProfileVO{User: u, Permissions: ps}, nil
}

// Authorize 供 RBAC 中间件调用
func (s *Service) Authorize(ctx context.Context, userID uint, method, path string) bool {
	return s.uc.Authorize(ctx, userID, method, path)
}

// ParseSubject 解析 access token
func (s *Service) ParseSubject(token string) (*bizauth.Subject, error) {
	return s.uc.ParseSubject(token)
}
