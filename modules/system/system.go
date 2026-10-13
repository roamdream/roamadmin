// Package system 系统管理模块：登录/JWT、用户/角色/权限(RBAC)、菜单、字典、多租户、操作日志。
//
// 标准动作与 demo 模块一致：新建目录 + 实现 Module 接口 + init() 里 Register + 在
// cmd/server/modules.go 里加一行 import。内核自动接管 Init/Start/Stop。
package system

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/module"
	"github.com/roamdream/roamadmin/modules/system/internal/consts"
	"github.com/roamdream/roamadmin/modules/system/internal/controller"
)

func init() {
	module.Register(New())
}

// Module system 模块。
type Module struct{}

// New 创建 system 模块实例。
func New() *Module { return &Module{} }

// Name 模块名。
func (m *Module) Name() string { return "system" }

// Version 版本号。
func (m *Module) Version() string { return "0.1.0" }

// Init 依赖注入阶段：system 无内部依赖，直接返回。
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
		// 公开接口：登录（无需鉴权，控制器方法路由由 api 请求结构体的 g.Meta 定义）。
		group.Bind(controller.NewAuth())

		// 受保护接口：登录即可（当前用户信息 / 登出）。
		group.Group("/", func(g *ghttp.RouterGroup) {
			g.Middleware(auth.RequireAuth)
			g.Bind(controller.NewProfile())
		})

		// 管理类接口：登录 + 对应权限码（RBAC，权限码见 internal/consts）。
		group.Group("/", func(g *ghttp.RouterGroup) {
			g.Middleware(auth.RequireAuth, auth.RequirePermission(consts.PermUserList))
			g.Bind(controller.NewUser())
		})
		group.Group("/", func(g *ghttp.RouterGroup) {
			g.Middleware(auth.RequireAuth, auth.RequirePermission(consts.PermRoleList))
			g.Bind(controller.NewRole())
		})
		group.Group("/", func(g *ghttp.RouterGroup) {
			g.Middleware(auth.RequireAuth, auth.RequirePermission(consts.PermMenuList))
			g.Bind(controller.NewMenu())
		})
		group.Group("/", func(g *ghttp.RouterGroup) {
			g.Middleware(auth.RequireAuth, auth.RequirePermission(consts.PermDictList))
			g.Bind(controller.NewDict())
		})
		group.Group("/", func(g *ghttp.RouterGroup) {
			g.Middleware(auth.RequireAuth, auth.RequirePermission(consts.PermTenantList))
			g.Bind(controller.NewTenant())
		})
		group.Group("/", func(g *ghttp.RouterGroup) {
			g.Middleware(auth.RequireAuth, auth.RequirePermission(consts.PermOperLogList))
			g.Bind(controller.NewOperLog())
		})
	})
}
