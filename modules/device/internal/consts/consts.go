// Package consts device 模块内部常量：表名、权限码、MQTT 主题约定、缓存键。
package consts

// 表名（统一 ra_ 前缀）。
const (
	TableDeviceModel         = "ra_device_model"
	TableDeviceModelProperty = "ra_device_model_property"
	TableDevice              = "ra_device"
	TableDeviceShadow        = "ra_device_shadow"
)

// 权限码（RBAC）。
const (
	PermModelList   = "device:model:list"
	PermModelCreate = "device:model:create"
	PermModelUpdate = "device:model:update"
	PermModelDelete = "device:model:delete"
	PermDeviceList  = "device:device:list"
	PermDeviceCreate = "device:device:create"
	PermDeviceUpdate = "device:device:update"
	PermDeviceDelete = "device:device:delete"
)

// 缓存键前缀（分布式锁）。
const (
	CachePrefixModelKey = "ra:device:model:key:"
	CachePrefixDeviceKey = "ra:device:device:key:"
)

// MQTT 主题约定（上行 / 下行）。
//
//	上行：ra/{productKey}/{deviceKey}/telemetry         —— 遥测数据
//	上行：ra/{productKey}/{deviceKey}/shadow/reported   —— 影子上报
//	下行：ra/{productKey}/{deviceKey}/command           —— 下发命令
//	下行：ra/{productKey}/{deviceKey}/shadow/desired    —— 影子期望态
const (
	TopicRoot          = "ra"
	TopicTelemetry     = "telemetry"
	TopicShadow        = "shadow"
	TopicShadowReported = "reported"
	TopicShadowDesired = "desired"
	TopicCommand       = "command"
)

// 时序 measurement 前缀：遥测写入 device:{deviceId}。
const MeasurementDevicePrefix = "device:"
