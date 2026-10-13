package service

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/system/internal/consts"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
)

// CreateMenu 创建菜单（雪花主键）。
func CreateMenu(ctx context.Context, menu *model.Menu) (int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return 0, errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	id := ra.DB().NextID()
	if _, err := ra.DB().Model(consts.TableMenu).Data(g.Map{
		"id":         id,
		"tenant_id":  tenantID,
		"parent_id":  menu.ParentId,
		"name":       menu.Name,
		"type":       menu.Type,
		"perm":       menu.Perm,
		"path":       menu.Path,
		"component":  menu.Component,
		"icon":       menu.Icon,
		"sort":       menu.Sort,
		"status":     menu.Status,
		"created_at": gtime.Now(),
	}).Insert(); err != nil {
		return 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return id, nil
}

// ListMenu 查询菜单树（按租户隔离，按 sort、id 升序）。
func ListMenu(ctx context.Context, name string) ([]*model.MenuNode, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return nil, errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	m := ra.DB().Model(consts.TableMenu).Where("tenant_id", tenantID)
	if name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}
	var list []*model.Menu
	if err := m.OrderAsc("sort").OrderAsc("id").Scan(&list); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	return buildMenuTree(list, 0), nil
}

// GetMenu 查询单个菜单。
func GetMenu(ctx context.Context, id int64) (*model.Menu, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return nil, errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	var menu model.Menu
	if err := ra.DB().Model(consts.TableMenu).Where("id", id).Where("tenant_id", tenantID).Scan(&menu); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if menu.Id == 0 {
		return nil, errcode.New(errcode.CodeNotFound, "菜单不存在")
	}
	return &menu, nil
}

// UpdateMenu 更新菜单。
func UpdateMenu(ctx context.Context, menu *model.Menu) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	if _, err := ra.DB().Model(consts.TableMenu).
		Where("id", menu.Id).
		Where("tenant_id", tenantID).
		Data(g.Map{
			"parent_id": menu.ParentId,
			"name":      menu.Name,
			"type":      menu.Type,
			"perm":      menu.Perm,
			"path":      menu.Path,
			"component": menu.Component,
			"icon":      menu.Icon,
			"sort":      menu.Sort,
			"status":    menu.Status,
		}).Update(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// DeleteMenu 删除菜单（连同其后代节点与角色关联一并清理）。
func DeleteMenu(ctx context.Context, id int64) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	err := ra.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		ids, err := collectMenuDescendants(ctx, tx, tenantID, id)
		if err != nil {
			return err
		}
		ids = append(ids, id)
		if _, err := tx.Model(consts.TableRoleMenu).
			Where("menu_id", ids).
			Where("tenant_id", tenantID).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(consts.TableMenu).
			Where("id", ids).
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

// buildMenuTree 将扁平菜单列表构建为树。
func buildMenuTree(list []*model.Menu, parentID int64) []*model.MenuNode {
	nodes := make([]*model.MenuNode, 0)
	for _, m := range list {
		if m.ParentId == parentID {
			nodes = append(nodes, &model.MenuNode{
				Menu:     m,
				Children: buildMenuTree(list, m.Id),
			})
		}
	}
	return nodes
}

// collectMenuDescendants 收集某菜单的全部后代 ID（含多级）。
func collectMenuDescendants(ctx context.Context, tx gdb.TX, tenantID, id int64) ([]int64, error) {
	var list []*model.Menu
	if err := tx.Model(consts.TableMenu).Where("tenant_id", tenantID).Scan(&list); err != nil {
		return nil, err
	}
	ids := make([]int64, 0)
	var walk func(pid int64)
	walk = func(pid int64) {
		for _, m := range list {
			if m.ParentId == pid {
				ids = append(ids, m.Id)
				walk(m.Id)
			}
		}
	}
	walk(id)
	return ids, nil
}
