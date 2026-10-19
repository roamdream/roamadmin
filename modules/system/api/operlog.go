package api

import "github.com/gogf/gf/v2/frame/g"

// OperLogItem 操作日志列表项（不含 req_body/detail 大字段，详情另查）。
type OperLogItem struct {
	Id         int64  `json:"id"         dc:"Log ID"`
	TenantId   int64  `json:"tenantId"   dc:"Tenant ID"`
	UserId     int64  `json:"userId"     dc:"User ID"`
	Username   string `json:"username"   dc:"Username"`
	Module     string `json:"module"     dc:"Module"`
	Action     string `json:"action"     dc:"Action"`
	Method     string `json:"method"     dc:"HTTP method"`
	Path       string `json:"path"       dc:"Request path"`
	Query      string `json:"query"      dc:"Query string"`
	IP         string `json:"ip"         dc:"Client IP"`
	UserAgent  string `json:"userAgent"  dc:"User agent"`
	Referer    string `json:"referer"    dc:"Referer"`
	Status     int    `json:"status"     dc:"Status (0 success, 1 failure)"`
	RespCode   int    `json:"respCode"   dc:"HTTP status code"`
	ErrMsg     string `json:"errMsg"     dc:"Error message"`
	DurationMs int64  `json:"durationMs" dc:"Duration in milliseconds"`
	TargetType string `json:"targetType" dc:"Target resource type"`
	TargetId   int64  `json:"targetId"   dc:"Target resource ID"`
	TraceId    string `json:"traceId"    dc:"Trace ID"`
	CreatedAt  string `json:"createdAt"  dc:"Creation time"`
}

// OperLogListReq 分页查询操作日志请求。
type OperLogListReq struct {
	g.Meta     `path:"/system/operlog" method:"get" tags:"system" summary:"List operation logs"`
	Page       int    `json:"page" d:"1" v:"min:1#page must be at least 1" dc:"Page number"`
	Size       int    `json:"size" d:"20" v:"min:1|max:200#page size must be at least 1|page size must be at most 200" dc:"Page size"`
	Username   string `json:"username" dc:"Fuzzy search by username"`
	Module     string `json:"module" dc:"Filter by module"`
	Action     string `json:"action" dc:"Filter by action"`
	Method     string `json:"method" dc:"Filter by HTTP method"`
	TargetType string `json:"targetType" dc:"Filter by target type"`
	TraceId    string `json:"traceId" dc:"Filter by trace ID"`
	Status     *int   `json:"status" dc:"Filter by status (0 success, 1 failure)"`
	StartAt    string `json:"startAt" dc:"Start time (RFC3339 UTC)"`
	EndAt      string `json:"endAt" dc:"End time (RFC3339 UTC)"`
}

// OperLogListRes 分页查询操作日志响应。
type OperLogListRes struct {
	List  []*OperLogItem `json:"list"  dc:"Operation log list"`
	Total int64          `json:"total" dc:"Total count"`
	Page  int            `json:"page"  dc:"Page number"`
	Size  int            `json:"size"  dc:"Page size"`
}

// OperLogDetailReq 操作日志详情请求。
type OperLogDetailReq struct {
	g.Meta `path:"/system/operlog/{id}" method:"get" tags:"system" summary:"Get operation log detail"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#id is required|id must be at least 1" dc:"Log ID"`
}

// OperLogDetailRes 操作日志详情响应（含请求体与变更详情大字段）。
type OperLogDetailRes struct {
	Item    *OperLogItem `json:"item"    dc:"Log summary"`
	ReqBody string       `json:"reqBody" dc:"Request body (sensitive fields masked)"`
	Detail  string       `json:"detail"  dc:"Change detail (JSON diff, if any)"`
}

// OperLogClearReq 清空操作日志请求。
type OperLogClearReq struct {
	g.Meta `path:"/system/operlog" method:"delete" tags:"system" summary:"Clear operation logs"`
}

// OperLogClearRes 清空操作日志响应。
type OperLogClearRes struct {
	Cleared int64 `json:"cleared" dc:"Number of cleared logs"`
}
