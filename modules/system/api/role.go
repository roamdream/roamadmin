package api

import "github.com/gogf/gf/v2/frame/g"

// RoleItem 角色列表项。
type RoleItem struct {
	Id        int64  `json:"id"        dc:"Role ID"`
	TenantId  int64  `json:"tenantId"  dc:"Tenant ID"`
	Name      string `json:"name"      dc:"Role name"`
	Code      string `json:"code"      dc:"Role code"`
	Status    int    `json:"status"    dc:"Status (1 enabled, 0 disabled)"`
	Remark    string `json:"remark"    dc:"Remark"`
	CreatedAt string `json:"createdAt" dc:"Creation time"`
}

// RoleCreateReq 创建角色请求。
type RoleCreateReq struct {
	g.Meta `path:"/system/role" method:"post" tags:"system" summary:"Create role"`
	Name   string `json:"name" v:"required|length:1,64#role name is required|role name length must be 1-64" dc:"Role name"`
	Code   string `json:"code" v:"required|length:1,64#role code is required|role code length must be 1-64" dc:"Role code"`
	Status int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
	Remark string `json:"remark" dc:"Remark"`
}

// RoleCreateRes 创建角色响应。
type RoleCreateRes struct {
	Id int64 `json:"id" dc:"New role ID"`
}

// RoleListReq 分页查询角色请求。
type RoleListReq struct {
	g.Meta `path:"/system/role" method:"get" tags:"system" summary:"List roles"`
	Page   int    `json:"page" d:"1" v:"min:1#page must be at least 1" dc:"Page number"`
	Size   int    `json:"size" d:"20" v:"min:1|max:200#page size must be at least 1|page size must be at most 200" dc:"Page size"`
	Name   string `json:"name" dc:"Fuzzy search by role name"`
}

// RoleListRes 分页查询角色响应。
type RoleListRes struct {
	List  []*RoleItem `json:"list"  dc:"Role list"`
	Total int64       `json:"total" dc:"Total count"`
	Page  int         `json:"page"  dc:"Page number"`
	Size  int         `json:"size"  dc:"Page size"`
}

// RoleGetReq 查询角色请求。
type RoleGetReq struct {
	g.Meta `path:"/system/role/{id}" method:"get" tags:"system" summary:"Get role"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Role ID"`
}

// RoleGetRes 查询角色响应。
type RoleGetRes struct {
	Role    *RoleItem `json:"role"    dc:"Role info"`
	MenuIds []int64   `json:"menuIds" dc:"Assigned menu ID list"`
}

// RoleUpdateReq 更新角色请求。
type RoleUpdateReq struct {
	g.Meta `path:"/system/role/{id}" method:"put" tags:"system" summary:"Update role"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Role ID"`
	Name   string `json:"name" v:"required|length:1,64#role name is required|role name length must be 1-64" dc:"Role name"`
	Status int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
	Remark string `json:"remark" dc:"Remark"`
}

// RoleUpdateRes 更新角色响应。
type RoleUpdateRes struct {
	Updated bool `json:"updated" dc:"Whether updated"`
}

// RoleDeleteReq 删除角色请求。
type RoleDeleteReq struct {
	g.Meta `path:"/system/role/{id}" method:"delete" tags:"system" summary:"Delete role"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Role ID"`
}

// RoleDeleteRes 删除角色响应。
type RoleDeleteRes struct {
	Deleted bool `json:"deleted" dc:"Whether deleted"`
}

// RoleAssignMenuReq 为角色分配菜单权限请求。
type RoleAssignMenuReq struct {
	g.Meta  `path:"/system/role/{id}/menu" method:"put" tags:"system" summary:"Assign menus to role"`
	Id      int64   `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Role ID"`
	MenuIds []int64 `json:"menuIds" dc:"Menu ID list"`
}

// RoleAssignMenuRes 为角色分配菜单权限响应。
type RoleAssignMenuRes struct {
	Assigned bool `json:"assigned" dc:"Whether assigned"`
}
