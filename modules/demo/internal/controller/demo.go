// Package controller demo 模块 HTTP 控制器。
package controller

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/modules/demo/api"
	"github.com/roamdream/roamadmin/modules/demo/internal/service"
)

// Controller demo 控制器。
type Controller struct{}

// New 创建控制器。
func New() *Controller { return &Controller{} }

// Ping 连通性测试（不依赖 DB，验证框架与模块装配）。
func (c *Controller) Ping(ctx context.Context, req *api.PingReq) (*api.PingRes, error) {
	return &api.PingRes{
		Module:  "demo",
		Version: "v0.1.0",
		Time:    gtime.Now().String(),
	}, nil
}

// Create 创建示例数据。
func (c *Controller) Create(ctx context.Context, req *api.CreateReq) (*api.CreateRes, error) {
	id, err := service.Create(ctx, req.Name, req.Remark)
	if err != nil {
		return nil, err
	}
	return &api.CreateRes{Id: id}, nil
}

// List 分页查询示例数据。
func (c *Controller) List(ctx context.Context, req *api.ListReq) (*api.ListRes, error) {
	items, total, err := service.List(ctx, req.Page, req.Size, req.Name)
	if err != nil {
		return nil, err
	}
	out := make([]*api.Item, 0, len(items))
	for _, it := range items {
		out = append(out, &api.Item{
			Id:        it.Id,
			Name:      it.Name,
			Remark:    it.Remark,
			CreatedAt: it.CreatedAt.String(),
		})
	}
	return &api.ListRes{List: out, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Get 查询单条示例数据。
func (c *Controller) Get(ctx context.Context, req *api.GetReq) (*api.Item, error) {
	it, err := service.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.Item{
		Id:        it.Id,
		Name:      it.Name,
		Remark:    it.Remark,
		CreatedAt: it.CreatedAt.String(),
	}, nil
}

// Delete 删除示例数据。
func (c *Controller) Delete(ctx context.Context, req *api.DeleteReq) (*api.DeleteRes, error) {
	if err := service.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &api.DeleteRes{Deleted: true}, nil
}
