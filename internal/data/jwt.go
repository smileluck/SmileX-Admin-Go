package data

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smilex/smilex-admin-gin/internal/biz/auth"
	"github.com/smilex/smilex-admin-gin/internal/conf"
)

type claims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type jwtIssuer struct {
	secret       []byte
	issuer       string
	expireHours  int
	refreshHours int
}

// NewJWTIssuer JWT 令牌签发实现
func NewJWTIssuer(c *conf.Bootstrap) auth.TokenIssuer {
	return &jwtIssuer{
		secret:       []byte(c.JWT.Secret),
		issuer:       c.JWT.Issuer,
		expireHours:  c.JWT.ExpireHours,
		refreshHours: c.JWT.RefreshHours,
	}
}

func (j *jwtIssuer) newClaims(s auth.Subject, ttl time.Duration) *claims {
	now := time.Now()
	return &claims{
		UserID:   s.UserID,
		Username: s.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   s.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
}

func (j *jwtIssuer) sign(c *claims) (string, time.Time, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := t.SignedString(j.secret)
	return token, c.ExpiresAt.Time, err
}

func (j *jwtIssuer) IssueAccessToken(s auth.Subject) (string, time.Time, error) {
	return j.sign(j.newClaims(s, time.Duration(j.expireHours)*time.Hour))
}

func (j *jwtIssuer) IssueRefreshToken(s auth.Subject) (string, error) {
	token, _, err := j.sign(j.newClaims(s, time.Duration(j.refreshHours)*time.Hour))
	return token, err
}

func (j *jwtIssuer) parse(token string) (*auth.Subject, error) {
	var c claims
	_, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	return &auth.Subject{UserID: c.UserID, Username: c.Username}, nil
}

func (j *jwtIssuer) ParseAccessToken(token string) (*auth.Subject, error) { return j.parse(token) }
func (j *jwtIssuer) ParseRefreshToken(token string) (*auth.Subject, error) { return j.parse(token) }
