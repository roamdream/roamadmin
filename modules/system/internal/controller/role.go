package controller

import (
	"context"

	"github.com/roamdream/roamadmin/modules/system/api"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
	"github.com/roamdream/roamadmin/modules/system/internal/service"
)

// Role 角色控制器（角色 CRUD + 分配权限）。
type Role struct{}

// NewRole 创建角色控制器。
func NewRole() *Role { return &Role{} }

// Create 创建角色。
func (c *Role) Create(ctx context.Context, req *api.RoleCreateReq) (*api.RoleCreateRes, error) {
	id, err := service.CreateRole(ctx, &model.Role{
		Name:   req.Name,
		Code:   req.Code,
		Status: req.Status,
		Remark: req.Remark,
	})
	writeOper(ctx, "角色", "创建", err)
	if err != nil {
		return nil, err
	}
	return &api.RoleCreateRes{Id: id}, nil
}

// List 分页查询角色。
func (c *Role) List(ctx context.Context, req *api.RoleListReq) (*api.RoleListRes, error) {
	list, total, err := service.ListRole(ctx, req.Page, req.Size, req.Name)
	if err != nil {
		return nil, err
	}
	out := make([]*api.RoleItem, 0, len(list))
	for _, it := range list {
		out = append(out, &api.RoleItem{
			Id:        it.Id,
			TenantId:  it.TenantId,
			Name:      it.Name,
			Code:      it.Code,
			Status:    it.Status,
			Remark:    it.Remark,
			CreatedAt: it.CreatedAt.String(),
		})
	}
	return &api.RoleListRes{List: out, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Get 查询角色。
func (c *Role) Get(ctx context.Context, req *api.RoleGetReq) (*api.RoleGetRes, error) {
	role, menuIds, err := service.GetRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.RoleGetRes{
		Role: &api.RoleItem{
			Id:        role.Id,
			TenantId:  role.TenantId,
			Name:      role.Name,
			Code:      role.Code,
			Status:    role.Status,
			Remark:    role.Remark,
			CreatedAt: role.CreatedAt.String(),
		},
		MenuIds: menuIds,
	}, nil
}

// Update 更新角色。
func (c *Role) Update(ctx context.Context, req *api.RoleUpdateReq) (*api.RoleUpdateRes, error) {
	err := service.UpdateRole(ctx, &model.Role{
		Id:     req.Id,
		Name:   req.Name,
		Status: req.Status,
		Remark: req.Remark,
	})
	writeOper(ctx, "角色", "更新", err)
	if err != nil {
		return nil, err
	}
	return &api.RoleUpdateRes{Updated: true}, nil
}

// Delete 删除角色。
func (c *Role) Delete(ctx context.Context, req *api.RoleDeleteReq) (*api.RoleDeleteRes, error) {
	err := service.DeleteRole(ctx, req.Id)
	writeOper(ctx, "角色", "删除", err)
	if err != nil {
		return nil, err
	}
	return &api.RoleDeleteRes{Deleted: true}, nil
}

// AssignMenus 为角色分配菜单权限。
func (c *Role) AssignMenus(ctx context.Context, req *api.RoleAssignMenuReq) (*api.RoleAssignMenuRes, error) {
	err := service.AssignRoleMenus(ctx, req.Id, req.MenuIds)
	writeOper(ctx, "角色", "分配权限", err)
	if err != nil {
		return nil, err
	}
	return &api.RoleAssignMenuRes{Assigned: true}, nil
}
