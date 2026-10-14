package api

import "github.com/gogf/gf/v2/frame/g"

// ShadowGetReq 查询设备影子请求。
type ShadowGetReq struct {
	g.Meta `path:"/device/{id}/shadow" method:"get" tags:"device" summary:"Get device shadow"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Device ID"`
}

// ShadowGetRes 查询设备影子响应。
type ShadowGetRes struct {
	DeviceId int64          `json:"deviceId" dc:"Device ID"`
	Reported map[string]any `json:"reported" dc:"Reported state"`
	Desired  map[string]any `json:"desired"  dc:"Desired state"`
	Delta    map[string]any `json:"delta"    dc:"Delta (desired - reported)"`
	Version  int64          `json:"version"  dc:"Shadow version"`
}

// ShadowUpdateDesiredReq 更新影子期望态请求。
type ShadowUpdateDesiredReq struct {
	g.Meta  `path:"/device/{id}/shadow/desired" method:"put" tags:"device" summary:"Update device shadow desired state"`
	Id      int64          `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Device ID"`
	Desired map[string]any `json:"desired" v:"required#desired state is required" dc:"Desired state to merge"`
}

// ShadowUpdateDesiredRes 更新影子期望态响应。
type ShadowUpdateDesiredRes struct {
	Version int64 `json:"version" dc:"New shadow version"`
}

// ShadowDeltaReq 查询影子增量请求。
type ShadowDeltaReq struct {
	g.Meta `path:"/device/{id}/shadow/delta" method:"get" tags:"device" summary:"Get device shadow delta"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Device ID"`
}

// ShadowDeltaRes 查询影子增量响应。
type ShadowDeltaRes struct {
	Delta map[string]any `json:"delta" dc:"Delta (desired - reported)"`
}

// TelemetryReportReq 上报遥测请求（模拟设备上报，MQTT 上行与 HTTP 同构）。
type TelemetryReportReq struct {
	g.Meta `path:"/device/{id}/telemetry" method:"post" tags:"device" summary:"Report device telemetry (simulate uplink)"`
	Id     int64         `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Device ID"`
	Fields map[string]any `json:"fields" v:"required#fields is required" dc:"Telemetry fields"`
}

// TelemetryReportRes 上报遥测响应。
type TelemetryReportRes struct {
	Ok bool `json:"ok" dc:"Whether reported"`
}

// TelemetryQueryReq 查询最近遥测请求。
type TelemetryQueryReq struct {
	g.Meta `path:"/device/{id}/telemetry" method:"get" tags:"device" summary:"Query recent device telemetry"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"Device ID"`
	Limit  int   `json:"limit" d:"100" v:"min:1|max:1000#limit must be at least 1|limit must be at most 1000" dc:"Max number of points"`
}

// TelemetryPoint 遥测数据点。
type TelemetryPoint struct {
	Measurement string            `json:"measurement" dc:"Measurement name"`
	Tags        map[string]string `json:"tags"        dc:"Tags"`
	Fields      map[string]any    `json:"fields"      dc:"Fields"`
	Time        string            `json:"time"        dc:"Timestamp (UTC RFC3339)"`
}

// TelemetryQueryRes 查询最近遥测响应。
type TelemetryQueryRes struct {
	List []*TelemetryPoint `json:"list" dc:"Telemetry points (time descending)"`
}
