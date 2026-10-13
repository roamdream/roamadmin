package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

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
	list, total, err := service.ListOperLog(ctx, req.Page, req.Size, req.Username, req.Module, req.Status)
	if err != nil {
		return nil, err
	}
	out := make([]*api.OperLogItem, 0, len(list))
	for _, it := range list {
		out = append(out, &api.OperLogItem{
			Id:        it.Id,
			TenantId:  it.TenantId,
			UserId:    it.UserId,
			Username:  it.Username,
			Module:    it.Module,
			Action:    it.Action,
			Method:    it.Method,
			Path:      it.Path,
			IP:        it.IP,
			Status:    it.Status,
			ErrMsg:    it.ErrMsg,
			CreatedAt: it.CreatedAt.String(),
		})
	}
	return &api.OperLogListRes{List: out, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Clear 清空操作日志。
func (c *OperLog) Clear(ctx context.Context, req *api.OperLogClearReq) (*api.OperLogClearRes, error) {
	n, err := service.ClearOperLog(ctx)
	writeOper(ctx, "操作日志", "清空", err)
	if err != nil {
		return nil, err
	}
	return &api.OperLogClearRes{Cleared: n}, nil
}

// recordOper 记录操作日志（登录与关键写操作调用，失败不阻断主流程）。
func recordOper(ctx context.Context, module, action string, status int, errMsg string) {
	var method, path, ip string
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		method, path, ip = r.Method, r.URL.Path, r.GetClientIp()
	}
	service.RecordOper(ctx, &model.OperLog{
		Module:  module,
		Action:  action,
		Method:  method,
		Path:    path,
		IP:      ip,
		Status:  status,
		ErrMsg:  errMsg,
	})
}

// writeOper 记录写操作结果（成功 status=0，失败 status=1 并带错误信息）。
func writeOper(ctx context.Context, module, action string, err error) {
	if err != nil {
		recordOper(ctx, module, action, 1, err.Error())
		return
	}
	recordOper(ctx, module, action, 0, "")
}

// recordLogin 记录登录结果（登录接口无 JWT，需显式传入用户信息）。
func recordLogin(ctx context.Context, status int, errMsg string, userID, tenantID int64, username string) {
	var method, path, ip string
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		method, path, ip = r.Method, r.URL.Path, r.GetClientIp()
	}
	service.RecordOper(ctx, &model.OperLog{
		Module:   "认证",
		Action:   "登录",
		Method:   method,
		Path:     path,
		IP:       ip,
		UserId:   userID,
		TenantId: tenantID,
		Username: username,
		Status:   status,
		ErrMsg:   errMsg,
	})
}
