// Package service system 模块业务逻辑。
package service

import (
	"context"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/kernel/utils"
	"github.com/roamdream/roamadmin/modules/system/internal/consts"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
)

// Login 用户登录：校验账号密码与状态，组装角色权限后签发 JWT。
func Login(ctx context.Context, username, password string) (string, *model.UserInfo, error) {
	var user model.User
	if err := ra.DB().Model(consts.TableUser).Where("username", username).Scan(&user); err != nil {
		return "", nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if user.Id == 0 || !utils.VerifyPassword(user.Password, password) {
		return "", nil, errcode.New(errcode.CodeLoginFailed, "用户名或密码错误")
	}
	if user.Status != 1 {
		return "", nil, errcode.New(errcode.CodeUserDisabled, "用户已被禁用")
	}

	roles, perms, err := loadRolesPerms(ctx, user.TenantId, user.Id)
	if err != nil {
		return "", nil, err
	}

	token, err := ra.Auth().IssueToken(&auth.Claims{
		UserID:   user.Id,
		TenantID: user.TenantId,
		Username: user.Username,
		Roles:    roles,
		Perms:    perms,
	})
	if err != nil {
		return "", nil, errcode.Wrap(errcode.CodeInternal, err)
	}

	return token, &model.UserInfo{User: &user, Roles: roles, Perms: perms}, nil
}

// GetUserInfo 查询当前登录用户信息（基本信息 + 角色 + 权限）。
func GetUserInfo(ctx context.Context) (*model.UserInfo, error) {
	c, ok := auth.UserFrom(ctx)
	if !ok {
		return nil, errcode.New(errcode.CodeUnauthorized, "未登录")
	}
	var user model.User
	if err := ra.DB().Model(consts.TableUser).Where("id", c.UserID).Where("tenant_id", c.TenantID).Scan(&user); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if user.Id == 0 {
		return nil, errcode.New(errcode.CodeNotFound, "用户不存在")
	}
	roles, perms, err := loadRolesPerms(ctx, c.TenantID, c.UserID)
	if err != nil {
		return nil, err
	}
	return &model.UserInfo{User: &user, Roles: roles, Perms: perms}, nil
}

// Logout 用户登出：JWT 无状态，令牌由客户端丢弃，此处仅作占位（日志留痕在控制器完成）。
func Logout(ctx context.Context) error { return nil }

// loadRolesPerms 加载用户的角色编码与权限码；超级管理员角色（super）附加通配 "*"。
func loadRolesPerms(ctx context.Context, tenantID, userID int64) ([]string, []string, error) {
	var roleRows []struct {
		Code string `orm:"code"`
	}
	if err := ra.DB().Model(consts.TableUserRole+" u").
		LeftJoin(consts.TableRole+" r", "r.id = u.role_id").
		Where("u.user_id", userID).
		Where("u.tenant_id", tenantID).
		Fields("r.code").
		Scan(&roleRows); err != nil {
		return nil, nil, errcode.Wrap(errcode.CodeDB, err)
	}
	roles := make([]string, 0, len(roleRows))
	for _, r := range roleRows {
		roles = append(roles, r.Code)
	}

	var permRows []struct {
		Perm string `orm:"perm"`
	}
	if err := ra.DB().Model(consts.TableUserRole+" u").
		LeftJoin(consts.TableRoleMenu+" rm", "rm.role_id = u.role_id").
		LeftJoin(consts.TableMenu+" m", "m.id = rm.menu_id").
		Where("u.user_id", userID).
		Where("u.tenant_id", tenantID).
		Fields("m.perm").
		Scan(&permRows); err != nil {
		return nil, nil, errcode.Wrap(errcode.CodeDB, err)
	}
	permSet := make(map[string]struct{}, len(permRows))
	for _, p := range permRows {
		if p.Perm != "" {
			permSet[p.Perm] = struct{}{}
		}
	}
	for _, role := range roles {
		if role == "super" {
			permSet["*"] = struct{}{}
			break
		}
	}
	perms := make([]string, 0, len(permSet))
	for p := range permSet {
		perms = append(perms, p)
	}
	return roles, perms, nil
}
