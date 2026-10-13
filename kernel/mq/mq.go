// Package mq 消息队列抽象（Producer / Consumer）。
//
// 内核定义统一接口，先落地本地内存实现；未来换 NSQ/Kafka/NATS 或拆微服务时
// 业务侧零改动（仅替换实现注入）。
package mq

import (
	"context"
	"sync"
)

// Message 消息体。
type Message struct {
	Topic   string            `json:"topic"`
	Key     string            `json:"key"` // 消息键（分区 / 幂等）
	Payload []byte            `json:"payload"`
	Headers map[string]string `json:"headers,omitempty"`
}

// Producer 生产者。
type Producer interface {
	Send(ctx context.Context, msg *Message) error
}

// ConsumerHandler 消费处理函数（返回 error 触发重试/告警）。
type ConsumerHandler func(ctx context.Context, msg *Message) error

// Consumer 消费者。
type Consumer interface {
	Subscribe(topic string, h ConsumerHandler) error
	Close() error
}

// Manager 消息管理器（兼具生产/消费能力），为内核门面 ra.MQ() 的统一形态。
type Manager interface {
	Producer
	Consumer
}

// Local 本地内存实现（单机默认，异步投递）。
type Local struct {
	mu     sync.RWMutex
	subs   map[string][]ConsumerHandler
	queue  chan *Message
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewLocal 创建本地消息管理器。
func NewLocal() *Local {
	ctx, cancel := context.WithCancel(context.Background())
	l := &Local{
		subs:   make(map[string][]ConsumerHandler),
		queue:  make(chan *Message, 1024),
		ctx:    ctx,
		cancel: cancel,
	}
	l.wg.Add(1)
	go l.dispatch()
	return l
}

// Send 投递消息（异步入队）。
func (l *Local) Send(ctx context.Context, msg *Message) error {
	if msg == nil {
		return nil
	}
	select {
	case l.queue <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-l.ctx.Done():
		return l.ctx.Err()
	}
}

// Subscribe 订阅主题（同一主题可有多个处理器，按订阅顺序串行派发）。
func (l *Local) Subscribe(topic string, h ConsumerHandler) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.subs[topic] = append(l.subs[topic], h)
	return nil
}

// Close 关闭管理器。
func (l *Local) Close() error {
	l.cancel()
	l.wg.Wait()
	return nil
}

func (l *Local) dispatch() {
	defer l.wg.Done()
	for {
		select {
		case msg := <-l.queue:
			l.mu.RLock()
			handlers := append([]ConsumerHandler(nil), l.subs[msg.Topic]...)
			l.mu.RUnlock()
			for _, h := range handlers {
				// 本地实现：消费失败仅记录，不影响后续消息。
				_ = h(l.ctx, msg)
			}
		case <-l.ctx.Done():
			return
		}
	}
}
