package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/system/internal/consts"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
)

// CreateRole 创建角色（角色编码在同一租户内唯一，分布式锁防并发重复）。
func CreateRole(ctx context.Context, role *model.Role) (int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return 0, errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}

	lockKey := fmt.Sprintf("%s%d:%s", consts.CachePrefixRole, tenantID, role.Code)
	unlock, err := ra.Lock().Lock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return 0, errcode.Wrap(errcode.CodeInternal, err)
	}
	defer unlock()

	if exists, err := roleCodeExists(ctx, tenantID, role.Code, 0); err != nil {
		return 0, err
	} else if exists {
		return 0, errcode.New(errcode.CodeDuplicate, "角色编码已存在")
	}

	id := ra.DB().NextID()
	now := gtime.Now()
	if _, err := ra.DB().Model(consts.TableRole).Data(g.Map{
		"id":         id,
		"tenant_id":  tenantID,
		"name":       role.Name,
		"code":       role.Code,
		"status":     role.Status,
		"remark":     role.Remark,
		"created_at": now,
		"updated_at": now,
	}).Insert(); err != nil {
		return 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return id, nil
}

// ListRole 分页查询角色（按租户隔离）。
func ListRole(ctx context.Context, page, size int, name string) ([]*model.Role, int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return nil, 0, errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 20
	}
	m := ra.DB().Model(consts.TableRole).Where("tenant_id", tenantID)
	if name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}
	total, err := m.Count()
	if err != nil || total == 0 {
		return []*model.Role{}, int64(total), err
	}
	var list []*model.Role
	if err = m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return list, int64(total), nil
}

// GetRole 查询角色及其已分配菜单 ID 列表。
func GetRole(ctx context.Context, id int64) (*model.Role, []int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return nil, nil, errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	var role model.Role
	if err := ra.DB().Model(consts.TableRole).Where("id", id).Where("tenant_id", tenantID).Scan(&role); err != nil {
		return nil, nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if role.Id == 0 {
		return nil, nil, errcode.New(errcode.CodeNotFound, "角色不存在")
	}

	var menuRows []struct {
		MenuId int64 `orm:"menu_id"`
	}
	if err := ra.DB().Model(consts.TableRoleMenu).
		Where("role_id", id).
		Where("tenant_id", tenantID).
		Fields("menu_id").
		Scan(&menuRows); err != nil {
		return nil, nil, errcode.Wrap(errcode.CodeDB, err)
	}
	menuIds := make([]int64, 0, len(menuRows))
	for _, m := range menuRows {
		menuIds = append(menuIds, m.MenuId)
	}
	return &role, menuIds, nil
}

// UpdateRole 更新角色（角色编码为身份标识，不可修改）。
func UpdateRole(ctx context.Context, role *model.Role) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	if _, err := ra.DB().Model(consts.TableRole).
		Where("id", role.Id).
		Where("tenant_id", tenantID).
		Data(g.Map{
			"name":       role.Name,
			"status":     role.Status,
			"remark":     role.Remark,
			"updated_at": gtime.Now(),
		}).Update(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// DeleteRole 删除角色（同时清理其菜单与用户关联）。
func DeleteRole(ctx context.Context, id int64) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	err := ra.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(consts.TableRoleMenu).
			Where("role_id", id).
			Where("tenant_id", tenantID).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(consts.TableUserRole).
			Where("role_id", id).
			Where("tenant_id", tenantID).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(consts.TableRole).
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

// AssignRoleMenus 为角色分配菜单权限（关联表整体替换）。
func AssignRoleMenus(ctx context.Context, roleID int64, menuIds []int64) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	n, err := ra.DB().Model(consts.TableRole).Where("id", roleID).Where("tenant_id", tenantID).Count()
	if err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	if n == 0 {
		return errcode.New(errcode.CodeNotFound, "角色不存在")
	}

	err = ra.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(consts.TableRoleMenu).
			Where("role_id", roleID).
			Where("tenant_id", tenantID).
			Delete(); err != nil {
			return err
		}
		for _, menuID := range menuIds {
			if _, err := tx.Model(consts.TableRoleMenu).Data(g.Map{
				"id":        ra.DB().NextID(),
				"role_id":   roleID,
				"menu_id":   menuID,
				"tenant_id": tenantID,
			}).Insert(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// roleCodeExists 判断角色编码在租户内是否已存在（excludeID 用于更新时排除自身）。
func roleCodeExists(ctx context.Context, tenantID int64, code string, excludeID int64) (bool, error) {
	m := ra.DB().Model(consts.TableRole).Where("tenant_id", tenantID).Where("code", code)
	if excludeID > 0 {
		m = m.WhereNot("id", excludeID)
	}
	n, err := m.Count()
	return n > 0, err
}
