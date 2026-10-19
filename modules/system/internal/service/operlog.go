package service

import (
	"context"
	"time"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/system/internal/consts"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
)

// OperLogFilter 操作日志查询过滤条件。
type OperLogFilter struct {
	Username   string
	Module     string
	Action     string
	Method     string
	TargetType string
	TraceId    string
	Status     *int
	StartAt    *time.Time
	EndAt      *time.Time
}

// ListOperLog 分页查询操作日志（按租户隔离）。
func ListOperLog(ctx context.Context, page, size int, f *OperLogFilter) ([]*model.OperLog, int64, error) {
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
	m := ra.DB().Model(consts.TableOperLog).Where("tenant_id", tenantID)
	if f != nil {
		if f.Username != "" {
			m = m.WhereLike("username", "%"+f.Username+"%")
		}
		if f.Module != "" {
			m = m.Where("module", f.Module)
		}
		if f.Action != "" {
			m = m.Where("action", f.Action)
		}
		if f.Method != "" {
			m = m.Where("method", f.Method)
		}
		if f.TargetType != "" {
			m = m.Where("target_type", f.TargetType)
		}
		if f.TraceId != "" {
			m = m.Where("trace_id", f.TraceId)
		}
		if f.Status != nil {
			m = m.Where("status", *f.Status)
		}
		if f.StartAt != nil {
			m = m.WhereGTE("created_at", *f.StartAt)
		}
		if f.EndAt != nil {
			m = m.WhereLTE("created_at", *f.EndAt)
		}
	}
	total, err := m.Count()
	if err != nil || total == 0 {
		return []*model.OperLog{}, int64(total), err
	}
	var list []*model.OperLog
	if err = m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, errcode.Wrap(errcode.CodeDB, err)
	}
	return list, int64(total), nil
}

// GetOperLog 查询单条操作日志详情（租户隔离，含 req_body/detail 大字段）。
func GetOperLog(ctx context.Context, id int64) (*model.OperLog, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return nil, errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	var out *model.OperLog
	if err := ra.DB().Model(consts.TableOperLog).Where("tenant_id", tenantID).Where("id", id).Scan(&out); err != nil {
		return nil, errcode.Wrap(errcode.CodeDB, err)
	}
	return out, nil
}

// ClearOperLog 清空当前租户的操作日志，返回清空条数。
func ClearOperLog(ctx context.Context) (int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return 0, errcode.New(errcode.CodeTenantMissing, "missing tenant context")
	}
	res, err := ra.DB().Model(consts.TableOperLog).Where("tenant_id", tenantID).Delete()
	if err != nil {
		return 0, errcode.Wrap(errcode.CodeDB, err)
	}
	n, _ := res.RowsAffected()
	return n, nil
}
