package controller

import (
	"context"

	"github.com/roamdream/roamadmin/kernel/utils"
	"github.com/roamdream/roamadmin/modules/system/api"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
	"github.com/roamdream/roamadmin/modules/system/internal/service"
)

// User 用户控制器（用户 CRUD）。
type User struct{}

// NewUser 创建用户控制器。
func NewUser() *User { return &User{} }

// Create 创建用户。
func (c *User) Create(ctx context.Context, req *api.UserCreateReq) (*api.UserCreateRes, error) {
	id, err := service.CreateUser(ctx, &model.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   req.Status,
		Remark:   req.Remark,
	}, req.RoleIds)
	writeOper(ctx, "user", "create", err)
	if err != nil {
		return nil, err
	}
	return &api.UserCreateRes{Id: id}, nil
}

// List 分页查询用户。
func (c *User) List(ctx context.Context, req *api.UserListReq) (*api.UserListRes, error) {
	list, total, err := service.ListUser(ctx, req.Page, req.Size, req.Username, req.Nickname, req.Status)
	if err != nil {
		return nil, err
	}
	out := make([]*api.UserItem, 0, len(list))
	for _, it := range list {
		out = append(out, &api.UserItem{
			Id:        it.Id,
			TenantId:  it.TenantId,
			Username:  it.Username,
			Nickname:  it.Nickname,
			Email:     it.Email,
			Phone:     it.Phone,
			Avatar:    it.Avatar,
			Status:    it.Status,
			Remark:    it.Remark,
			CreatedAt: utils.TimeRFC3339(it.CreatedAt),
		})
	}
	return &api.UserListRes{List: out, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Get 查询用户。
func (c *User) Get(ctx context.Context, req *api.UserGetReq) (*api.UserGetRes, error) {
	user, roleIds, err := service.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.UserGetRes{
		User: &api.UserItem{
			Id:        user.Id,
			TenantId:  user.TenantId,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			Status:    user.Status,
			Remark:    user.Remark,
			CreatedAt: utils.TimeRFC3339(user.CreatedAt),
		},
		RoleIds: roleIds,
	}, nil
}

// Update 更新用户。
func (c *User) Update(ctx context.Context, req *api.UserUpdateReq) (*api.UserUpdateRes, error) {
	err := service.UpdateUser(ctx, &model.User{
		Id:       req.Id,
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   req.Status,
		Remark:   req.Remark,
	}, req.RoleIds)
	writeOper(ctx, "user", "update", err)
	if err != nil {
		return nil, err
	}
	return &api.UserUpdateRes{Updated: true}, nil
}

// Delete 删除用户。
func (c *User) Delete(ctx context.Context, req *api.UserDeleteReq) (*api.UserDeleteRes, error) {
	err := service.DeleteUser(ctx, req.Id)
	writeOper(ctx, "user", "delete", err)
	if err != nil {
		return nil, err
	}
	return &api.UserDeleteRes{Deleted: true}, nil
}

// ResetPassword 重置用户密码。
func (c *User) ResetPassword(ctx context.Context, req *api.UserResetPasswordReq) (*api.UserResetPasswordRes, error) {
	err := service.ResetUserPassword(ctx, req.Id, req.Password)
	writeOper(ctx, "user", "reset password", err)
	if err != nil {
		return nil, err
	}
	return &api.UserResetPasswordRes{Reset: true}, nil
}

// SetStatus 启用/禁用用户。
func (c *User) SetStatus(ctx context.Context, req *api.UserSetStatusReq) (*api.UserSetStatusRes, error) {
	err := service.SetUserStatus(ctx, req.Id, req.Status)
	writeOper(ctx, "user", "enable or disable", err)
	if err != nil {
		return nil, err
	}
	return &api.UserSetStatusRes{Updated: true}, nil
}
