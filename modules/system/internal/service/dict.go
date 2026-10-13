package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/system/internal/consts"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
)

// CreateDictType 创建字典类型（类型编码在同一租户内唯一，分布式锁防并发重复）。
func CreateDictType(ctx context.Context, in *model.DictType) (int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return 0, errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}

	lockKey := fmt.Sprintf("%s%d:%s", consts.CachePrefixDict, tenantID, in.Type)
	unlock, err := ra.Lock().Lock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return 0, errcode.Wrap(errcode.CodeInternal, err)
	}
	defer unlock()

	if exists, err := dictTypeExists(ctx, tenantID, in.Type, 0); err != nil {
		return 0, err
	} else if exists {
		return 0, errcode.New(errcode.CodeDuplicate, "字典类型编码已存在")
	}

	id := ra.DB().NextID()
	if _, err := ra.DB().Model(consts.TableDictType).Data(g.Map{
		"id":         id,
		"tenant_id":  tenantID,
		"name":       in.Name,
		"type":       in.Type,
		"status":     in.Status,
		"remark":     in.Remark,
		"created_at": gtime.Now(),
	}).Insert(); err != nil {
		return 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return id, nil
}

// ListDictType 分页查询字典类型（按租户隔离）。
func ListDictType(ctx context.Context, page, size int, name string) ([]*model.DictType, int64, error) {
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
	m := ra.DB().Model(consts.TableDictType).Where("tenant_id", tenantID)
	if name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}
	total, err := m.Count()
	if err != nil || total == 0 {
		return []*model.DictType{}, int64(total), err
	}
	var list []*model.DictType
	if err = m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return list, int64(total), nil
}

// GetDictType 查询单个字典类型。
func GetDictType(ctx context.Context, id int64) (*model.DictType, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return nil, errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	var in model.DictType
	if err := ra.DB().Model(consts.TableDictType).Where("id", id).Where("tenant_id", tenantID).Scan(&in); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if in.Id == 0 {
		return nil, errcode.New(errcode.CodeNotFound, "字典类型不存在")
	}
	return &in, nil
}

// UpdateDictType 更新字典类型（类型编码为身份标识，不可修改）。
func UpdateDictType(ctx context.Context, in *model.DictType) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	if _, err := ra.DB().Model(consts.TableDictType).
		Where("id", in.Id).
		Where("tenant_id", tenantID).
		Data(g.Map{"name": in.Name, "status": in.Status, "remark": in.Remark}).
		Update(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// DeleteDictType 删除字典类型（连同其下字典数据一并清理）。
func DeleteDictType(ctx context.Context, id int64) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	var in model.DictType
	if err := ra.DB().Model(consts.TableDictType).Where("id", id).Where("tenant_id", tenantID).Scan(&in); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	if in.Id == 0 {
		return errcode.New(errcode.CodeNotFound, "字典类型不存在")
	}
	if _, err := ra.DB().Model(consts.TableDictData).
		Where("dict_type", in.Type).
		Where("tenant_id", tenantID).
		Delete(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	if _, err := ra.DB().Model(consts.TableDictType).
		Where("id", id).
		Where("tenant_id", tenantID).
		Delete(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// CreateDictData 创建字典数据。
func CreateDictData(ctx context.Context, in *model.DictData) (int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return 0, errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	id := ra.DB().NextID()
	if _, err := ra.DB().Model(consts.TableDictData).Data(g.Map{
		"id":         id,
		"tenant_id":  tenantID,
		"dict_type":  in.DictType,
		"label":      in.Label,
		"value":      in.Value,
		"sort":       in.Sort,
		"status":     in.Status,
		"remark":     in.Remark,
		"created_at": gtime.Now(),
	}).Insert(); err != nil {
		return 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return id, nil
}

// ListDictData 分页查询字典数据（按租户隔离）。
func ListDictData(ctx context.Context, page, size int, dictType, label string) ([]*model.DictData, int64, error) {
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
	m := ra.DB().Model(consts.TableDictData).Where("tenant_id", tenantID)
	if dictType != "" {
		m = m.Where("dict_type", dictType)
	}
	if label != "" {
		m = m.WhereLike("label", "%"+label+"%")
	}
	total, err := m.Count()
	if err != nil || total == 0 {
		return []*model.DictData{}, int64(total), err
	}
	var list []*model.DictData
	if err = m.Page(page, size).OrderAsc("sort").OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return list, int64(total), nil
}

// GetDictData 查询单个字典数据。
func GetDictData(ctx context.Context, id int64) (*model.DictData, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return nil, errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	var in model.DictData
	if err := ra.DB().Model(consts.TableDictData).Where("id", id).Where("tenant_id", tenantID).Scan(&in); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if in.Id == 0 {
		return nil, errcode.New(errcode.CodeNotFound, "字典数据不存在")
	}
	return &in, nil
}

// UpdateDictData 更新字典数据。
func UpdateDictData(ctx context.Context, in *model.DictData) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	if _, err := ra.DB().Model(consts.TableDictData).
		Where("id", in.Id).
		Where("tenant_id", tenantID).
		Data(g.Map{
			"dict_type": in.DictType,
			"label":     in.Label,
			"value":     in.Value,
			"sort":      in.Sort,
			"status":    in.Status,
			"remark":    in.Remark,
		}).Update(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// DeleteDictData 删除字典数据。
func DeleteDictData(ctx context.Context, id int64) error {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	if _, err := ra.DB().Model(consts.TableDictData).
		Where("id", id).
		Where("tenant_id", tenantID).
		Delete(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// dictTypeExists 判断字典类型编码在租户内是否已存在。
func dictTypeExists(ctx context.Context, tenantID int64, typ string, excludeID int64) (bool, error) {
	m := ra.DB().Model(consts.TableDictType).Where("tenant_id", tenantID).Where("type", typ)
	if excludeID > 0 {
		m = m.WhereNot("id", excludeID)
	}
	n, err := m.Count()
	return n > 0, err
}
