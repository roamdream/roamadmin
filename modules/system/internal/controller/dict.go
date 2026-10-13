package controller

import (
	"context"

	"github.com/roamdream/roamadmin/modules/system/api"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
	"github.com/roamdream/roamadmin/modules/system/internal/service"
)

// Dict 字典控制器（字典类型 + 字典数据）。
type Dict struct{}

// NewDict 创建字典控制器。
func NewDict() *Dict { return &Dict{} }

// —— 字典类型 ——

// CreateType 创建字典类型。
func (c *Dict) CreateType(ctx context.Context, req *api.DictTypeCreateReq) (*api.DictTypeCreateRes, error) {
	id, err := service.CreateDictType(ctx, &model.DictType{
		Name:   req.Name,
		Type:   req.Type,
		Status: req.Status,
		Remark: req.Remark,
	})
	writeOper(ctx, "字典", "创建类型", err)
	if err != nil {
		return nil, err
	}
	return &api.DictTypeCreateRes{Id: id}, nil
}

// ListType 分页查询字典类型。
func (c *Dict) ListType(ctx context.Context, req *api.DictTypeListReq) (*api.DictTypeListRes, error) {
	list, total, err := service.ListDictType(ctx, req.Page, req.Size, req.Name)
	if err != nil {
		return nil, err
	}
	out := make([]*api.DictTypeItem, 0, len(list))
	for _, it := range list {
		out = append(out, &api.DictTypeItem{
			Id:        it.Id,
			TenantId:  it.TenantId,
			Name:      it.Name,
			Type:      it.Type,
			Status:    it.Status,
			Remark:    it.Remark,
			CreatedAt: it.CreatedAt.String(),
		})
	}
	return &api.DictTypeListRes{List: out, Total: total, Page: req.Page, Size: req.Size}, nil
}

// GetType 查询字典类型。
func (c *Dict) GetType(ctx context.Context, req *api.DictTypeGetReq) (*api.DictTypeGetRes, error) {
	it, err := service.GetDictType(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DictTypeGetRes{DictType: &api.DictTypeItem{
		Id:        it.Id,
		TenantId:  it.TenantId,
		Name:      it.Name,
		Type:      it.Type,
		Status:    it.Status,
		Remark:    it.Remark,
		CreatedAt: it.CreatedAt.String(),
	}}, nil
}

// UpdateType 更新字典类型。
func (c *Dict) UpdateType(ctx context.Context, req *api.DictTypeUpdateReq) (*api.DictTypeUpdateRes, error) {
	err := service.UpdateDictType(ctx, &model.DictType{
		Id:     req.Id,
		Name:   req.Name,
		Status: req.Status,
		Remark: req.Remark,
	})
	writeOper(ctx, "字典", "更新类型", err)
	if err != nil {
		return nil, err
	}
	return &api.DictTypeUpdateRes{Updated: true}, nil
}

// DeleteType 删除字典类型。
func (c *Dict) DeleteType(ctx context.Context, req *api.DictTypeDeleteReq) (*api.DictTypeDeleteRes, error) {
	err := service.DeleteDictType(ctx, req.Id)
	writeOper(ctx, "字典", "删除类型", err)
	if err != nil {
		return nil, err
	}
	return &api.DictTypeDeleteRes{Deleted: true}, nil
}

// —— 字典数据 ——

// CreateData 创建字典数据。
func (c *Dict) CreateData(ctx context.Context, req *api.DictDataCreateReq) (*api.DictDataCreateRes, error) {
	id, err := service.CreateDictData(ctx, &model.DictData{
		DictType: req.DictType,
		Label:    req.Label,
		Value:    req.Value,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
	})
	writeOper(ctx, "字典", "创建数据", err)
	if err != nil {
		return nil, err
	}
	return &api.DictDataCreateRes{Id: id}, nil
}

// ListData 分页查询字典数据。
func (c *Dict) ListData(ctx context.Context, req *api.DictDataListReq) (*api.DictDataListRes, error) {
	list, total, err := service.ListDictData(ctx, req.Page, req.Size, req.DictType, req.Label)
	if err != nil {
		return nil, err
	}
	out := make([]*api.DictDataItem, 0, len(list))
	for _, it := range list {
		out = append(out, &api.DictDataItem{
			Id:       it.Id,
			TenantId: it.TenantId,
			DictType: it.DictType,
			Label:    it.Label,
			Value:    it.Value,
			Sort:     it.Sort,
			Status:   it.Status,
			Remark:   it.Remark,
		})
	}
	return &api.DictDataListRes{List: out, Total: total, Page: req.Page, Size: req.Size}, nil
}

// GetData 查询字典数据。
func (c *Dict) GetData(ctx context.Context, req *api.DictDataGetReq) (*api.DictDataGetRes, error) {
	it, err := service.GetDictData(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DictDataGetRes{DictData: &api.DictDataItem{
		Id:       it.Id,
		TenantId: it.TenantId,
		DictType: it.DictType,
		Label:    it.Label,
		Value:    it.Value,
		Sort:     it.Sort,
		Status:   it.Status,
		Remark:   it.Remark,
	}}, nil
}

// UpdateData 更新字典数据。
func (c *Dict) UpdateData(ctx context.Context, req *api.DictDataUpdateReq) (*api.DictDataUpdateRes, error) {
	err := service.UpdateDictData(ctx, &model.DictData{
		Id:       req.Id,
		DictType: req.DictType,
		Label:    req.Label,
		Value:    req.Value,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
	})
	writeOper(ctx, "字典", "更新数据", err)
	if err != nil {
		return nil, err
	}
	return &api.DictDataUpdateRes{Updated: true}, nil
}

// DeleteData 删除字典数据。
func (c *Dict) DeleteData(ctx context.Context, req *api.DictDataDeleteReq) (*api.DictDataDeleteRes, error) {
	err := service.DeleteDictData(ctx, req.Id)
	writeOper(ctx, "字典", "删除数据", err)
	if err != nil {
		return nil, err
	}
	return &api.DictDataDeleteRes{Deleted: true}, nil
}
