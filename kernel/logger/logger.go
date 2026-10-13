// Package logger 结构化日志封装。
//
// 基于 GoFrame glog，全局统一前缀 [ra]；通过 context 自动携带 traceId
// （glog 会从 ctx 中提取 traceId 写入每条日志），配合 observ 的埋点做链路关联。
package logger

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

// BrandPrefix 框架统一日志前缀。
const BrandPrefix = "[ra]"

// Logger 封装 glog.Logger，保留其全部方法（Debug/Info/Error 等）。
type Logger struct {
	*glog.Logger
}

// New 包装一个 glog.Logger。
func New(inner *glog.Logger) *Logger { return &Logger{Logger: inner} }

// Default 返回全局默认 logger。
func Default() *Logger { return New(g.Log()) }

// Named 返回带模块前缀的子 logger，例如 [ra][system]。
func (l *Logger) Named(name string) *Logger {
	child := l.Clone()
	child.SetPrefix(BrandPrefix + "[" + name + "]")
	return &Logger{Logger: child}
}
