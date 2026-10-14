package api

import "github.com/gogf/gf/v2/frame/g"

// DeviceItem 设备列表项。
type DeviceItem struct {
	Id         int64  `json:"id"         dc:"Device ID"`
	TenantId   int64  `json:"tenantId"   dc:"Tenant ID"`
	ModelId    int64  `json:"modelId"    dc:"Model ID"`
	ProductKey string `json:"productKey" dc:"Product key"`
	DeviceKey  string `json:"deviceKey"  dc:"Device key"`
	Name       string `json:"name"       dc:"Device name"`
	Secret     string `json:"secret"     dc:"Device secret"`
	Status     int    `json:"status"     dc:"Status (1 enabled, 0 disabled)"`
	CreatedAt  string `json:"createdAt"  dc:"Creation time"`
}

// DeviceCreateReq 创建设备请求。
type DeviceCreateReq struct {
	g.Meta    `path:"/device" method:"post" tags:"device" summary:"Create device"`
	ModelId   int64  `json:"modelId" v:"required|min:1#model ID is required|model ID is invalid" dc:"Model ID"`
	DeviceKey string `json:"deviceKey" v:"required|length:1,64#device key is required|device key length must be 1-64" dc:"Device key (unique within product)"`
	Name      string `json:"name" v:"required|length:1,64#device name is required|device name length must be 1-64" dc:"Device name"`
	Secret    string `json:"secret" dc:"Device secret (auto-generated if empty)"`
	Status    int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
}

// DeviceCreateRes 创建设备响应。
type DeviceCreateRes struct {
	Id int64 `json:"id" dc:"New device ID"`
}

// DeviceListReq 分页查询设备请求。
type DeviceListReq struct {
	g.Meta    `path:"/device" method:"get" tags:"device" summary:"List devices"`
	Page      int    `json:"page" d:"1" v:"min:1#page must be at least 1" dc:"Page number"`
	Size      int    `json:"size" d:"20" v:"min:1|max:200#page size must be at least 1|page size must be at most 200" dc:"Page size"`
	ModelId   int64  `json:"modelId" dc:"Filter by model ID"`
	Name      string `json:"name" dc:"Fuzzy search by device name"`
}

// DeviceListRes 分页查询设备响应。
type DeviceListRes struct {
	List  []*DeviceItem `json:"list"  dc:"Device list"`
	Total int64         `json:"total" dc:"Total count"`
	Page  int           `json:"page"  dc:"Page number"`
	Size  int           `json:"size"  dc:"Page size"`
}

// DeviceGetReq 查询设备请求。
type DeviceGetReq struct {
	g.Meta `path:"/device/{id}" method:"get" tags:"device" summary:"Get device"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Device ID"`
}

// DeviceGetRes 查询设备响应。
type DeviceGetRes struct {
	Device *DeviceItem `json:"device" dc:"Device"`
}

// DeviceUpdateReq 更新设备请求。
type DeviceUpdateReq struct {
	g.Meta `path:"/device/{id}" method:"put" tags:"device" summary:"Update device"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Device ID"`
	Name   string `json:"name" v:"required|length:1,64#device name is required|device name length must be 1-64" dc:"Device name"`
	Status int    `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
}

// DeviceUpdateRes 更新设备响应。
type DeviceUpdateRes struct {
	Updated bool `json:"updated" dc:"Whether updated"`
}

// DeviceDeleteReq 删除设备请求。
type DeviceDeleteReq struct {
	g.Meta `path:"/device/{id}" method:"delete" tags:"device" summary:"Delete device"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Device ID"`
}

// DeviceDeleteRes 删除设备响应。
type DeviceDeleteRes struct {
	Deleted bool `json:"deleted" dc:"Whether deleted"`
}

// DeviceDownlinkReq 设备下行命令请求。
type DeviceDownlinkReq struct {
	g.Meta `path:"/device/{id}/command" method:"post" tags:"device" summary:"Send command to device"`
	Id     int64           `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Device ID"`
	Data   map[string]any  `json:"data" v:"required#command data is required" dc:"Command payload"`
}

// DeviceDownlinkRes 设备下行命令响应。
type DeviceDownlinkRes struct {
	Topic   string `json:"topic"   dc:"Published topic"`
	Payload string `json:"payload" dc:"Published payload"`
}
