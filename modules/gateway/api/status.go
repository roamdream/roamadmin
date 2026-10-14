package api

import "github.com/gogf/gf/v2/frame/g"

// GatewayStatusReq 查询网关状态请求。
type GatewayStatusReq struct {
	g.Meta `path:"/gateway/status" method:"get" tags:"gateway" summary:"Get MQTT gateway status"`
}

// GatewayStatusRes 查询网关状态响应。
type GatewayStatusRes struct {
	Running bool   `json:"running" dc:"Whether MQTT broker is running"`
	Address string `json:"address" dc:"Broker listen address"`
}
