// Package config 统一配置入口。
//
// 基于 GoFrame gcfg：多环境 yaml、命令行/环境变量覆盖、RA_ 前缀注入均由其承担。
// 内核各子包与业务模块统一通过本包（或 ra 门面）读取配置，不直接触碰 gcfg。
package config

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

// Manager 配置管理器。
type Manager struct {
	inner *gcfg.Config
}

// New 包装一个 gcfg.Config。
func New(inner *gcfg.Config) *Manager { return &Manager{inner: inner} }

// Default 返回全局默认配置实例。
func Default() *Manager { return &Manager{inner: g.Cfg()} }

// Inner 返回底层 gcfg.Config（供高级用法）。
func (m *Manager) Inner() *gcfg.Config { return m.inner }

// MustGet 读取配置项，不存在则 panic。
func (m *Manager) MustGet(ctx context.Context, pattern string) *gvar.Var {
	return m.inner.MustGet(ctx, pattern)
}

// Get 读取配置项。
func (m *Manager) Get(ctx context.Context, pattern string) (*gvar.Var, error) {
	return m.inner.Get(ctx, pattern)
}

// MustString 读取字符串配置，缺省时返回默认值。
func (m *Manager) MustString(ctx context.Context, pattern string, def ...string) string {
	v, err := m.inner.Get(ctx, pattern)
	if err != nil || v.IsNil() {
		if len(def) > 0 {
			return def[0]
		}
		return ""
	}
	return v.String()
}

// MustInt 读取整型配置，缺省时返回默认值。
func (m *Manager) MustInt(ctx context.Context, pattern string, def ...int) int {
	v, err := m.inner.Get(ctx, pattern)
	if err != nil || v.IsNil() {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	return v.Int()
}

// MustBool 读取布尔配置，缺省时返回默认值。
func (m *Manager) MustBool(ctx context.Context, pattern string, def ...bool) bool {
	v, err := m.inner.Get(ctx, pattern)
	if err != nil || v.IsNil() {
		if len(def) > 0 {
			return def[0]
		}
		return false
	}
	return v.Bool()
}
