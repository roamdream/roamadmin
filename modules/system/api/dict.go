package api

import "github.com/gogf/gf/v2/frame/g"

// DictTypeItem 字典类型列表项。
type DictTypeItem struct {
	Id        int64  `json:"id"        dc:"Dict type ID"`
	TenantId  int64  `json:"tenantId"  dc:"Tenant ID"`
	Name      string `json:"name"      dc:"Dict name"`
	Type      string `json:"type"      dc:"Dict type code"`
	Status    int    `json:"status"    dc:"Status (1 enabled, 0 disabled)"`
	Remark    string `json:"remark"    dc:"Remark"`
	CreatedAt string `json:"createdAt" dc:"Creation time"`
}

// DictTypeCreateReq 创建字典类型请求。
type DictTypeCreateReq struct {
	g.Meta `path:"/system/dict/type" method:"post" tags:"system" summary:"Create dict type"`
	Name   string `json:"name" v:"required|length:1,64#dict name is required|dict name length must be 1-64" dc:"Dict name"`
	Type   string `json:"type" v:"required|length:1,64#dict type code is required|dict type code length must be 1-64" dc:"Dict type code"`
	Status int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
	Remark string `json:"remark" dc:"Remark"`
}

// DictTypeCreateRes 创建字典类型响应。
type DictTypeCreateRes struct {
	Id int64 `json:"id" dc:"New dict type ID"`
}

// DictTypeListReq 分页查询字典类型请求。
type DictTypeListReq struct {
	g.Meta `path:"/system/dict/type" method:"get" tags:"system" summary:"List dict types"`
	Page   int    `json:"page" d:"1" v:"min:1#page must be at least 1" dc:"Page number"`
	Size   int    `json:"size" d:"20" v:"min:1|max:200#page size must be at least 1|page size must be at most 200" dc:"Page size"`
	Name   string `json:"name" dc:"Fuzzy search by dict name"`
}

// DictTypeListRes 分页查询字典类型响应。
type DictTypeListRes struct {
	List  []*DictTypeItem `json:"list"  dc:"Dict type list"`
	Total int64           `json:"total" dc:"Total count"`
	Page  int             `json:"page"  dc:"Page number"`
	Size  int             `json:"size"  dc:"Page size"`
}

// DictTypeGetReq 查询字典类型请求。
type DictTypeGetReq struct {
	g.Meta `path:"/system/dict/type/{id}" method:"get" tags:"system" summary:"Get dict type"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Dict type ID"`
}

// DictTypeGetRes 查询字典类型响应。
type DictTypeGetRes struct {
	DictType *DictTypeItem `json:"dictType" dc:"Dict type"`
}

// DictTypeUpdateReq 更新字典类型请求。
type DictTypeUpdateReq struct {
	g.Meta `path:"/system/dict/type/{id}" method:"put" tags:"system" summary:"Update dict type"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Dict type ID"`
	Name   string `json:"name" v:"required|length:1,64#dict name is required|dict name length must be 1-64" dc:"Dict name"`
	Status int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
	Remark string `json:"remark" dc:"Remark"`
}

// DictTypeUpdateRes 更新字典类型响应。
type DictTypeUpdateRes struct {
	Updated bool `json:"updated" dc:"Whether updated"`
}

// DictTypeDeleteReq 删除字典类型请求。
type DictTypeDeleteReq struct {
	g.Meta `path:"/system/dict/type/{id}" method:"delete" tags:"system" summary:"Delete dict type"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Dict type ID"`
}

// DictTypeDeleteRes 删除字典类型响应。
type DictTypeDeleteRes struct {
	Deleted bool `json:"deleted" dc:"Whether deleted"`
}

