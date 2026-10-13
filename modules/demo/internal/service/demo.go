// Package service demo 模块业务逻辑。
package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/demo/internal/model"
)

const table = "ra_demo_item"

// Create 创建示例数据（雪花主键）。
func Create(ctx context.Context, name, remark string) (int64, error) {
	id := ra.DB().NextID()
	_, err := ra.DB().Model(table).Data(g.Map{
		"id":         id,
		"name":       name,
		"remark":     remark,
		"created_at": gtime.Now(),
	}).Insert()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// List 分页查询。
func List(ctx context.Context, page, size int, name string) ([]*model.Item, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 20
	}
	m := ra.DB().Model(table)
	if name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}
	total, err := m.Count()
	if err != nil || total == 0 {
		return []*model.Item{}, int64(total), err
	}
	var list []*model.Item
	if err = m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, err
	}
	return list, int64(total), nil
}

// Get 查询单条。
func Get(ctx context.Context, id int64) (*model.Item, error) {
	var item model.Item
	if err := ra.DB().Model(table).Where("id", id).Scan(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

// Delete 删除单条。
func Delete(ctx context.Context, id int64) error {
	_, err := ra.DB().Model(table).Where("id", id).Delete()
	return err
}
