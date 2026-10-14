package api

import "github.com/gogf/gf/v2/frame/g"

// DeviceModelItem 物模型列表项。
type DeviceModelItem struct {
	Id          int64  `json:"id"          dc:"Model ID"`
	TenantId    int64  `json:"tenantId"    dc:"Tenant ID"`
	Name        string `json:"name"        dc:"Model name"`
	Key         string `json:"key"         dc:"Product key (globally unique)"`
	Description string `json:"description" dc:"Description"`
	Status      int    `json:"status"      dc:"Status (1 enabled, 0 disabled)"`
	CreatedAt   string `json:"createdAt"   dc:"Creation time"`
}

// DeviceModelCreateReq 创建物模型请求。
type DeviceModelCreateReq struct {
	g.Meta      `path:"/device/model" method:"post" tags:"device" summary:"Create device model"`
	Name        string `json:"name" v:"required|length:1,64#model name is required|model name length must be 1-64" dc:"Model name"`
	Key         string `json:"key" v:"required|length:1,64#product key is required|product key length must be 1-64" dc:"Product key (globally unique)"`
	Description string `json:"description" dc:"Description"`
	Status      int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
}

// DeviceModelCreateRes 创建物模型响应。
type DeviceModelCreateRes struct {
	Id int64 `json:"id" dc:"New model ID"`
}

// DeviceModelListReq 分页查询物模型请求。
type DeviceModelListReq struct {
	g.Meta `path:"/device/model" method:"get" tags:"device" summary:"List device models"`
	Page   int    `json:"page" d:"1" v:"min:1#page must be at least 1" dc:"Page number"`
	Size   int    `json:"size" d:"20" v:"min:1|max:200#page size must be at least 1|page size must be at most 200" dc:"Page size"`
	Name   string `json:"name" dc:"Fuzzy search by model name"`
}

// DeviceModelListRes 分页查询物模型响应。
type DeviceModelListRes struct {
	List  []*DeviceModelItem `json:"list"  dc:"Model list"`
	Total int64              `json:"total" dc:"Total count"`
	Page  int                `json:"page"  dc:"Page number"`
	Size  int                `json:"size"  dc:"Page size"`
}

// DeviceModelGetReq 查询物模型请求。
type DeviceModelGetReq struct {
	g.Meta `path:"/device/model/{id}" method:"get" tags:"device" summary:"Get device model"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Model ID"`
}

// DeviceModelGetRes 查询物模型响应。
type DeviceModelGetRes struct {
	Model *DeviceModelItem `json:"model" dc:"Model"`
}

// DeviceModelUpdateReq 更新物模型请求。
type DeviceModelUpdateReq struct {
	g.Meta      `path:"/device/model/{id}" method:"put" tags:"device" summary:"Update device model"`
	Id          int64  `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Model ID"`
	Name        string `json:"name" v:"required|length:1,64#model name is required|model name length must be 1-64" dc:"Model name"`
	Description string `json:"description" dc:"Description"`
	Status      int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
}

// DeviceModelUpdateRes 更新物模型响应。
type DeviceModelUpdateRes struct {
	Updated bool `json:"updated" dc:"Whether updated"`
}

// DeviceModelDeleteReq 删除物模型请求。
type DeviceModelDeleteReq struct {
	g.Meta `path:"/device/model/{id}" method:"delete" tags:"device" summary:"Delete device model"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Model ID"`
}

// DeviceModelDeleteRes 删除物模型响应。
type DeviceModelDeleteRes struct {
	Deleted bool `json:"deleted" dc:"Whether deleted"`
}

// PropertyItem 物模型属性列表项。
type PropertyItem struct {
	Id         int64  `json:"id"         dc:"Property ID"`
	ModelId    int64  `json:"modelId"    dc:"Model ID"`
	Identifier string `json:"identifier" dc:"Property identifier"`
	Name       string `json:"name"       dc:"Property name"`
	DataType   string `json:"dataType"   dc:"Data type (float/int/string/bool)"`
	Unit       string `json:"unit"       dc:"Unit"`
	Access     string `json:"access"     dc:"Access mode (r read-only, rw read-write)"`
	Sort       int    `json:"sort"       dc:"Sort order"`
}

// PropertyCreateReq 创建属性请求。
type PropertyCreateReq struct {
	g.Meta     `path:"/device/model/{modelId}/property" method:"post" tags:"device" summary:"Create model property"`
	ModelId    int64  `json:"modelId" in:"path" v:"required|min:1#model ID is required|model ID is invalid" dc:"Model ID"`
	Identifier string `json:"identifier" v:"required|length:1,64#identifier is required|identifier length must be 1-64" dc:"Property identifier"`
	Name       string `json:"name" v:"required|length:1,64#property name is required|property name length must be 1-64" dc:"Property name"`
	DataType   string `json:"dataType" v:"required|in:float,int,string,bool#data type is required|data type must be one of float/int/string/bool" dc:"Data type"`
	Unit       string `json:"unit" dc:"Unit"`
	Access     string `json:"access" d:"r" v:"in:r,rw#access must be r or rw" dc:"Access mode"`
	Sort       int    `json:"sort" d:"0" dc:"Sort order"`
}

// PropertyCreateRes 创建属性响应。
type PropertyCreateRes struct {
	Id int64 `json:"id" dc:"New property ID"`
}

// PropertyListReq 查询属性列表请求。
type PropertyListReq struct {
	g.Meta  `path:"/device/model/{modelId}/property" method:"get" tags:"device" summary:"List model properties"`
	ModelId int64 `json:"modelId" in:"path" v:"required|min:1#model ID is required|model ID is invalid" dc:"Model ID"`
}

// PropertyListRes 查询属性列表响应。
type PropertyListRes struct {
	List []*PropertyItem `json:"list" dc:"Property list"`
}

// PropertyUpdateReq 更新属性请求。
type PropertyUpdateReq struct {
	g.Meta     `path:"/device/model/{modelId}/property/{id}" method:"put" tags:"device" summary:"Update model property"`
	ModelId    int64  `json:"modelId" in:"path" v:"required|min:1#model ID is required|model ID is invalid" dc:"Model ID"`
	Id         int64  `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Property ID"`
	Name       string `json:"name" v:"required|length:1,64#property name is required|property name length must be 1-64" dc:"Property name"`
	DataType   string `json:"dataType" v:"required|in:float,int,string,bool#data type is required|data type must be one of float/int/string/bool" dc:"Data type"`
	Unit       string `json:"unit" dc:"Unit"`
	Access     string `json:"access" d:"r" v:"in:r,rw#access must be r or rw" dc:"Access mode"`
	Sort       int    `json:"sort" d:"0" dc:"Sort order"`
}

// PropertyUpdateRes 更新属性响应。
type PropertyUpdateRes struct {
	Updated bool `json:"updated" dc:"Whether updated"`
}

// PropertyDeleteReq 删除属性请求。
type PropertyDeleteReq struct {
	g.Meta  `path:"/device/model/{modelId}/property/{id}" method:"delete" tags:"device" summary:"Delete model property"`
	ModelId int64 `json:"modelId" in:"path" v:"required|min:1#model ID is required|model ID is invalid" dc:"Model ID"`
	Id      int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Property ID"`
}

// PropertyDeleteRes 删除属性响应。
type PropertyDeleteRes struct {
	Deleted bool `json:"deleted" dc:"Whether deleted"`
}
