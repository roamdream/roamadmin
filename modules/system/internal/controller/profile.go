package controller

import (
	"context"

	"github.com/roamdream/roamadmin/modules/system/api"
	"github.com/roamdream/roamadmin/modules/system/internal/service"
)

// Profile 当前用户控制器（登录即可访问）。
type Profile struct{}

// NewProfile 创建当前用户控制器。
func NewProfile() *Profile { return &Profile{} }

// UserInfo 当前用户信息。
func (c *Profile) UserInfo(ctx context.Context, req *api.UserInfoReq) (*api.UserInfoRes, error) {
	info, err := service.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return &api.UserInfoRes{
		UserID:   info.User.Id,
		Username: info.User.Username,
		Nickname: info.User.Nickname,
		Email:    info.User.Email,
		Phone:    info.User.Phone,
		Avatar:   info.User.Avatar,
		Roles:    info.Roles,
		Perms:    info.Perms,
	}, nil
}

// Logout 用户登出。
func (c *Profile) Logout(ctx context.Context, req *api.LogoutReq) (*api.LogoutRes, error) {
	if err := service.Logout(ctx); err != nil {
		return nil, err
	}
	recordOper(ctx, "auth", "logout", 0, "")
	return &api.LogoutRes{Logout: true}, nil
}
