// Package api demo 模块对外契约（请求/响应 DTO）。
package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PingReq 连通性测试请求。
type PingReq struct {
	g.Meta `path:"/demo/ping" method:"get" tags:"demo" summary:"连通性测试"`
}

// PingRes 连通性测试响应。
type PingRes struct {
	Module  string `json:"module" dc:"模块名"`
	Version string `json:"version" dc:"版本号"`
	Time    string `json:"time" dc:"当前时间"`
}

// Item 示例数据。
type Item struct {
	Id        int64  `json:"id" dc:"ID"`
	Name      string `json:"name" dc:"名称"`
	Remark    string `json:"remark" dc:"备注"`
	CreatedAt string `json:"createdAt" dc:"创建时间"`
}

// CreateReq 创建示例数据请求。
type CreateReq struct {
	g.Meta `path:"/demo/item" method:"post" tags:"demo" summary:"创建示例数据"`
	Name   string `json:"name" v:"required|length:1,64#名称不能为空|名称长度1-64" dc:"名称"`
	Remark string `json:"remark" dc:"备注"`
}

// CreateRes 创建示例数据响应。
type CreateRes struct {
	Id int64 `json:"id" dc:"新数据 ID"`
}

// ListReq 分页查询示例数据请求。
type ListReq struct {
	g.Meta `path:"/demo/item" method:"get" tags:"demo" summary:"分页查询示例数据"`
	Page   int    `json:"page" d:"1" v:"min:1#页码最小为1" dc:"页码"`
	Size   int    `json:"size" d:"20" v:"min:1|max:200#每页最小为1|每页最大为200" dc:"每页条数"`
	Name   string `json:"name" dc:"名称模糊搜索"`
}

// ListRes 分页查询示例数据响应。
type ListRes struct {
	List  []*Item `json:"list" dc:"数据列表"`
	Total int64   `json:"total" dc:"总数"`
	Page  int     `json:"page" dc:"页码"`
	Size  int     `json:"size" dc:"每页条数"`
}

// GetReq 查询单条示例数据请求。
type GetReq struct {
	g.Meta `path:"/demo/item/{id}" method:"get" tags:"demo" summary:"查询单条示例数据"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"ID"`
}

// DeleteReq 删除示例数据请求。
type DeleteReq struct {
	g.Meta `path:"/demo/item/{id}" method:"delete" tags:"demo" summary:"删除示例数据"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"ID"`
}

// DeleteRes 删除示例数据响应。
type DeleteRes struct {
	Deleted bool `json:"deleted" dc:"是否已删除"`
}
