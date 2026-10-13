package controller

import (
	"context"

	"github.com/roamdream/roamadmin/kernel/utils"
	"github.com/roamdream/roamadmin/modules/system/api"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
	"github.com/roamdream/roamadmin/modules/system/internal/service"
)

// Tenant 租户控制器（租户 CRUD，多租户）。
type Tenant struct{}

// NewTenant 创建租户控制器。
func NewTenant() *Tenant { return &Tenant{} }

// Create 创建租户。
func (c *Tenant) Create(ctx context.Context, req *api.TenantCreateReq) (*api.TenantCreateRes, error) {
	id, err := service.CreateTenant(ctx, &model.Tenant{
		Name:   req.Name,
		Code:   req.Code,
		Status: req.Status,
		Remark: req.Remark,
	})
	writeOper(ctx, "tenant", "create", err)
	if err != nil {
		return nil, err
	}
	return &api.TenantCreateRes{Id: id}, nil
}

// List 分页查询租户。
func (c *Tenant) List(ctx context.Context, req *api.TenantListReq) (*api.TenantListRes, error) {
	list, total, err := service.ListTenant(ctx, req.Page, req.Size, req.Name)
	if err != nil {
		return nil, err
	}
	out := make([]*api.TenantItem, 0, len(list))
	for _, it := range list {
		out = append(out, &api.TenantItem{
			Id:        it.Id,
			Name:      it.Name,
			Code:      it.Code,
			Status:    it.Status,
			Remark:    it.Remark,
			CreatedAt: utils.TimeRFC3339(it.CreatedAt),
		})
	}
	return &api.TenantListRes{List: out, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Get 查询租户。
func (c *Tenant) Get(ctx context.Context, req *api.TenantGetReq) (*api.TenantGetRes, error) {
	it, err := service.GetTenant(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.TenantGetRes{Tenant: &api.TenantItem{
		Id:        it.Id,
		Name:      it.Name,
		Code:      it.Code,
		Status:    it.Status,
		Remark:    it.Remark,
		CreatedAt: utils.TimeRFC3339(it.CreatedAt),
	}}, nil
}

// Update 更新租户。
func (c *Tenant) Update(ctx context.Context, req *api.TenantUpdateReq) (*api.TenantUpdateRes, error) {
	err := service.UpdateTenant(ctx, &model.Tenant{
		Id:     req.Id,
		Name:   req.Name,
		Status: req.Status,
		Remark: req.Remark,
	})
	writeOper(ctx, "tenant", "update", err)
	if err != nil {
		return nil, err
	}
	return &api.TenantUpdateRes{Updated: true}, nil
}

// Delete 删除租户。
func (c *Tenant) Delete(ctx context.Context, req *api.TenantDeleteReq) (*api.TenantDeleteRes, error) {
	err := service.DeleteTenant(ctx, req.Id)
	writeOper(ctx, "tenant", "delete", err)
	if err != nil {
		return nil, err
	}
	return &api.TenantDeleteRes{Deleted: true}, nil
}
