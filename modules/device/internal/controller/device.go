package controller

import (
	"context"

	"github.com/roamdream/roamadmin/kernel/utils"
	"github.com/roamdream/roamadmin/modules/device/api"
	"github.com/roamdream/roamadmin/modules/device/internal/model"
	"github.com/roamdream/roamadmin/modules/device/internal/service"
)

// Device 设备控制器。
type Device struct{}

// NewDevice 创建设备控制器。
func NewDevice() *Device { return &Device{} }

// Create 创建设备。
func (c *Device) Create(ctx context.Context, req *api.DeviceCreateReq) (*api.DeviceCreateRes, error) {
	id, err := service.CreateDevice(ctx, &model.Device{
		ModelId:   req.ModelId,
		DeviceKey: req.DeviceKey,
		Name:      req.Name,
		Secret:    req.Secret,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &api.DeviceCreateRes{Id: id}, nil
}

// List 分页查询设备。
func (c *Device) List(ctx context.Context, req *api.DeviceListReq) (*api.DeviceListRes, error) {
	list, total, err := service.ListDevice(ctx, req.Page, req.Size, req.ModelId, req.Name)
	if err != nil {
		return nil, err
	}
	out := make([]*api.DeviceItem, 0, len(list))
	for _, it := range list {
		out = append(out, toDeviceItem(it))
	}
	return &api.DeviceListRes{List: out, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Get 查询设备。
func (c *Device) Get(ctx context.Context, req *api.DeviceGetReq) (*api.DeviceGetRes, error) {
	it, err := service.GetDevice(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DeviceGetRes{Device: toDeviceItem(it)}, nil
}

// Update 更新设备。
func (c *Device) Update(ctx context.Context, req *api.DeviceUpdateReq) (*api.DeviceUpdateRes, error) {
	err := service.UpdateDevice(ctx, &model.Device{Id: req.Id, Name: req.Name, Status: req.Status})
	if err != nil {
		return nil, err
	}
	return &api.DeviceUpdateRes{Updated: true}, nil
}

// Delete 删除设备。
func (c *Device) Delete(ctx context.Context, req *api.DeviceDeleteReq) (*api.DeviceDeleteRes, error) {
	if err := service.DeleteDevice(ctx, req.Id); err != nil {
		return nil, err
	}
	return &api.DeviceDeleteRes{Deleted: true}, nil
}

// SendCommand 下发命令。
func (c *Device) SendCommand(ctx context.Context, req *api.DeviceDownlinkReq) (*api.DeviceDownlinkRes, error) {
	topic, payload, err := service.PublishCommand(ctx, req.Id, req.Data)
	if err != nil {
		return nil, err
	}
	return &api.DeviceDownlinkRes{Topic: topic, Payload: payload}, nil
}

func toDeviceItem(it *model.Device) *api.DeviceItem {
	return &api.DeviceItem{
		Id:         it.Id,
		TenantId:   it.TenantId,
		ModelId:    it.ModelId,
		ProductKey: it.ProductKey,
		DeviceKey:  it.DeviceKey,
		Name:       it.Name,
		Secret:     it.Secret,
		Status:     it.Status,
		CreatedAt:  utils.TimeRFC3339(it.CreatedAt),
	}
}
