// Package api 跨模块共享的对外契约（请求/响应结构）。
//
// 各业务模块在自身 modules/<name>/api 中定义专有 DTO；此处只放被多个模块
// 复用的通用结构，避免重复定义。
package api

// IdReq 通用主键请求（按 ID 查询 / 删除单条）。
type IdReq struct {
	Id int64 `json:"id" v:"required|min:1#ID is required|ID is invalid"`
}

// IdsReq 通用主键列表请求（批量删除）。
type IdsReq struct {
	Ids []int64 `json:"ids" v:"required#ID list is required"`
}
