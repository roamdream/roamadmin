// Package auth 认证鉴权：JWT + RBAC + 多租户数据隔离。
//
// 内核负责令牌签发/校验、用户上下文注入、角色权限判断与租户提取；
// 具体的用户/角色/权限/租户数据由 system 模块落地，登录时组装成 Claims。
package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/roamdream/roamadmin/kernel/config"
	"github.com/roamdream/roamadmin/kernel/errcode"
)

// Claims JWT 载荷：身份 + 角色权限 + 租户。
type Claims struct {
	UserID   int64    `json:"uid"`
	TenantID int64    `json:"tid"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Perms    []string `json:"perms"`
	jwt.RegisteredClaims
}

// HasRole 判断是否具备某角色。
func (c *Claims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasPermission 判断是否具备某权限码；拥有 "*" 通配视为超管全通过。
func (c *Claims) HasPermission(perm string) bool {
	for _, p := range c.Perms {
		if p == "*" || p == perm {
			return true
		}
	}
	return false
}

// IsSuper 是否超级管理员（角色 super 或权限通配）。
func (c *Claims) IsSuper() bool {
	return c.HasRole("super") || c.HasPermission("*")
}

// Manager JWT 令牌管理器。
type Manager struct {
	secret []byte
	issuer string
	expire time.Duration
}

// New 创建 JWT 管理器。
func New(secret, issuer string, expire time.Duration) *Manager {
	return &Manager{secret: []byte(secret), issuer: issuer, expire: expire}
}

// Default 从配置读取 JWT 参数（auth.jwt.*）。
func Default() *Manager {
	ctx := context.Background()
	cfg := config.Default()
	secret := cfg.MustString(ctx, "auth.jwt.secret", "roamadmin-default-secret")
	issuer := cfg.MustString(ctx, "auth.jwt.issuer", "roamadmin")
	expire := time.Duration(cfg.MustInt(ctx, "auth.jwt.expireSeconds", 7200)) * time.Second
	return New(secret, issuer, expire)
}

// IssueToken 签发令牌。
func (m *Manager) IssueToken(c *Claims) (string, error) {
	now := time.Now()
	c.Issuer = m.issuer
	c.IssuedAt = jwt.NewNumericDate(now)
	c.ExpiresAt = jwt.NewNumericDate(now.Add(m.expire))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(m.secret)
}

// ParseToken 校验并解析令牌。
func (m *Manager) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errcode.New(errcode.CodeUnauthorized, "unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, errcode.Wrap(errcode.CodeUnauthorized, err)
	}
	if c, ok := token.Claims.(*Claims); ok && token.Valid {
		return c, nil
	}
	return nil, errcode.New(errcode.CodeUnauthorized, "invalid token")
}

// —— 用户 / 租户上下文 ——

type userCtxKey struct{}

// WithUser 将用户注入 context。
func WithUser(ctx context.Context, c *Claims) context.Context {
	return context.WithValue(ctx, userCtxKey{}, c)
}

// UserFrom 从 context 取出当前用户。
func UserFrom(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(userCtxKey{}).(*Claims)
	return c, ok && c != nil
}

// TenantFrom 从 context 取出当前租户 ID。
func TenantFrom(ctx context.Context) (int64, bool) {
	c, ok := UserFrom(ctx)
	if !ok {
		return 0, false
	}
	return c.TenantID, true
}