// DictDataItem 字典数据列表项。
type DictDataItem struct {
	Id       int64  `json:"id"       dc:"Dict data ID"`
	TenantId int64  `json:"tenantId" dc:"Tenant ID"`
	DictType string `json:"dictType" dc:"Dict type code"`
	Label    string `json:"label"    dc:"Label"`
	Value    string `json:"value"    dc:"Value"`
	Sort     int    `json:"sort"     dc:"Sort order"`
	Status   int    `json:"status"   dc:"Status (1 enabled, 0 disabled)"`
	Remark   string `json:"remark"   dc:"Remark"`
}

// DictDataCreateReq 创建字典数据请求。
type DictDataCreateReq struct {
	g.Meta   `path:"/system/dict/data" method:"post" tags:"system" summary:"Create dict data"`
	DictType string `json:"dictType" v:"required|length:1,64#dict type code is required|dict type code length must be 1-64" dc:"Dict type code"`
	Label    string `json:"label" v:"required|length:1,64#label is required|label length must be 1-64" dc:"Label"`
	Value    string `json:"value" v:"required|length:1,64#value is required|value length must be 1-64" dc:"Value"`
	Sort     int    `json:"sort" d:"0" dc:"Sort order"`
	Status   int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
	Remark   string `json:"remark" dc:"Remark"`
}

// DictDataCreateRes 创建字典数据响应。
type DictDataCreateRes struct {
	Id int64 `json:"id" dc:"New dict data ID"`
}

// DictDataListReq 分页查询字典数据请求。
type DictDataListReq struct {
	g.Meta   `path:"/system/dict/data" method:"get" tags:"system" summary:"List dict data"`
	Page     int    `json:"page" d:"1" v:"min:1#page must be at least 1" dc:"Page number"`
	Size     int    `json:"size" d:"20" v:"min:1|max:200#page size must be at least 1|page size must be at most 200" dc:"Page size"`
	DictType string `json:"dictType" dc:"Filter by dict type code"`
	Label    string `json:"label" dc:"Fuzzy search by label"`
}

// DictDataListRes 分页查询字典数据响应。
type DictDataListRes struct {
	List  []*DictDataItem `json:"list"  dc:"Dict data list"`
	Total int64           `json:"total" dc:"Total count"`
	Page  int             `json:"page"  dc:"Page number"`
	Size  int             `json:"size"  dc:"Page size"`
}

// DictDataGetReq 查询字典数据请求。
type DictDataGetReq struct {
	g.Meta `path:"/system/dict/data/{id}" method:"get" tags:"system" summary:"Get dict data"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Dict data ID"`
}

// DictDataGetRes 查询字典数据响应。
type DictDataGetRes struct {
	DictData *DictDataItem `json:"dictData" dc:"Dict data"`
}

// DictDataUpdateReq 更新字典数据请求。
type DictDataUpdateReq struct {
	g.Meta   `path:"/system/dict/data/{id}" method:"put" tags:"system" summary:"Update dict data"`
	Id       int64  `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Dict data ID"`
	DictType string `json:"dictType" v:"required|length:1,64#dict type code is required|dict type code length must be 1-64" dc:"Dict type code"`
	Label    string `json:"label" v:"required|length:1,64#label is required|label length must be 1-64" dc:"Label"`
	Value    string `json:"value" v:"required|length:1,64#value is required|value length must be 1-64" dc:"Value"`
	Sort     int    `json:"sort" d:"0" dc:"Sort order"`
	Status   int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
	Remark   string `json:"remark" dc:"Remark"`
}

// DictDataUpdateRes 更新字典数据响应。
type DictDataUpdateRes struct {
	Updated bool `json:"updated" dc:"Whether updated"`
}

// DictDataDeleteReq 删除字典数据请求。
type DictDataDeleteReq struct {
	g.Meta `path:"/system/dict/data/{id}" method:"delete" tags:"system" summary:"Delete dict data"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Dict data ID"`
}

// DictDataDeleteRes 删除字典数据响应。
type DictDataDeleteRes struct {
	Deleted bool `json:"deleted" dc:"Whether deleted"`
}
