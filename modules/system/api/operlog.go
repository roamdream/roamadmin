package api

import "github.com/gogf/gf/v2/frame/g"

// OperLogItem 操作日志列表项。
type OperLogItem struct {
	Id        int64  `json:"id"        dc:"日志 ID"`
	TenantId  int64  `json:"tenantId"  dc:"租户 ID"`
	UserId    int64  `json:"userId"    dc:"用户 ID"`
	Username  string `json:"username"  dc:"用户名"`
	Module    string `json:"module"    dc:"模块"`
	Action    string `json:"action"    dc:"操作"`
	Method    string `json:"method"    dc:"请求方法"`
	Path      string `json:"path"      dc:"请求路径"`
	IP        string `json:"ip"        dc:"客户端 IP"`
	Status    int    `json:"status"    dc:"状态（0成功 1失败）"`
	ErrMsg    string `json:"errMsg"    dc:"错误信息"`
	CreatedAt string `json:"createdAt" dc:"创建时间"`
}

// OperLogListReq 分页查询操作日志请求。
type OperLogListReq struct {
	g.Meta   `path:"/system/operlog" method:"get" tags:"系统" summary:"分页查询操作日志"`
	Page     int    `json:"page" d:"1" v:"min:1#页码最小为1" dc:"页码"`
	Size     int    `json:"size" d:"20" v:"min:1|max:200#每页最小为1|每页最大为200" dc:"每页条数"`
	Username string `json:"username" dc:"用户名模糊搜索"`
	Module   string `json:"module" dc:"模块筛选"`
	Status   *int   `json:"status" dc:"状态筛选（0成功 1失败）"`
}

// OperLogListRes 分页查询操作日志响应。
type OperLogListRes struct {
	List  []*OperLogItem `json:"list"  dc:"操作日志列表"`
	Total int64          `json:"total" dc:"总数"`
	Page  int            `json:"page"  dc:"页码"`
	Size  int            `json:"size"  dc:"每页条数"`
}

// OperLogClearReq 清空操作日志请求。
type OperLogClearReq struct {
	g.Meta `path:"/system/operlog" method:"delete" tags:"系统" summary:"清空操作日志"`
}

// OperLogClearRes 清空操作日志响应。
type OperLogClearRes struct {
	Cleared int64 `json:"cleared" dc:"清空条数"`
}
