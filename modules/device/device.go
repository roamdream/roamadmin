// Package device 设备接入模块：物模型（模型+属性）、设备、设备影子、遥测查询与 MQTT 上行消化。
//
// 职责边界：device 负责“设备域”数据建模与遥测/影子落地；MQTT broker 由 gateway 模块
// 启动并把上行消息经消息总线（mqtt.message）转发，device 订阅后调用 service.Ingest 消化，
// 两个模块零 import 依赖，仅通过消息总线解耦。
package device

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
	"github.com/roamdream/roamadmin/modules/device/internal/consts"
	"github.com/roamdream/roamadmin/modules/device/internal/controller"
	"github.com/roamdream/roamadmin/modules/device/internal/service"
)

func init() {
	module.Register(New())
}

// Module device 模块。
type Module struct{}

// New 创建 device 模块实例。
func New() *Module { return &Module{} }

// Name 模块名。
func (m *Module) Name() string { return "device" }

// Version 版本号。
func (m *Module) Version() string { return "0.1.0" }

// Init 依赖注入阶段。
func (m *Module) Init(ctx context.Context) error { return nil }

// Start 启动阶段：订阅 MQTT 上行事件 + 注册 HTTP 路由。
func (m *Module) Start(ctx context.Context) error {
	// 幂等建表（与 manifest/sql/iot.sql 一致，首启自建，无需手动执行）。
	if err := service.EnsureSchema(ctx); err != nil {
		return err
	}

	// 订阅 gateway 转发的 MQTT 上行消息（经消息总线，payload 为 JSON 编码的 {topic, payload}）。
	_ = ra.MQ().Subscribe(mqtt.EventMessage, func(ctx context.Context, m *mq.Message) error {
		var msg struct {
			Topic   string `json:"topic"`
			Payload string `json:"payload"`
		}
		if err := json.Unmarshal(m.Payload, &msg); err != nil {
			return err
		}
		return service.Ingest(ctx, msg.Topic, []byte(msg.Payload))
	})
	m.registerRoutes(ctx)
	return nil
}

// Stop 优雅关闭阶段。
func (m *Module) Stop(ctx context.Context) error { return nil }

func (m *Module) registerRoutes(ctx context.Context) {
	s := g.Server()
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		// 物模型：登录 + device:model:list 权限。
		group.Group("/", func(g *ghttp.RouterGroup) {
			g.Middleware(auth.RequireAuth, auth.RequirePermission(consts.PermModelList))
			g.Bind(controller.NewModel())
		})
		// 设备：登录 + device:device:list 权限。
		group.Group("/", func(g *ghttp.RouterGroup) {
			g.Middleware(auth.RequireAuth, auth.RequirePermission(consts.PermDeviceList))
			g.Bind(controller.NewDevice())
		})
		// 影子 + 遥测：登录即可（设备运行态视图，非管理配置）。
		group.Group("/", func(g *ghttp.RouterGroup) {
			g.Middleware(auth.RequireAuth)
			g.Bind(controller.NewShadow())
		})
	})
}
