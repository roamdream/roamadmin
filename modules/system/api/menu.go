package api

import "github.com/gogf/gf/v2/frame/g"

// MenuItem 菜单列表项。
type MenuItem struct {
	Id        int64  `json:"id"        dc:"菜单 ID"`
	TenantId  int64  `json:"tenantId"  dc:"租户 ID"`
	ParentId  int64  `json:"parentId"  dc:"父菜单 ID"`
	Name      string `json:"name"      dc:"菜单名称"`
	Type      string `json:"type"      dc:"类型（dir/menu/button）"`
	Perm      string `json:"perm"      dc:"权限码"`
	Path      string `json:"path"      dc:"路由地址"`
	Component string `json:"component" dc:"组件路径"`
	Icon      string `json:"icon"      dc:"图标"`
	Sort      int    `json:"sort"      dc:"排序"`
	Status    int    `json:"status"    dc:"状态（1启用 0禁用）"`
	CreatedAt string `json:"createdAt" dc:"创建时间"`
}

// MenuTreeItem 菜单树节点。
type MenuTreeItem struct {
	*MenuItem
	Children []*MenuTreeItem `json:"children" dc:"子菜单"`
}

// MenuCreateReq 创建菜单请求。
type MenuCreateReq struct {
	g.Meta    `path:"/system/menu" method:"post" tags:"系统" summary:"创建菜单"`
	ParentId  int64  `json:"parentId" d:"0" dc:"父菜单 ID（0 为根）"`
	Name      string `json:"name" v:"required|length:1,64#菜单名称不能为空|菜单名称长度1-64" dc:"菜单名称"`
	Type      string `json:"type" v:"required|in:dir,menu,button#类型不能为空|类型不合法" dc:"类型（dir/menu/button）"`
	Perm      string `json:"perm" dc:"权限码"`
	Path      string `json:"path" dc:"路由地址"`
	Component string `json:"component" dc:"组件路径"`
	Icon      string `json:"icon" dc:"图标"`
	Sort      int    `json:"sort" d:"0" dc:"排序"`
	Status    int    `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
}

// MenuCreateRes 创建菜单响应。
type MenuCreateRes struct {
	Id int64 `json:"id" dc:"新菜单 ID"`
}

// MenuListReq 查询菜单树请求。
type MenuListReq struct {
	g.Meta `path:"/system/menu" method:"get" tags:"系统" summary:"查询菜单树"`
	Name   string `json:"name" dc:"菜单名称模糊搜索"`
}

// MenuListRes 查询菜单树响应。
type MenuListRes struct {
	List []*MenuTreeItem `json:"list" dc:"菜单树"`
}

// MenuGetReq 查询菜单请求。
type MenuGetReq struct {
	g.Meta `path:"/system/menu/{id}" method:"get" tags:"系统" summary:"查询菜单"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"菜单 ID"`
}

// MenuGetRes 查询菜单响应。
type MenuGetRes struct {
	Menu *MenuItem `json:"menu" dc:"菜单信息"`
}

// MenuUpdateReq 更新菜单请求。
type MenuUpdateReq struct {
	g.Meta    `path:"/system/menu/{id}" method:"put" tags:"系统" summary:"更新菜单"`
	Id        int64  `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"菜单 ID"`
	ParentId  int64  `json:"parentId" d:"0" dc:"父菜单 ID（0 为根）"`
	Name      string `json:"name" v:"required|length:1,64#菜单名称不能为空|菜单名称长度1-64" dc:"菜单名称"`
	Type      string `json:"type" v:"required|in:dir,menu,button#类型不能为空|类型不合法" dc:"类型（dir/menu/button）"`
	Perm      string `json:"perm" dc:"权限码"`
	Path      string `json:"path" dc:"路由地址"`
	Component string `json:"component" dc:"组件路径"`
	Icon      string `json:"icon" dc:"图标"`
	Sort      int    `json:"sort" d:"0" dc:"排序"`
	Status    int    `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
}

// MenuUpdateRes 更新菜单响应。
type MenuUpdateRes struct {
	Updated bool `json:"updated" dc:"是否已更新"`
}

// MenuDeleteReq 删除菜单请求。
type MenuDeleteReq struct {
	g.Meta `path:"/system/menu/{id}" method:"delete" tags:"系统" summary:"删除菜单"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"菜单 ID"`
}

// MenuDeleteRes 删除菜单响应。
type MenuDeleteRes struct {
	Deleted bool `json:"deleted" dc:"是否已删除"`
}
