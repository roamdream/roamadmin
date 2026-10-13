package api

import "github.com/gogf/gf/v2/frame/g"

// MenuItem 菜单列表项。
type MenuItem struct {
	Id        int64  `json:"id"        dc:"Menu ID"`
	TenantId  int64  `json:"tenantId"  dc:"Tenant ID"`
	ParentId  int64  `json:"parentId"  dc:"Parent menu ID"`
	Name      string `json:"name"      dc:"Menu name"`
	Type      string `json:"type"      dc:"Type (dir/menu/button)"`
	Perm      string `json:"perm"      dc:"Permission code"`
	Path      string `json:"path"      dc:"Route path"`
	Component string `json:"component" dc:"Component path"`
	Icon      string `json:"icon"      dc:"Icon"`
	Sort      int    `json:"sort"      dc:"Sort order"`
	Status    int    `json:"status"    dc:"Status (1 enabled, 0 disabled)"`
	CreatedAt string `json:"createdAt" dc:"Creation time"`
}

// MenuTreeItem 菜单树节点。
type MenuTreeItem struct {
	*MenuItem
	Children []*MenuTreeItem `json:"children" dc:"Child menus"`
}

// MenuCreateReq 创建菜单请求。
type MenuCreateReq struct {
	g.Meta    `path:"/system/menu" method:"post" tags:"system" summary:"Create menu"`
	ParentId  int64  `json:"parentId" d:"0" dc:"Parent menu ID (0 for root)"`
	Name      string `json:"name" v:"required|length:1,64#menu name is required|menu name length must be 1-64" dc:"Menu name"`
	Type      string `json:"type" v:"required|in:dir,menu,button#type is required|type is invalid" dc:"Type (dir/menu/button)"`
	Perm      string `json:"perm" dc:"Permission code"`
	Path      string `json:"path" dc:"Route path"`
	Component string `json:"component" dc:"Component path"`
	Icon      string `json:"icon" dc:"Icon"`
	Sort      int    `json:"sort" d:"0" dc:"Sort order"`
	Status    int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
}

// MenuCreateRes 创建菜单响应。
type MenuCreateRes struct {
	Id int64 `json:"id" dc:"New menu ID"`
}

// MenuListReq 查询菜单树请求。
type MenuListReq struct {
	g.Meta `path:"/system/menu" method:"get" tags:"system" summary:"List menu tree"`
	Name   string `json:"name" dc:"Fuzzy search by menu name"`
}

// MenuListRes 查询菜单树响应。
type MenuListRes struct {
	List []*MenuTreeItem `json:"list" dc:"Menu tree"`
}

// MenuGetReq 查询菜单请求。
type MenuGetReq struct {
	g.Meta `path:"/system/menu/{id}" method:"get" tags:"system" summary:"Get menu"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Menu ID"`
}

// MenuGetRes 查询菜单响应。
type MenuGetRes struct {
	Menu *MenuItem `json:"menu" dc:"Menu info"`
}

// MenuUpdateReq 更新菜单请求。
type MenuUpdateReq struct {
	g.Meta    `path:"/system/menu/{id}" method:"put" tags:"system" summary:"Update menu"`
	Id        int64  `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Menu ID"`
	ParentId  int64  `json:"parentId" d:"0" dc:"Parent menu ID (0 for root)"`
	Name      string `json:"name" v:"required|length:1,64#menu name is required|menu name length must be 1-64" dc:"Menu name"`
	Type      string `json:"type" v:"required|in:dir,menu,button#type is required|type is invalid" dc:"Type (dir/menu/button)"`
	Perm      string `json:"perm" dc:"Permission code"`
	Path      string `json:"path" dc:"Route path"`
	Component string `json:"component" dc:"Component path"`
	Icon      string `json:"icon" dc:"Icon"`
	Sort      int    `json:"sort" d:"0" dc:"Sort order"`
	Status    int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
}

// MenuUpdateRes 更新菜单响应。
type MenuUpdateRes struct {
	Updated bool `json:"updated" dc:"Whether updated"`
}

// MenuDeleteReq 删除菜单请求。
type MenuDeleteReq struct {
	g.Meta `path:"/system/menu/{id}" method:"delete" tags:"system" summary:"Delete menu"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Menu ID"`
}

// MenuDeleteRes 删除菜单响应。
type MenuDeleteRes struct {
	Deleted bool `json:"deleted" dc:"Whether deleted"`
}
