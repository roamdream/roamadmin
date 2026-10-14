package service

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/device/internal/consts"
	"github.com/roamdream/roamadmin/modules/device/internal/model"
)

// CreateModel 创建物模型（ProductKey 全局唯一，分布式锁防并发重复）。
func CreateModel(ctx context.Context, in *model.DeviceModel) (int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return 0, errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}

	lockKey := consts.CachePrefixModelKey + in.Key
	unlock, err := ra.Lock().Lock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return 0, errcode.Wrap(errcode.CodeInternal, err)
	}
	defer unlock()

	if exists, err := modelKeyExists(ctx, in.Key, 0); err != nil {
		return 0, err
	} else if exists {
		return 0, errcode.New(errcode.CodeDuplicate, "product key already exists")
	}

	id := ra.DB().NextID()
	if _, err := ra.DB().Model(consts.TableDeviceModel).Data(g.Map{
		"id":          id,
		"tenant_id":   tenantID,
		"name":        in.Name,
		"product_key": in.Key,
		"description": in.Description,
		"status":      in.Status,
		"created_at":  gtime.Now(),
	}).Insert(); err != nil {
		return 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return id, nil
}

// ListModel 分页查询物模型（租户隔离）。
func ListModel(ctx context.Context, page, size int, name string) ([]*model.DeviceModel, int64, error) {
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
	m := ra.DB().Model(consts.TableDeviceModel).Where("tenant_id", tenantID)
	if name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}
	total, err := m.Count()
	if err != nil || total == 0 {
		return []*model.DeviceModel{}, int64(total), err
	}
	var list []*model.DeviceModel
	if err = m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return list, int64(total), nil
}

// GetModel 查询物模型。
func GetModel(ctx context.Context, id int64) (*model.DeviceModel, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return nil, errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	var list []*model.DeviceModel
	if err := ra.DB().Model(consts.TableDeviceModel).Where("id", id).Where("tenant_id", tenantID).Scan(&list); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if len(list) == 0 {
		return nil, errcode.New(errcode.CodeNotFound, "device model not found")
	}
	return list[0], nil
}

// UpdateModel 更新物模型（Key 为身份标识，不可修改）。
func UpdateModel(ctx context.Context, in *model.DeviceModel) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	if _, err := ra.DB().Model(consts.TableDeviceModel).
		Where("id", in.Id).
		Where("tenant_id", tenantID).
		Data(g.Map{"name": in.Name, "description": in.Description, "status": in.Status}).
		Update(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// DeleteModel 删除物模型（存在设备时禁止删除；连带删除属性）。
func DeleteModel(ctx context.Context, id int64) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	var list []*model.DeviceModel
	if err := ra.DB().Model(consts.TableDeviceModel).Where("id", id).Where("tenant_id", tenantID).Scan(&list); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	if len(list) == 0 {
		return errcode.New(errcode.CodeNotFound, "device model not found")
	}
	if n, err := ra.DB().Model(consts.TableDevice).Where("model_id", id).Count(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	} else if n > 0 {
		return errcode.New(errcode.CodeConflict, "cannot delete model with existing devices")
	}
	if _, err := ra.DB().Model(consts.TableDeviceModelProperty).Where("model_id", id).Delete(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	if _, err := ra.DB().Model(consts.TableDeviceModel).Where("id", id).Where("tenant_id", tenantID).Delete(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// GetModelByKey 按 ProductKey 全局查询物模型（MQTT 上行解析用，无租户上下文）。
func GetModelByKey(ctx context.Context, key string) (*model.DeviceModel, error) {
	var list []*model.DeviceModel
	if err := ra.DB().Model(consts.TableDeviceModel).Where("product_key", key).Scan(&list); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if len(list) == 0 {
		return nil, errcode.New(errcode.CodeNotFound, "device model not found")
	}
	return list[0], nil
}

// CreateProperty 创建属性（Identifier 在同一物模型内唯一）。
func CreateProperty(ctx context.Context, in *model.DeviceModelProperty) (int64, error) {
	if _, ok := auth.TenantFrom(ctx); !ok {
		return 0, errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	if _, err := GetModel(ctx, in.ModelId); err != nil {
		return 0, err
	}
	if exists, err := propertyExists(ctx, in.ModelId, in.Identifier, 0); err != nil {
		return 0, err
	} else if exists {
		return 0, errcode.New(errcode.CodeDuplicate, "property identifier already exists in this model")
	}
	id := ra.DB().NextID()
	if _, err := ra.DB().Model(consts.TableDeviceModelProperty).Data(g.Map{
		"id":         id,
		"model_id":   in.ModelId,
		"identifier": in.Identifier,
		"name":       in.Name,
		"data_type":  in.DataType,
		"unit":       in.Unit,
		"access":     in.Access,
		"sort":       in.Sort,
	}).Insert(); err != nil {
		return 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return id, nil
}

// ListProperty 查询物模型属性列表。
func ListProperty(ctx context.Context, modelID int64) ([]*model.DeviceModelProperty, error) {
	if _, err := GetModel(ctx, modelID); err != nil {
		return nil, err
	}
	var list []*model.DeviceModelProperty
	if err := ra.DB().Model(consts.TableDeviceModelProperty).
		Where("model_id", modelID).
		OrderAsc("sort").OrderAsc("id").
		Scan(&list); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	return list, nil
}

// UpdateProperty 更新属性（Identifier 为身份标识，不可修改）。
func UpdateProperty(ctx context.Context, in *model.DeviceModelProperty) error {
	if _, err := GetModel(ctx, in.ModelId); err != nil {
		return err
	}
	if _, err := ra.DB().Model(consts.TableDeviceModelProperty).
		Where("id", in.Id).
		Where("model_id", in.ModelId).
		Data(g.Map{
			"name":      in.Name,
			"data_type": in.DataType,
			"unit":      in.Unit,
			"access":    in.Access,
			"sort":      in.Sort,
		}).Update(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// DeleteProperty 删除属性。
func DeleteProperty(ctx context.Context, modelID, id int64) error {
	if _, err := GetModel(ctx, modelID); err != nil {
		return err
	}
	if _, err := ra.DB().Model(consts.TableDeviceModelProperty).
		Where("id", id).
		Where("model_id", modelID).
		Delete(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// modelKeyExists 判断 ProductKey 是否已存在（全局，不区分租户）。
func modelKeyExists(ctx context.Context, key string, excludeID int64) (bool, error) {
	m := ra.DB().Model(consts.TableDeviceModel).Where("product_key", key)
	if excludeID > 0 {
		m = m.WhereNot("id", excludeID)
	}
	n, err := m.Count()
	return n > 0, err
}

// propertyExists 判断属性标识在同一物模型内是否已存在。
func propertyExists(ctx context.Context, modelID int64, identifier string, excludeID int64) (bool, error) {
	m := ra.DB().Model(consts.TableDeviceModelProperty).Where("model_id", modelID).Where("identifier", identifier)
	if excludeID > 0 {
		m = m.WhereNot("id", excludeID)
	}
	n, err := m.Count()
	return n > 0, err
}
