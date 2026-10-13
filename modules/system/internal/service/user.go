package service

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/kernel/utils"
	"github.com/roamdream/roamadmin/modules/system/internal/consts"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
)

// CreateUser 创建用户（雪花主键；用户名全局唯一，用分布式锁防并发重复；密码 bcrypt 加密）。
func CreateUser(ctx context.Context, user *model.User, roleIds []int64) (int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return 0, errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}

	// 分布式锁：同一用户名串行化，防止并发创建同名用户。
	unlock, err := ra.Lock().Lock(ctx, consts.CachePrefixUser+user.Username, 5*time.Second)
	if err != nil {
		return 0, errcode.Wrap(errcode.CodeInternal, err)
	}
	defer unlock()

	if exists, err := userExists(ctx, user.Username, 0); err != nil {
		return 0, err
	} else if exists {
		return 0, errcode.New(errcode.CodeDuplicate, "username already exists")
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return 0, errcode.Wrap(errcode.CodeInternal, err)
	}

	id := ra.DB().NextID()
	now := gtime.Now()
	err = ra.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(consts.TableUser).Data(g.Map{
			"id":         id,
			"tenant_id":  tenantID,
			"username":   user.Username,
			"password":   hash,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"phone":      user.Phone,
			"avatar":     user.Avatar,
			"status":     user.Status,
			"remark":     user.Remark,
			"created_at": now,
			"updated_at": now,
		}).Insert(); err != nil {
			return err
		}
		return insertUserRoles(ctx, tx, tenantID, id, roleIds)
	})
	if err != nil {
		return 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return id, nil
}

// UpdateUser 更新用户（用户名唯一校验同创建；角色关联整体替换）。
func UpdateUser(ctx context.Context, user *model.User, roleIds []int64) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}

	unlock, err := ra.Lock().Lock(ctx, consts.CachePrefixUser+user.Username, 5*time.Second)
	if err != nil {
		return errcode.Wrap(errcode.CodeInternal, err)
	}
	defer unlock()

	if exists, err := userExists(ctx, user.Username, user.Id); err != nil {
		return err
	} else if exists {
		return errcode.New(errcode.CodeDuplicate, "username already exists")
	}

	now := gtime.Now()
	err = ra.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(consts.TableUser).
			Where("id", user.Id).
			Where("tenant_id", tenantID).
			Data(g.Map{
				"username":   user.Username,
				"nickname":   user.Nickname,
				"email":      user.Email,
				"phone":      user.Phone,
				"avatar":     user.Avatar,
				"status":     user.Status,
				"remark":     user.Remark,
				"updated_at": now,
			}).Update(); err != nil {
			return err
		}
		if _, err := tx.Model(consts.TableUserRole).
			Where("user_id", user.Id).
			Where("tenant_id", tenantID).
			Delete(); err != nil {
			return err
		}
		return insertUserRoles(ctx, tx, tenantID, user.Id, roleIds)
	})
	if err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// ListUser 分页查询用户（按租户隔离）。
func ListUser(ctx context.Context, page, size int, username, nickname string, status *int) ([]*model.User, int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return nil, 0, errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 20
	}
	m := ra.DB().Model(consts.TableUser).Where("tenant_id", tenantID)
	if username != "" {
		m = m.WhereLike("username", "%"+username+"%")
	}
	if nickname != "" {
		m = m.WhereLike("nickname", "%"+nickname+"%")
	}
	if status != nil {
		m = m.Where("status", *status)
	}
	total, err := m.Count()
	if err != nil || total == 0 {
		return []*model.User{}, int64(total), err
	}
	var list []*model.User
	if err = m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return list, int64(total), nil
}

// GetUser 查询单个用户及其角色 ID 列表（按租户隔离）。
func GetUser(ctx context.Context, id int64) (*model.User, []int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return nil, nil, errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	var user model.User
	if err := ra.DB().Model(consts.TableUser).Where("id", id).Where("tenant_id", tenantID).Scan(&user); err != nil {
		return nil, nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if user.Id == 0 {
		return nil, nil, errcode.New(errcode.CodeNotFound, "user not found")
	}

	var roleRows []struct {
		RoleId int64 `orm:"role_id"`
	}
	if err := ra.DB().Model(consts.TableUserRole).
		Where("user_id", id).
		Where("tenant_id", tenantID).
		Fields("role_id").
		Scan(&roleRows); err != nil {
		return nil, nil, errcode.Wrap(errcode.CodeDB, err)
	}
	roleIds := make([]int64, 0, len(roleRows))
	for _, r := range roleRows {
		roleIds = append(roleIds, r.RoleId)
	}
	return &user, roleIds, nil
}

// DeleteUser 删除用户（同时清理其角色关联）。
func DeleteUser(ctx context.Context, id int64) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	err := ra.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(consts.TableUserRole).
			Where("user_id", id).
			Where("tenant_id", tenantID).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(consts.TableUser).
			Where("id", id).
			Where("tenant_id", tenantID).
			Delete(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// ResetUserPassword 重置用户密码。
func ResetUserPassword(ctx context.Context, id int64, password string) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return errcode.Wrap(errcode.CodeInternal, err)
	}
	if _, err := ra.DB().Model(consts.TableUser).
		Where("id", id).
		Where("tenant_id", tenantID).
		Data(g.Map{"password": hash, "updated_at": gtime.Now()}).
		Update(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// SetUserStatus 启用/禁用用户。
func SetUserStatus(ctx context.Context, id int64, status int) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	if _, err := ra.DB().Model(consts.TableUser).
		Where("id", id).
		Where("tenant_id", tenantID).
		Data(g.Map{"status": status, "updated_at": gtime.Now()}).
		Update(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// userExists 判断用户名是否已存在（全局唯一，excludeID 用于更新时排除自身）。
func userExists(ctx context.Context, username string, excludeID int64) (bool, error) {
	m := ra.DB().Model(consts.TableUser).Where("username", username)
	if excludeID > 0 {
		m = m.WhereNot("id", excludeID)
	}
	n, err := m.Count()
	return n > 0, err
}

// insertUserRoles 在事务内插入用户-角色关联。
func insertUserRoles(ctx context.Context, tx gdb.TX, tenantID, userID int64, roleIds []int64) error {
	for _, roleID := range roleIds {
		if _, err := tx.Model(consts.TableUserRole).Data(g.Map{
			"id":        ra.DB().NextID(),
			"user_id":   userID,
			"role_id":   roleID,
			"tenant_id": tenantID,
		}).Insert(); err != nil {
			return err
		}
	}
	return nil
}
