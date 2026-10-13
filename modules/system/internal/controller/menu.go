package controller

import (
	"context"

	"github.com/roamdream/roamadmin/modules/system/api"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
	"github.com/roamdream/roamadmin/modules/system/internal/service"
)

// Menu 菜单控制器（菜单 CRUD，树形）。
type Menu struct{}

// NewMenu 创建菜单控制器。
func NewMenu() *Menu { return &Menu{} }

// Create 创建菜单。
func (c *Menu) Create(ctx context.Context, req *api.MenuCreateReq) (*api.MenuCreateRes, error) {
	id, err := service.CreateMenu(ctx, &model.Menu{
		ParentId:  req.ParentId,
		Name:      req.Name,
		Type:      req.Type,
		Perm:      req.Perm,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Status:    req.Status,
	})
	writeOper(ctx, "菜单", "创建", err)
	if err != nil {
		return nil, err
	}
	return &api.MenuCreateRes{Id: id}, nil
}

// List 查询菜单树。
func (c *Menu) List(ctx context.Context, req *api.MenuListReq) (*api.MenuListRes, error) {
	nodes, err := service.ListMenu(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &api.MenuListRes{List: toMenuTreeItems(nodes)}, nil
}

// Get 查询菜单。
func (c *Menu) Get(ctx context.Context, req *api.MenuGetReq) (*api.MenuGetRes, error) {
	menu, err := service.GetMenu(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.MenuGetRes{Menu: toMenuItem(menu)}, nil
}

// Update 更新菜单。
func (c *Menu) Update(ctx context.Context, req *api.MenuUpdateReq) (*api.MenuUpdateRes, error) {
	err := service.UpdateMenu(ctx, &model.Menu{
		Id:        req.Id,
		ParentId:  req.ParentId,
		Name:      req.Name,
		Type:      req.Type,
		Perm:      req.Perm,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Status:    req.Status,
	})
	writeOper(ctx, "菜单", "更新", err)
	if err != nil {
		return nil, err
	}
	return &api.MenuUpdateRes{Updated: true}, nil
}

// Delete 删除菜单。
func (c *Menu) Delete(ctx context.Context, req *api.MenuDeleteReq) (*api.MenuDeleteRes, error) {
	err := service.DeleteMenu(ctx, req.Id)
	writeOper(ctx, "菜单", "删除", err)
	if err != nil {
		return nil, err
	}
	return &api.MenuDeleteRes{Deleted: true}, nil
}

// toMenuItem 模型菜单转 DTO。
func toMenuItem(m *model.Menu) *api.MenuItem {
	return &api.MenuItem{
		Id:        m.Id,
		TenantId:  m.TenantId,
		ParentId:  m.ParentId,
		Name:      m.Name,
		Type:      m.Type,
		Perm:      m.Perm,
		Path:      m.Path,
		Component: m.Component,
		Icon:      m.Icon,
		Sort:      m.Sort,
		Status:    m.Status,
		CreatedAt: m.CreatedAt.String(),
	}
}

// toMenuTreeItems 模型菜单树转 DTO。
func toMenuTreeItems(nodes []*model.MenuNode) []*api.MenuTreeItem {
	out := make([]*api.MenuTreeItem, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, &api.MenuTreeItem{
			MenuItem: toMenuItem(n.Menu),
			Children: toMenuTreeItems(n.Children),
		})
	}
	return out
}
