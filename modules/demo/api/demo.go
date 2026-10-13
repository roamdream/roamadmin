// Package api demo 模块对外契约（请求/响应 DTO）。
package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PingReq 连通性测试请求。
type PingReq struct {
	g.Meta `path:"/demo/ping" method:"get" tags:"demo" summary:"Connectivity check"`
}

// PingRes 连通性测试响应。
type PingRes struct {
	Module  string `json:"module" dc:"Module name"`
	Version string `json:"version" dc:"Version"`
	Time    string `json:"time" dc:"Current time (RFC3339 UTC)"`
}

// Item 示例数据。
type Item struct {
	Id        int64  `json:"id" dc:"ID"`
	Name      string `json:"name" dc:"Name"`
	Remark    string `json:"remark" dc:"Remark"`
	CreatedAt string `json:"createdAt" dc:"Creation time"`
}

// CreateReq 创建示例数据请求。
type CreateReq struct {
	g.Meta `path:"/demo/item" method:"post" tags:"demo" summary:"Create item"`
	Name   string `json:"name" v:"required|length:1,64#name is required|name length must be 1-64" dc:"Name"`
	Remark string `json:"remark" dc:"Remark"`
}

// CreateRes 创建示例数据响应。
type CreateRes struct {
	Id int64 `json:"id" dc:"New item ID"`
}

// ListReq 分页查询示例数据请求。
type ListReq struct {
	g.Meta `path:"/demo/item" method:"get" tags:"demo" summary:"List items"`
	Page   int    `json:"page" d:"1" v:"min:1#page must be at least 1" dc:"Page number"`
	Size   int    `json:"size" d:"20" v:"min:1|max:200#page size must be at least 1|page size must be at most 200" dc:"Page size"`
	Name   string `json:"name" dc:"Fuzzy search by name"`
}

// ListRes 分页查询示例数据响应。
type ListRes struct {
	List  []*Item `json:"list" dc:"Item list"`
	Total int64   `json:"total" dc:"Total count"`
	Page  int     `json:"page" dc:"Page number"`
	Size  int     `json:"size" dc:"Page size"`
}

// GetReq 查询单条示例数据请求。
type GetReq struct {
	g.Meta `path:"/demo/item/{id}" method:"get" tags:"demo" summary:"Get item"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"ID"`
}

// DeleteReq 删除示例数据请求。
type DeleteReq struct {
	g.Meta `path:"/demo/item/{id}" method:"delete" tags:"demo" summary:"Delete item"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"ID"`
}

// DeleteRes 删除示例数据响应。
type DeleteRes struct {
	Deleted bool `json:"deleted" dc:"Whether deleted"`
}
