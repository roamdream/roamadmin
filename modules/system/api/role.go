package api

import "github.com/gogf/gf/v2/frame/g"

// RoleItem 角色列表项。
type RoleItem struct {
	Id        int64  `json:"id"        dc:"角色 ID"`
	TenantId  int64  `json:"tenantId"  dc:"租户 ID"`
	Name      string `json:"name"      dc:"角色名称"`
	Code      string `json:"code"      dc:"角色编码"`
	Status    int    `json:"status"    dc:"状态（1启用 0禁用）"`
	Remark    string `json:"remark"    dc:"备注"`
	CreatedAt string `json:"createdAt" dc:"创建时间"`
}

// RoleCreateReq 创建角色请求。
type RoleCreateReq struct {
	g.Meta `path:"/system/role" method:"post" tags:"系统" summary:"创建角色"`
	Name   string `json:"name" v:"required|length:1,64#角色名称不能为空|角色名称长度1-64" dc:"角色名称"`
	Code   string `json:"code" v:"required|length:1,64#角色编码不能为空|角色编码长度1-64" dc:"角色编码"`
	Status int    `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
	Remark string `json:"remark" dc:"备注"`
}

// RoleCreateRes 创建角色响应。
type RoleCreateRes struct {
	Id int64 `json:"id" dc:"新角色 ID"`
}

// RoleListReq 分页查询角色请求。
type RoleListReq struct {
	g.Meta `path:"/system/role" method:"get" tags:"系统" summary:"分页查询角色"`
	Page   int    `json:"page" d:"1" v:"min:1#页码最小为1" dc:"页码"`
	Size   int    `json:"size" d:"20" v:"min:1|max:200#每页最小为1|每页最大为200" dc:"每页条数"`
	Name   string `json:"name" dc:"角色名称模糊搜索"`
}

// RoleListRes 分页查询角色响应。
type RoleListRes struct {
	List  []*RoleItem `json:"list"  dc:"角色列表"`
	Total int64       `json:"total" dc:"总数"`
	Page  int         `json:"page"  dc:"页码"`
	Size  int         `json:"size"  dc:"每页条数"`
}

// RoleGetReq 查询角色请求。
type RoleGetReq struct {
	g.Meta `path:"/system/role/{id}" method:"get" tags:"系统" summary:"查询角色"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"角色 ID"`
}

// RoleGetRes 查询角色响应。
type RoleGetRes struct {
	Role    *RoleItem `json:"role"    dc:"角色信息"`
	MenuIds []int64   `json:"menuIds" dc:"已分配菜单 ID 列表"`
}

// RoleUpdateReq 更新角色请求。
type RoleUpdateReq struct {
	g.Meta `path:"/system/role/{id}" method:"put" tags:"系统" summary:"更新角色"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"角色 ID"`
	Name   string `json:"name" v:"required|length:1,64#角色名称不能为空|角色名称长度1-64" dc:"角色名称"`
	Status int    `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
	Remark string `json:"remark" dc:"备注"`
}

// RoleUpdateRes 更新角色响应。
type RoleUpdateRes struct {
	Updated bool `json:"updated" dc:"是否已更新"`
}

// RoleDeleteReq 删除角色请求。
type RoleDeleteReq struct {
	g.Meta `path:"/system/role/{id}" method:"delete" tags:"系统" summary:"删除角色"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"角色 ID"`
}

// RoleDeleteRes 删除角色响应。
type RoleDeleteRes struct {
	Deleted bool `json:"deleted" dc:"是否已删除"`
}

// RoleAssignMenuReq 为角色分配菜单权限请求。
type RoleAssignMenuReq struct {
	g.Meta  `path:"/system/role/{id}/menu" method:"put" tags:"系统" summary:"为角色分配菜单权限"`
	Id      int64   `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"角色 ID"`
	MenuIds []int64 `json:"menuIds" dc:"菜单 ID 列表"`
}

// RoleAssignMenuRes 为角色分配菜单权限响应。
type RoleAssignMenuRes struct {
	Assigned bool `json:"assigned" dc:"是否已分配"`
}
