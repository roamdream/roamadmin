// Package module 定义框架的模块 SPI（Service Provider Interface）与注册表。
//
// 每一个业务模块（modules/*）都实现 Module 接口，并在其 init() 中调用
// Register 完成编译期注册。内核通过注册表统一管理各模块的生命周期：
// Init（依赖注入）→ Start（注册路由/启动消费者）→ Stop（优雅关闭）。
package module

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// Module 业务模块必须实现的统一接口。
type Module interface {
	// Name 返回模块名（唯一，用于注册 / 日志 / 文档分组）。
	Name() string
	// Version 返回模块版本号。
	Version() string
	// Init 依赖注入阶段：从内核门面获取能力、初始化模块内部。
	Init(ctx context.Context) error
	// Start 启动阶段：注册 HTTP/gRPC 路由、启动消费者/定时任务。
	Start(ctx context.Context) error
	// Stop 优雅关闭阶段：排空流量、回收资源。
	Stop(ctx context.Context) error
}

var (
	mu      sync.RWMutex
	modules = make(map[string]Module)
)

// Register 注册一个模块，通常在模块包的 init() 中调用。
// 重复注册同名模块会 panic（编译期尽早暴露错误）。
func Register(m Module) {
	if m == nil {
		panic("module: Register(nil)")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, dup := modules[m.Name()]; dup {
		panic(fmt.Sprintf("module: duplicate registration of %q", m.Name()))
	}
	modules[m.Name()] = m
}

// MustGet 按名称取模块，不存在则 panic。
func MustGet(name string) Module {
	mu.RLock()
	defer mu.RUnlock()
	m, ok := modules[name]
	if !ok {
		panic(fmt.Sprintf("module: %q not registered", name))
	}
	return m
}

// Get 按名称取模块，返回是否注册。
func Get(name string) (Module, bool) {
	mu.RLock()
	defer mu.RUnlock()
	m, ok := modules[name]
	return m, ok
}

// List 返回全部已注册模块（按名称排序，保证启动顺序稳定）。
func List() []Module {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Module, 0, len(modules))
	for _, m := range modules {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name() < out[j].Name() })
	return out
}
