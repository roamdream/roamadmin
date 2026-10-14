package controller

import (
	"context"

	"github.com/roamdream/roamadmin/kernel/utils"
	"github.com/roamdream/roamadmin/modules/device/api"
	"github.com/roamdream/roamadmin/modules/device/internal/model"
	"github.com/roamdream/roamadmin/modules/device/internal/service"
)

// Model 物模型控制器（物模型 + 属性）。
type Model struct{}

// NewModel 创建物模型控制器。
func NewModel() *Model { return &Model{} }

// Create 创建物模型。
func (c *Model) Create(ctx context.Context, req *api.DeviceModelCreateReq) (*api.DeviceModelCreateRes, error) {
	id, err := service.CreateModel(ctx, &model.DeviceModel{
		Name:        req.Name,
		Key:         req.Key,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &api.DeviceModelCreateRes{Id: id}, nil
}

// List 分页查询物模型。
func (c *Model) List(ctx context.Context, req *api.DeviceModelListReq) (*api.DeviceModelListRes, error) {
	list, total, err := service.ListModel(ctx, req.Page, req.Size, req.Name)
	if err != nil {
		return nil, err
	}
	out := make([]*api.DeviceModelItem, 0, len(list))
	for _, it := range list {
		out = append(out, toModelItem(it))
	}
	return &api.DeviceModelListRes{List: out, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Get 查询物模型。
func (c *Model) Get(ctx context.Context, req *api.DeviceModelGetReq) (*api.DeviceModelGetRes, error) {
	it, err := service.GetModel(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DeviceModelGetRes{Model: toModelItem(it)}, nil
}

// Update 更新物模型。
func (c *Model) Update(ctx context.Context, req *api.DeviceModelUpdateReq) (*api.DeviceModelUpdateRes, error) {
	err := service.UpdateModel(ctx, &model.DeviceModel{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &api.DeviceModelUpdateRes{Updated: true}, nil
}

// Delete 删除物模型。
func (c *Model) Delete(ctx context.Context, req *api.DeviceModelDeleteReq) (*api.DeviceModelDeleteRes, error) {
	if err := service.DeleteModel(ctx, req.Id); err != nil {
		return nil, err
	}
	return &api.DeviceModelDeleteRes{Deleted: true}, nil
}

// CreateProperty 创建属性。
func (c *Model) CreateProperty(ctx context.Context, req *api.PropertyCreateReq) (*api.PropertyCreateRes, error) {
	id, err := service.CreateProperty(ctx, &model.DeviceModelProperty{
		ModelId:    req.ModelId,
		Identifier: req.Identifier,
		Name:       req.Name,
		DataType:   req.DataType,
		Unit:       req.Unit,
		Access:     req.Access,
		Sort:       req.Sort,
	})
	if err != nil {
		return nil, err
	}
	return &api.PropertyCreateRes{Id: id}, nil
}

// ListProperty 查询属性列表。
func (c *Model) ListProperty(ctx context.Context, req *api.PropertyListReq) (*api.PropertyListRes, error) {
	list, err := service.ListProperty(ctx, req.ModelId)
	if err != nil {
		return nil, err
	}
	out := make([]*api.PropertyItem, 0, len(list))
	for _, it := range list {
		out = append(out, toPropertyItem(it))
	}
	return &api.PropertyListRes{List: out}, nil
}

// UpdateProperty 更新属性。
func (c *Model) UpdateProperty(ctx context.Context, req *api.PropertyUpdateReq) (*api.PropertyUpdateRes, error) {
	err := service.UpdateProperty(ctx, &model.DeviceModelProperty{
		Id:       req.Id,
		ModelId:  req.ModelId,
		Name:     req.Name,
		DataType: req.DataType,
		Unit:     req.Unit,
		Access:   req.Access,
		Sort:     req.Sort,
	})
	if err != nil {
		return nil, err
	}
	return &api.PropertyUpdateRes{Updated: true}, nil
}

// DeleteProperty 删除属性。
func (c *Model) DeleteProperty(ctx context.Context, req *api.PropertyDeleteReq) (*api.PropertyDeleteRes, error) {
	if err := service.DeleteProperty(ctx, req.ModelId, req.Id); err != nil {
		return nil, err
	}
	return &api.PropertyDeleteRes{Deleted: true}, nil
}

func toModelItem(it *model.DeviceModel) *api.DeviceModelItem {
	return &api.DeviceModelItem{
		Id:          it.Id,
		TenantId:    it.TenantId,
		Name:        it.Name,
		Key:         it.Key,
		Description: it.Description,
		Status:      it.Status,
		CreatedAt:   utils.TimeRFC3339(it.CreatedAt),
	}
}

func toPropertyItem(it *model.DeviceModelProperty) *api.PropertyItem {
	return &api.PropertyItem{
		Id:         it.Id,
		ModelId:    it.ModelId,
		Identifier: it.Identifier,
		Name:       it.Name,
		DataType:   it.DataType,
		Unit:       it.Unit,
		Access:     it.Access,
		Sort:       it.Sort,
	}
}
