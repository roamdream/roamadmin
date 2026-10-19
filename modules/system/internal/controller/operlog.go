package controller

import (
	"context"
	"time"

	"github.com/roamdream/roamadmin/kernel/audit"
	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/utils"
	"github.com/roamdream/roamadmin/modules/system/api"
	"github.com/roamdream/roamadmin/modules/system/internal/model"
	"github.com/roamdream/roamadmin/modules/system/internal/service"
)

// OperLog 操作日志控制器。
type OperLog struct{}

// NewOperLog 创建操作日志控制器。
func NewOperLog() *OperLog { return &OperLog{} }

// List 分页查询操作日志。
func (c *OperLog) List(ctx context.Context, req *api.OperLogListReq) (*api.OperLogListRes, error) {
	f := &service.OperLogFilter{
		Username:   req.Username,
		Module:     req.Module,
		Action:     req.Action,
		Method:     req.Method,
		TargetType: req.TargetType,
		TraceId:    req.TraceId,
		Status:     req.Status,
	}
	var err error
	if req.StartAt != "" {
		if f.StartAt, err = parseTime(req.StartAt); err != nil {
			return nil, err
		}
	}
	if req.EndAt != "" {
		if f.EndAt, err = parseTime(req.EndAt); err != nil {
			return nil, err
		}
	}
	list, total, err := service.ListOperLog(ctx, req.Page, req.Size, f)
	if err != nil {
		return nil, err
	}
	out := make([]*api.OperLogItem, 0, len(list))
	for _, it := range list {
		out = append(out, toOperLogItem(it))
	}
	return &api.OperLogListRes{List: out, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Detail 查询单条操作日志详情（含请求体与变更详情）。
func (c *OperLog) Detail(ctx context.Context, req *api.OperLogDetailReq) (*api.OperLogDetailRes, error) {
	it, err := service.GetOperLog(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if it == nil {
		return nil, errcode.New(errcode.CodeNotFound, "operation log not found")
	}
	return &api.OperLogDetailRes{
		Item:    toOperLogItem(it),
		ReqBody: it.ReqBody,
		Detail:  it.Detail,
	}, nil
}

// Clear 清空操作日志。
func (c *OperLog) Clear(ctx context.Context, req *api.OperLogClearReq) (*api.OperLogClearRes, error) {
	n, err := service.ClearOperLog(ctx)
	writeOper(ctx, "operlog", "clear", err)
	if err != nil {
		return nil, err
	}
	return &api.OperLogClearRes{Cleared: n}, nil
}

// writeOper 标记写操作的 module/action 语义；其余字段与结果由审计中间件自动采集。
func writeOper(ctx context.Context, module, action string, _ error) {
	audit.Annotate(ctx, module, action)
}

// recordLogin 记录登录结果（登录接口无 JWT，需显式传入用户信息）。
func recordLogin(ctx context.Context, status int, errMsg string, userID, tenantID int64, username string) {
	audit.Record(ctx, &audit.Entry{
		Module:   "auth",
		Action:   "login",
		Status:   status,
		ErrMsg:   errMsg,
		UserID:   userID,
		TenantID: tenantID,
		Username: username,
	})
}

// toOperLogItem 将模型转为列表项 DTO（不含 req_body/detail 大字段）。
func toOperLogItem(it *model.OperLog) *api.OperLogItem {
	return &api.OperLogItem{
		Id:         it.Id,
		TenantId:   it.TenantId,
		UserId:     it.UserId,
		Username:   it.Username,
		Module:     it.Module,
		Action:     it.Action,
		Method:     it.Method,
		Path:       it.Path,
		Query:      it.Query,
		IP:         it.IP,
		UserAgent:  it.UserAgent,
		Referer:    it.Referer,
		Status:     it.Status,
		RespCode:   it.RespCode,
		ErrMsg:     it.ErrMsg,
		DurationMs: it.DurationMs,
		TargetType: it.TargetType,
		TargetId:   it.TargetId,
		TraceId:    it.TraceId,
		CreatedAt:  utils.TimeRFC3339(it.CreatedAt),
	}
}

// parseTime 解析 RFC3339 时间参数。
func parseTime(s string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, errcode.New(errcode.CodeParam, "invalid time, expect RFC3339 (e.g. 2026-10-18T00:00:00Z)")
	}
	return &t, nil
}
