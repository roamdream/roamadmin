// Package pool 协程池与对象池（高并发资源治理，防 OOM）。
//
// 基于 ants 封装：设备消息处理、批量任务、消费者统一走池；默认池自带 panic 隔离。
package pool

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/panjf2000/ants/v2"
)

// Pool 协程池。
type Pool struct {
	inner *ants.Pool
}

// New 创建协程池，opts 透传给 ants（如 ants.WithNonblocking）。
func New(size int, opts ...ants.Option) (*Pool, error) {
	p, err := ants.NewPool(size, opts...)
	if err != nil {
		return nil, err
	}
	return &Pool{inner: p}, nil
}

// Submit 提交任务（池满时行为取决于是否 WithNonblocking）。
func (p *Pool) Submit(task func()) error { return p.inner.Submit(task) }

// Running 正在运行的协程数。
func (p *Pool) Running() int { return p.inner.Running() }

// Free 空闲协程数。
func (p *Pool) Free() int { return p.inner.Free() }

// Cap 池容量。
func (p *Pool) Cap() int { return p.inner.Cap() }

// Waiting 等待队列长度。
func (p *Pool) Waiting() int { return p.inner.Waiting() }

// Release 释放池资源。
func (p *Pool) Release() { p.inner.Release() }

// defaultPool 全局默认池：容量 10000，panic 自动 recover，避免单任务拖垮进程。
var defaultPool = newDefault()

func newDefault() *Pool {
	p, err := New(10000, ants.WithPanicHandler(func(v any) {
		g.Log().Errorf(context.Background(), "[ra] goroutine panic recovered: %v", v)
	}))
	if err != nil {
		panic(err)
	}
	return p
}

// Default 返回全局默认协程池。
func Default() *Pool { return defaultPool }

// Submit 提交任务到默认池。
func Submit(task func()) error { return defaultPool.Submit(task) }
