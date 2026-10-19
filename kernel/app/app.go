// Package app 应用生命周期管理。
//
// 负责按 module 注册表驱动各业务模块的 Init → Start → Stop，
// 并托管 HTTP Server 的启动与优雅关闭（信号由 GoFrame ghttp 内部处理）。
package app

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"

	"github.com/roamdream/roamadmin/kernel/audit"
	"github.com/roamdream/roamadmin/kernel/logger"
	"github.com/roamdream/roamadmin/kernel/module"
)

// Application 应用实例。
type Application struct {
	name    string
	version string
}

// New 创建应用实例。
func New(name, version string) *Application {
	return &Application{name: name, version: version}
}

// Name 应用名。
func (a *Application) Name() string { return a.name }

// Version 版本号。
func (a *Application) Version() string { return a.version }

// Server 返回全局 HTTP Server（模块在其上注册路由）。
func (a *Application) Server() *ghttp.Server { return g.Server() }

// Run 完整生命周期：
//
//	Init 全部模块 → Start 全部模块 → 启动操作审计 → 启动 HTTP Server（阻塞，内部处理信号与优雅关闭）
//	→ 关闭审计（flush 剩余日志）→ Stop 模块。
func (a *Application) Run(ctx context.Context) error {
	// 统一日志前缀 [ra]。
	g.Log().SetPrefix(logger.BrandPrefix)

	if err := a.init(ctx); err != nil {
		return err
	}
	if err := a.start(ctx); err != nil {
		_ = a.stop(ctx)
		return err
	}

	// 启动操作审计（异步批量落库 + 过期清理）。
	audit.Default().Start(ctx)

	g.Log().Infof(ctx, "%s v%s started, %d module(s) registered",
		a.name, a.version, len(module.List()))

	// 阻塞运行 HTTP Server；收到退出信号后完成 HTTP 优雅关闭再返回。
	a.Server().Run()

	// 关闭审计：停止后台协程并 flush 剩余日志。
	audit.Default().Close(ctx)

	g.Log().Info(ctx, "http server stopped, shutting down modules")
	return a.stop(ctx)
}

func (a *Application) init(ctx context.Context) error {
	for _, m := range module.List() {
		glog.Infof(ctx, "[ra] init module %s v%s", m.Name(), m.Version())
		if err := m.Init(ctx); err != nil {
			return fmt.Errorf("init module %q: %w", m.Name(), err)
		}
	}
	return nil
}

func (a *Application) start(ctx context.Context) error {
	for _, m := range module.List() {
		glog.Infof(ctx, "[ra] start module %s", m.Name())
		if err := m.Start(ctx); err != nil {
			return fmt.Errorf("start module %q: %w", m.Name(), err)
		}
	}
	return nil
}

// stop 逆序停止（后启动的先停），单个失败不阻断其余模块。
func (a *Application) stop(ctx context.Context) error {
	list := module.List()
	for i := len(list) - 1; i >= 0; i-- {
		m := list[i]
		if err := m.Stop(ctx); err != nil {
			glog.Errorf(ctx, "[ra] stop module %q: %v", m.Name(), err)
		}
	}
	return nil
}
