// Package controller system 模块 HTTP 控制器。
package controller

import (
	"context"

	"github.com/roamdream/roamadmin/modules/system/api"
	"github.com/roamdream/roamadmin/modules/system/internal/service"
)

// Auth 认证控制器（登录，公开接口）。
type Auth struct{}

// NewAuth 创建认证控制器。
func NewAuth() *Auth { return &Auth{} }

// Login 用户登录。
func (c *Auth) Login(ctx context.Context, req *api.LoginReq) (*api.LoginRes, error) {
	token, info, err := service.Login(ctx, req.Username, req.Password)
	if err != nil {
		recordLogin(ctx, 1, err.Error(), 0, 0, req.Username)
		return nil, err
	}
	recordLogin(ctx, 0, "", info.User.Id, info.User.TenantId, info.User.Username)
	return &api.LoginRes{
		Token:  token,
		UserID: info.User.Id,
		Name:   info.User.Nickname,
		Roles:  info.Roles,
		Perms:  info.Perms,
	}, nil
}
