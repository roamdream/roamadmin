// Package event 进程内事件总线（模块间解耦），含背压策略。
//
// 设备数据到达、用户登录等事件在此发布/订阅；处理器默认走内核协程池异步执行，
// 单个处理器 panic 被隔离，不拖垮总线。
package event

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/roamdream/roamadmin/kernel/pool"
)

// Handler 事件处理器。
type Handler func(ctx context.Context, topic string, payload any)

// Policy 背压策略：总线队列满时的行为。
type Policy int

const (
	// PolicyDrop 满则丢弃（默认，适合高频遥测，宁可丢点不可拖慢主流程）。
	PolicyDrop Policy = iota
	// PolicyBlock 满则阻塞（适合关键事件，保证不丢）。
	PolicyBlock
)

type message struct {
	topic   string
	payload any
}

type subscription struct {
	id uint64
	h  Handler
}

// Bus 事件总线。
type Bus struct {
	mu       sync.RWMutex
	subs     map[string][]*subscription
	queue    chan *message
	seq      atomic.Uint64
	policy   Policy
	pool     *pool.Pool
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	closed   atomic.Bool
}

// Option 总线配置项。
type Option func(*Bus)

// WithPolicy 设置背压策略。
func WithPolicy(p Policy) Option { return func(b *Bus) { b.policy = p } }

// WithChanSize 设置总线缓冲队列大小。
func WithChanSize(n int) Option { return func(b *Bus) { b.queue = make(chan *message, n) } }

// WithPool 指定处理器执行协程池（默认用内核默认池）。
func WithPool(p *pool.Pool) Option { return func(b *Bus) { b.pool = p } }

// New 创建事件总线。
func New(opts ...Option) *Bus {
	ctx, cancel := context.WithCancel(context.Background())
	b := &Bus{
		subs:   make(map[string][]*subscription),
		queue:  make(chan *message, 1024),
		policy: PolicyDrop,
		pool:   pool.Default(),
		ctx:    ctx,
		cancel: cancel,
	}
	for _, o := range opts {
		o(b)
	}
	b.wg.Add(1)
	go b.dispatch()
	return b
}

// Subscribe 订阅主题，返回取消订阅函数。
func (b *Bus) Subscribe(topic string, h Handler) (unsubscribe func()) {
	if h == nil {
		return func() {}
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	sub := &subscription{id: b.seq.Add(1), h: h}
	b.subs[topic] = append(b.subs[topic], sub)
	return func() { b.unsubscribe(topic, sub.id) }
}

// Publish 发布事件。
func (b *Bus) Publish(ctx context.Context, topic string, payload any) {
	if b.closed.Load() {
		return
	}
	msg := &message{topic: topic, payload: payload}
	switch b.policy {
	case PolicyDrop:
		select {
		case b.queue <- msg:
		default:
			g.Log().Warningf(ctx, "[ra][event] drop message on topic %q (queue full)", topic)
		}
	case PolicyBlock:
		select {
		case b.queue <- msg:
		case <-ctx.Done():
		case <-b.ctx.Done():
		}
	}
}

// Close 关闭总线，停止派发。
func (b *Bus) Close() {
	if !b.closed.CompareAndSwap(false, true) {
		return
	}
	b.cancel()
	b.wg.Wait()
}

func (b *Bus) unsubscribe(topic string, id uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	list := b.subs[topic]
	for i, s := range list {
		if s.id == id {
			b.subs[topic] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

func (b *Bus) dispatch() {
	defer b.wg.Done()
	for {
		select {
		case msg := <-b.queue:
			b.mu.RLock()
			src := b.subs[msg.topic]
			handlers := make([]Handler, 0, len(src))
			for _, s := range src {
				handlers = append(handlers, s.h)
			}
			b.mu.RUnlock()
			for _, h := range handlers {
				b.run(h, msg)
			}
		case <-b.ctx.Done():
			return
		}
	}
}

func (b *Bus) run(h Handler, msg *message) {
	task := func() {
		defer func() {
			if r := recover(); r != nil {
				g.Log().Errorf(context.Background(), "[ra][event] handler panic on topic %q: %v", msg.topic, r)
			}
		}()
		h(b.ctx, msg.topic, msg.payload)
	}
	if b.pool != nil {
		if err := b.pool.Submit(task); err != nil {
			g.Log().Errorf(context.Background(), "[ra][event] submit handler failed: %v", err)
		}
		return
	}
	task()
}
