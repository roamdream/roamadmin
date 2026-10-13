package api

import "github.com/gogf/gf/v2/frame/g"

// TenantItem 租户列表项。
type TenantItem struct {
	Id        int64  `json:"id"        dc:"租户 ID"`
	Name      string `json:"name"      dc:"租户名称"`
	Code      string `json:"code"      dc:"租户编码"`
	Status    int    `json:"status"    dc:"状态（1启用 0禁用）"`
	Remark    string `json:"remark"    dc:"备注"`
	CreatedAt string `json:"createdAt" dc:"创建时间"`
}

// TenantCreateReq 创建租户请求。
type TenantCreateReq struct {
	g.Meta `path:"/system/tenant" method:"post" tags:"系统" summary:"创建租户"`
	Name   string `json:"name" v:"required|length:1,64#租户名称不能为空|租户名称长度1-64" dc:"租户名称"`
	Code   string `json:"code" v:"required|length:1,64#租户编码不能为空|租户编码长度1-64" dc:"租户编码"`
	Status int    `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
	Remark string `json:"remark" dc:"备注"`
}

// TenantCreateRes 创建租户响应。
type TenantCreateRes struct {
	Id int64 `json:"id" dc:"新租户 ID"`
}

// TenantListReq 分页查询租户请求。
type TenantListReq struct {
	g.Meta `path:"/system/tenant" method:"get" tags:"系统" summary:"分页查询租户"`
	Page   int    `json:"page" d:"1" v:"min:1#页码最小为1" dc:"页码"`
	Size   int    `json:"size" d:"20" v:"min:1|max:200#每页最小为1|每页最大为200" dc:"每页条数"`
	Name   string `json:"name" dc:"租户名称模糊搜索"`
}

// TenantListRes 分页查询租户响应。
type TenantListRes struct {
	List  []*TenantItem `json:"list"  dc:"租户列表"`
	Total int64         `json:"total" dc:"总数"`
	Page  int           `json:"page"  dc:"页码"`
	Size  int           `json:"size"  dc:"每页条数"`
}

// TenantGetReq 查询租户请求。
type TenantGetReq struct {
	g.Meta `path:"/system/tenant/{id}" method:"get" tags:"系统" summary:"查询租户"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"租户 ID"`
}

// TenantGetRes 查询租户响应。
type TenantGetRes struct {
	Tenant *TenantItem `json:"tenant" dc:"租户信息"`
}

// TenantUpdateReq 更新租户请求。
type TenantUpdateReq struct {
	g.Meta `path:"/system/tenant/{id}" method:"put" tags:"系统" summary:"更新租户"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"租户 ID"`
	Name   string `json:"name" v:"required|length:1,64#租户名称不能为空|租户名称长度1-64" dc:"租户名称"`
	Status int    `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
	Remark string `json:"remark" dc:"备注"`
}

// TenantUpdateRes 更新租户响应。
type TenantUpdateRes struct {
	Updated bool `json:"updated" dc:"是否已更新"`
}

// TenantDeleteReq 删除租户请求。
type TenantDeleteReq struct {
	g.Meta `path:"/system/tenant/{id}" method:"delete" tags:"系统" summary:"删除租户"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"租户 ID"`
}

// TenantDeleteRes 删除租户响应。
type TenantDeleteRes struct {
	Deleted bool `json:"deleted" dc:"是否已删除"`
}
