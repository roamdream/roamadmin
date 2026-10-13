package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/system/internal/consts"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
)

// CreateTenant 创建租户（租户编码全局唯一，分布式锁防并发重复）。
func CreateTenant(ctx context.Context, in *model.Tenant) (int64, error) {
	lockKey := fmt.Sprintf("%s%s", consts.CachePrefixTenant, in.Code)
	unlock, err := ra.Lock().Lock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return 0, errcode.Wrap(errcode.CodeInternal, err)
	}
	defer unlock()

	if exists, err := tenantCodeExists(ctx, in.Code, 0); err != nil {
		return 0, err
	} else if exists {
		return 0, errcode.New(errcode.CodeDuplicate, "租户编码已存在")
	}

	id := ra.DB().NextID()
	now := gtime.Now()
	if _, err := ra.DB().Model(consts.TableTenant).Data(g.Map{
		"id":         id,
		"name":       in.Name,
		"code":       in.Code,
		"status":     in.Status,
		"remark":     in.Remark,
		"created_at": now,
		"updated_at": now,
	}).Insert(); err != nil {
		return 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return id, nil
}

// ListTenant 分页查询租户（租户注册表全局可见，不按租户过滤）。
func ListTenant(ctx context.Context, page, size int, name string) ([]*model.Tenant, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 20
	}
	m := ra.DB().Model(consts.TableTenant)
	if name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}
	total, err := m.Count()
	if err != nil || total == 0 {
		return []*model.Tenant{}, int64(total), err
	}
	var list []*model.Tenant
	if err = m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return list, int64(total), nil
}

// GetTenant 查询单个租户。
func GetTenant(ctx context.Context, id int64) (*model.Tenant, error) {
	var in model.Tenant
	if err := ra.DB().Model(consts.TableTenant).Where("id", id).Scan(&in); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	if in.Id == 0 {
		return nil, errcode.New(errcode.CodeNotFound, "租户不存在")
	}
	return &in, nil
}

// UpdateTenant 更新租户（租户编码为身份标识，不可修改）。
func UpdateTenant(ctx context.Context, in *model.Tenant) error {
	if _, err := ra.DB().Model(consts.TableTenant).
		Where("id", in.Id).
		Data(g.Map{
			"name":       in.Name,
			"status":     in.Status,
			"remark":     in.Remark,
			"updated_at": gtime.Now(),
		}).Update(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// DeleteTenant 删除租户。
func DeleteTenant(ctx context.Context, id int64) error {
	if _, err := ra.DB().Model(consts.TableTenant).Where("id", id).Delete(); err != nil {
		return errcode.Wrap(errcode.CodeDB, err)
	}
	return nil
}

// tenantCodeExists 判断租户编码是否已存在。
func tenantCodeExists(ctx context.Context, code string, excludeID int64) (bool, error) {
	m := ra.DB().Model(consts.TableTenant).Where("code", code)
	if excludeID > 0 {
		m = m.WhereNot("id", excludeID)
	}
	n, err := m.Count()
	return n > 0, err
}
