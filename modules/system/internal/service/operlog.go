package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/system/internal/consts"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
)

// ListOperLog 分页查询操作日志（按租户隔离）。
func ListOperLog(ctx context.Context, page, size int, username, module string, status *int) ([]*model.OperLog, int64, error) {
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
	m := ra.DB().Model(consts.TableOperLog).Where("tenant_id", tenantID)
	if username != "" {
		m = m.WhereLike("username", "%"+username+"%")
	}
	if module != "" {
		m = m.Where("module", module)
	}
	if status != nil {
		m = m.Where("status", *status)
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

// ClearOperLog 清空当前租户的操作日志，返回清空条数。
func ClearOperLog(ctx context.Context) (int64, error) {
	tenantID, ok := auth.TenantFrom(ctx)
	if !ok {
		return 0, errcode.New(errcode.CodeTenantMissing, "缺失租户上下文")
	}
	res, err := ra.DB().Model(consts.TableOperLog).Where("tenant_id", tenantID).Delete()
	if err != nil {
		return 0, errcode.Wrap(errcode.CodeDB, err)
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// RecordOper 记录一条操作日志（登录与关键写操作调用；失败仅记录，不阻断主流程）。
func RecordOper(ctx context.Context, in *model.OperLog) {
	if c, ok := auth.UserFrom(ctx); ok {
		in.TenantId, in.UserId, in.Username = c.TenantID, c.UserID, c.Username
	}
	in.Id = ra.DB().NextID()
	in.CreatedAt = gtime.Now()
	if _, err := ra.DB().Model(consts.TableOperLog).Data(g.Map{
		"id":         in.Id,
		"tenant_id":  in.TenantId,
		"user_id":    in.UserId,
		"username":   in.Username,
		"module":     in.Module,
		"action":     in.Action,
		"method":     in.Method,
		"path":       in.Path,
		"ip":         in.IP,
		"status":     in.Status,
		"err_msg":    in.ErrMsg,
		"created_at": in.CreatedAt,
	}).Insert(); err != nil {
		ra.Log().Errorf(ctx, "[system] record operlog failed: %v", err)
	}
}
