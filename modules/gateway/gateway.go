// Package gateway MQTT 接入网关模块：启动内嵌 Broker，把上行消息经消息总线转发给 device。
//
// 职责边界：gateway 只负责“接进来”（启动 broker + 转发），不关心业务含义；
// 下行（ra.MQTT().Publish）由 device 模块在需要时直接调用内核门面。
package gateway

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/module"
	"github.com/roamdream/roamadmin/kernel/mq"
	"github.com/roamdream/roamadmin/kernel/mqtt"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/gateway/internal/controller"
)

func init() {
	module.Register(New())
}

// Module gateway 模块。
type Module struct{}

// New 创建 gateway 模块实例。
func New() *Module { return &Module{} }

// Name 模块名。
func (m *Module) Name() string { return "gateway" }

// Version 版本号。
func (m *Module) Version() string { return "0.1.0" }

// Init 依赖注入阶段。
func (m *Module) Init(ctx context.Context) error { return nil }

// Start 启动阶段：启动内嵌 MQTT Broker + 注册状态路由。
func (m *Module) Start(ctx context.Context) error {
	m.registerRoutes(ctx)

	broker := ra.MQTT()
	enabled := ra.Cfg().MustBool(ctx, "mqtt.enabled", true)
	if !broker.Available() || !enabled {
		if !broker.Available() {
			g.Log().Info(ctx, "[gateway] mqtt disabled at build time (build with -tags mqtt to enable), skip broker")
		}
		return nil
	}
	addr := ra.Cfg().MustString(ctx, "mqtt.address", ":1883")
	if err := broker.Start(addr); err != nil {
		return err
	}
	// 上行消息（含 telemetry / shadow/reported）经消息总线转发，device 订阅消化。
	// local 模式等价进程内投递；redis-streams 模式落库并跨实例消费。
	broker.OnMessage(func(topic string, payload []byte) {
		b, _ := json.Marshal(map[string]string{"topic": topic, "payload": string(payload)})
		_ = ra.MQ().Send(ctx, &mq.Message{Topic: mqtt.EventMessage, Payload: b})
	})
	return nil
}

// Stop 优雅关闭阶段：关闭 Broker。
func (m *Module) Stop(ctx context.Context) error { return ra.MQTT().Close() }

func (m *Module) registerRoutes(ctx context.Context) {
	s := g.Server()
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Group("/", func(g *ghttp.RouterGroup) {
			g.Middleware(auth.RequireAuth)
			g.Bind(controller.NewGateway())
		})
	})
}
