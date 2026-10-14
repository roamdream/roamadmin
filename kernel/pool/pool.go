// Package pool 协程池与对象池（高并发资源治理，防 OOM）。
//
// 基于 ants 封装：设备消息处理、批量任务、消费者统一走池；默认池自带 panic 隔离。
package pool

import (
	"context"
	"sync"

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

var (
	defaultOnce sync.Once
	defaultPool *Pool
)

// Default 返回全局默认协程池（懒初始化：容量读配置 pool.size，缺省 10000）。
func Default() *Pool {
	defaultOnce.Do(func() {
		size := g.Cfg().MustGet(context.Background(), "pool.size", 10000).Int()
		defaultPool = newDefault(size)
	})
	return defaultPool
}

// newDefault 创建默认池：panic 自动 recover，避免单任务拖垮进程；size<=0 时回退 10000。
func newDefault(size int) *Pool {
	if size <= 0 {
		size = 10000
	}
	p, err := New(size, ants.WithPanicHandler(func(v any) {
		g.Log().Errorf(context.Background(), "[ra] goroutine panic recovered: %v", v)
	}))
	if err != nil {
		panic(err)
	}
	return p
}

// Submit 提交任务到默认池。
func Submit(task func()) error { return Default().Submit(task) }
