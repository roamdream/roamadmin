package controller

import (
	"context"

	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/gateway/api"
)

// Gateway 网关控制器。
type Gateway struct{}

// NewGateway 创建网关控制器。
func NewGateway() *Gateway { return &Gateway{} }

// Status 查询 MQTT 网关状态。
func (c *Gateway) Status(ctx context.Context, req *api.GatewayStatusReq) (*api.GatewayStatusRes, error) {
	broker := ra.MQTT()
	return &api.GatewayStatusRes{
		Running: broker.Running(),
		Address: broker.Addr(),
	}, nil
}
