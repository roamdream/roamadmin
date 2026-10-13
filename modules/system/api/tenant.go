package api

import "github.com/gogf/gf/v2/frame/g"

// TenantItem 租户列表项。
type TenantItem struct {
	Id        int64  `json:"id"        dc:"Tenant ID"`
	Name      string `json:"name"      dc:"Tenant name"`
	Code      string `json:"code"      dc:"Tenant code"`
	Status    int    `json:"status"    dc:"Status (1 enabled, 0 disabled)"`
	Remark    string `json:"remark"    dc:"Remark"`
	CreatedAt string `json:"createdAt" dc:"Creation time"`
}

// TenantCreateReq 创建租户请求。
type TenantCreateReq struct {
	g.Meta `path:"/system/tenant" method:"post" tags:"system" summary:"Create tenant"`
	Name   string `json:"name" v:"required|length:1,64#tenant name is required|tenant name length must be 1-64" dc:"Tenant name"`
	Code   string `json:"code" v:"required|length:1,64#tenant code is required|tenant code length must be 1-64" dc:"Tenant code"`
	Status int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
	Remark string `json:"remark" dc:"Remark"`
}

// TenantCreateRes 创建租户响应。
type TenantCreateRes struct {
	Id int64 `json:"id" dc:"New tenant ID"`
}

// TenantListReq 分页查询租户请求。
type TenantListReq struct {
	g.Meta `path:"/system/tenant" method:"get" tags:"system" summary:"List tenants"`
	Page   int    `json:"page" d:"1" v:"min:1#page must be at least 1" dc:"Page number"`
	Size   int    `json:"size" d:"20" v:"min:1|max:200#page size must be at least 1|page size must be at most 200" dc:"Page size"`
	Name   string `json:"name" dc:"Fuzzy search by tenant name"`
}

// TenantListRes 分页查询租户响应。
type TenantListRes struct {
	List  []*TenantItem `json:"list"  dc:"Tenant list"`
	Total int64         `json:"total" dc:"Total count"`
	Page  int           `json:"page"  dc:"Page number"`
	Size  int           `json:"size"  dc:"Page size"`
}

// TenantGetReq 查询租户请求。
type TenantGetReq struct {
	g.Meta `path:"/system/tenant/{id}" method:"get" tags:"system" summary:"Get tenant"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Tenant ID"`
}

// TenantGetRes 查询租户响应。
type TenantGetRes struct {
	Tenant *TenantItem `json:"tenant" dc:"Tenant info"`
}

// TenantUpdateReq 更新租户请求。
type TenantUpdateReq struct {
	g.Meta `path:"/system/tenant/{id}" method:"put" tags:"system" summary:"Update tenant"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Tenant ID"`
	Name   string `json:"name" v:"required|length:1,64#tenant name is required|tenant name length must be 1-64" dc:"Tenant name"`
	Status int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
	Remark string `json:"remark" dc:"Remark"`
}

// TenantUpdateRes 更新租户响应。
type TenantUpdateRes struct {
	Updated bool `json:"updated" dc:"Whether updated"`
}

// TenantDeleteReq 删除租户请求。
type TenantDeleteReq struct {
	g.Meta `path:"/system/tenant/{id}" method:"delete" tags:"system" summary:"Delete tenant"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Tenant ID"`
}

// TenantDeleteRes 删除租户响应。
type TenantDeleteRes struct {
	Deleted bool `json:"deleted" dc:"Whether deleted"`
}
