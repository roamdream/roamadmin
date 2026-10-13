// Package demo 示例模块：演示并验证"增加一个模块"的标准动作。
//
// 标准动作 = 新建目录 + 实现 Module 接口 + init() 里 Register + 在
// cmd/server/modules.go 里加一行 import。内核自动接管 Init/Start/Stop。
package demo

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/roamdream/roamadmin/kernel/module"
	"github.com/roamdream/roamadmin/modules/demo/internal/controller"
)

func init() {
	module.Register(New())
}

// Module demo 模块。
type Module struct{}

// New 创建 demo 模块实例。
func New() *Module { return &Module{} }

// Name 模块名。
func (m *Module) Name() string { return "demo" }

// Version 版本号。
func (m *Module) Version() string { return "0.1.0" }

// Init 依赖注入阶段：demo 无内部依赖，直接返回。
func (m *Module) Init(ctx context.Context) error { return nil }

// Start 启动阶段：注册 HTTP 路由。
func (m *Module) Start(ctx context.Context) error {
	m.registerRoutes(ctx)
	return nil
}

// Stop 优雅关闭阶段。
func (m *Module) Stop(ctx context.Context) error { return nil }

func (m *Module) registerRoutes(ctx context.Context) {
	s := g.Server()
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		// 控制器方法的路由由 api 中请求结构体的 g.Meta 定义。
		group.Bind(controller.New())
	})
}
