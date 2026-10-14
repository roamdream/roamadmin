// Package tsdb 时序存储抽象（IoT 命门，接口先行、实现后置）。
//
// P3 阶段落 TDengine / ClickHouse 具体实现；内核与业务模块只面向本接口编程。
package tsdb

import (
	"context"
	"time"
)

// Point 时序数据点。
type Point struct {
	Measurement string            `json:"measurement"`
	Tags        map[string]string `json:"tags,omitempty"`
	Fields      map[string]any    `json:"fields"`
	Time        time.Time         `json:"time"`
}

// Engine 时序存储抽象接口。
type Engine interface {
	// Write 批量写入数据点。
	Write(ctx context.Context, points ...*Point) error
	// Query 执行查询，返回实现相关的原始结果（各后端自行解析）。
	Query(ctx context.Context, statement string, args ...any) (any, error)
	// Close 关闭连接。
	Close(ctx context.Context) error
}

// RecentQuerier 按 measurement 查询最近数据点的能力（可选，具体后端可额外实现）。
//
// 设备模块据此查询某设备最近遥测；换后端时通过类型断言按需取用。
type RecentQuerier interface {
	QueryRecent(ctx context.Context, measurement string, limit int) ([]*Point, error)
}
