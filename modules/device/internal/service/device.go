package service

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/device/internal/consts"
	"github.com/roamdream/roamadmin/modules/device/internal/model"
)

// CreateDevice 创建设备（ProductKey 取自物模型，ProductKey+DeviceKey 全局唯一）。
func CreateDevice(ctx context.Context, in *model.Device) (int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return 0, errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	mdl, err := GetModel(ctx, in.ModelId)
	if err != nil {
		return 0, err
	}
	in.ProductKey = mdl.Key
	if in.Secret == "" {
		in.Secret = grand.S(32)
	}

	lockKey := consts.CachePrefixDeviceKey + in.ProductKey + ":" + in.DeviceKey
	unlock, err := ra.Lock().Lock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return 0, errcode.Wrap(errcode.CodeInternal, err)
	}
	defer unlock()

	if exists, err := deviceKeyExists(ctx, in.ProductKey, in.DeviceKey, 0); err != nil {
		return 0, err
	} else if exists {
		return 0, errcode.New(errcode.CodeDuplicate, "device key already exists in this product")
	}

	id := ra.DB().NextID()
	if _, err := ra.DB().Model(consts.TableDevice).Data(g.Map{
		"id":          id,
		"tenant_id":   tenantID,
		"model_id":    in.ModelId,
		"product_key": in.ProductKey,
		"device_key":  in.DeviceKey,
		"name":        in.Name,
		"secret":      in.Secret,
		"status":      in.Status,
		"created_at":  gtime.Now(),
	}).Insert(); err != nil {
		return 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return id, nil
}

// ListDevice 分页查询设备（租户隔离）。
func ListDevice(ctx context.Context, page, size int, modelID int64, name string) ([]*model.Device, int64, error) {
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
	m := ra.DB().Model(consts.TableDevice).Where("tenant_id", tenantID)
	if modelID > 0 {
		m = m.Where("model_id", modelID)
	}
	if name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}
	total, err := m.Count()
	if err != nil || total == 0 {
		return []*model.Device{}, int64(total), err
	}
	var list []*model.Device
	if err = m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return list, int64(total), nil
}

// GetDevice 查询设备。
func GetDevice(ctx context.Context, id int64) (*model.Device, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return nil, errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	var list []*model.Device
	if err := ra.DB().Model(consts.TableDevice).Where("id", id).Where("tenant_id", tenantID).Scan(&list); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if len(list) == 0 {
		return nil, errcode.New(errcode.CodeNotFound, "device not found")
	}
	return list[0], nil
}

// GetDeviceByKey 按 ProductKey+DeviceKey 全局查询设备（MQTT 上行解析用，无租户上下文）。
func GetDeviceByKey(ctx context.Context, productKey, deviceKey string) (*model.Device, error) {
	var list []*model.Device
	if err := ra.DB().Model(consts.TableDevice).
		Where("product_key", productKey).
		Where("device_key", deviceKey).
		Scan(&list); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if len(list) == 0 {
		return nil, errcode.New(errcode.CodeNotFound, "device not found")
	}
	return list[0], nil
}

// UpdateDevice 更新设备（DeviceKey / Secret 不可修改）。
func UpdateDevice(ctx context.Context, in *model.Device) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	if _, err := ra.DB().Model(consts.TableDevice).
		Where("id", in.Id).
		Where("tenant_id", tenantID).
		Data(g.Map{"name": in.Name, "status": in.Status}).
		Update(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// DeleteDevice 删除设备（连带删除影子）。
func DeleteDevice(ctx context.Context, id int64) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	var list []*model.Device
	if err := ra.DB().Model(consts.TableDevice).Where("id", id).Where("tenant_id", tenantID).Scan(&list); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	if len(list) == 0 {
		return errcode.New(errcode.CodeNotFound, "device not found")
	}
	if _, err := ra.DB().Model(consts.TableDeviceShadow).Where("device_id", id).Delete(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	if _, err := ra.DB().Model(consts.TableDevice).Where("id", id).Where("tenant_id", tenantID).Delete(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// deviceKeyExists 判断 ProductKey+DeviceKey 是否已存在（全局）。
func deviceKeyExists(ctx context.Context, productKey, deviceKey string, excludeID int64) (bool, error) {
	m := ra.DB().Model(consts.TableDevice).Where("product_key", productKey).Where("device_key", deviceKey)
	if excludeID > 0 {
		m = m.WhereNot("id", excludeID)
	}
	n, err := m.Count()
	return n > 0, err
}
