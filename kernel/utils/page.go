package utils

// 分页默认值。
const (
	DefaultPageSize = 20
	MaxPageSize     = 200
)

// PageReq 通用分页请求。
type PageReq struct {
	Page int `json:"page" d:"1"  v:"min:1#页码最小为1"`
	Size int `json:"size" d:"20" v:"min:1|max:200#每页最小为1|每页最大为200"`
}

// Normalize 归一化非法分页参数。
func (p *PageReq) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = DefaultPageSize
	}
	if p.Size > MaxPageSize {
		p.Size = MaxPageSize
	}
}

// Offset 返回 SQL 查询偏移量（调用前请先 Normalize）。
func (p *PageReq) Offset() int { return (p.Page - 1) * p.Size }

// PageRes 通用分页响应。
type PageRes[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

// NewPageRes 构造分页响应。
func NewPageRes[T any](list []T, total int64, page, size int) *PageRes[T] {
	return &PageRes[T]{List: list, Total: total, Page: page, Size: size}
}
