// Package model device 模块数据模型。
package model

import "github.com/gogf/gf/v2/os/gtime"

// DeviceModel 物模型/产品实体（对应表 ra_device_model，租户隔离）。
//
// Key 为产品标识（ProductKey），MQTT 上行主题段，需全局唯一。
type DeviceModel struct {
	Id          int64       `json:"id"          orm:"id,primary"`
	TenantId    int64       `json:"tenantId"    orm:"tenant_id"`
	Name        string      `json:"name"        orm:"name"`
	Key         string      `json:"key"         orm:"product_key"`
	Description string      `json:"description" orm:"description"`
	Status      int         `json:"status"      orm:"status"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}

// TableName 返回表名。
func (DeviceModel) TableName() string { return "ra_device_model" }

// DeviceModelProperty 物模型属性实体（对应表 ra_device_model_property）。
type DeviceModelProperty struct {
	Id         int64  `json:"id"         orm:"id,primary"`
	ModelId    int64  `json:"modelId"    orm:"model_id"`
	Identifier string `json:"identifier" orm:"identifier"`
	Name       string `json:"name"       orm:"name"`
	DataType   string `json:"dataType"   orm:"data_type"`
	Unit       string `json:"unit"       orm:"unit"`
	Access     string `json:"access"     orm:"access"`
	Sort       int    `json:"sort"       orm:"sort"`
}

// TableName 返回表名。
func (DeviceModelProperty) TableName() string { return "ra_device_model_property" }

// Device 设备实体（对应表 ra_device，租户隔离）。
//
// ProductKey 冗余自物模型，与 DeviceKey 组成 MQTT 上行主题段，需全局唯一。
type Device struct {
	Id         int64       `json:"id"         orm:"id,primary"`
	TenantId   int64       `json:"tenantId"   orm:"tenant_id"`
	ModelId    int64       `json:"modelId"    orm:"model_id"`
	ProductKey string      `json:"productKey" orm:"product_key"`
	DeviceKey  string      `json:"deviceKey"  orm:"device_key"`
	Name       string      `json:"name"       orm:"name"`
	Secret     string      `json:"secret"     orm:"secret"`
	Status     int         `json:"status"     orm:"status"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"`
}

// TableName 返回表名。
func (Device) TableName() string { return "ra_device" }

// DeviceShadow 设备影子实体（对应表 ra_device_shadow，device_id 唯一）。
//
// Reported / Desired 存储 JSON 文本（map[string]any）。
type DeviceShadow struct {
	Id        int64       `json:"id"        orm:"id,primary"`
	DeviceId  int64       `json:"deviceId"  orm:"device_id"`
	Reported  string      `json:"reported"  orm:"reported"`
	Desired   string      `json:"desired"   orm:"desired"`
	Version   int64       `json:"version"   orm:"version"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}

// TableName 返回表名。
func (DeviceShadow) TableName() string { return "ra_device_shadow" }

// ShadowState 影子运行时视图（内存态，reported/desired 解析为 map）。
type ShadowState struct {
	DeviceId int64          `json:"deviceId"`
	Reported map[string]any `json:"reported"`
	Desired  map[string]any `json:"desired"`
	Version  int64          `json:"version"`
}
