// Package mq 消息队列抽象（Producer / Consumer）。
//
// 内核定义统一接口，提供两类实现：
//   - Local：本地内存（单机默认，异步投递）；
//   - Streams：Redis Streams 持久化（消息落库、跨实例消费、ack + 死信）。
//
// 通过配置 mq.type 切换（local | redis-streams），业务侧只面向 Manager 接口编程零改动。
package mq

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/roamdream/roamadmin/kernel/observ"
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

// NewFromConfig 按配置构造消息管理器：mq.type 为 redis-streams 时用持久化实现，否则本地内存实现。
func NewFromConfig(ctx context.Context) Manager {
	typ := g.Cfg().MustGet(ctx, "mq.type", "local").String()
	if typ == "redis-streams" {
		consumer := g.Cfg().MustGet(ctx, "mq.streams.consumer", "").String()
		if consumer == "" {
			host, _ := os.Hostname()
			consumer = host + "-" + strconv.Itoa(os.Getpid())
		}
		return NewStreams(StreamsConfig{
			Redis:        g.Redis(),
			Group:        g.Cfg().MustGet(ctx, "mq.streams.group", "ra").String(),
			Consumer:     consumer,
			MaxLen:       int64(g.Cfg().MustGet(ctx, "mq.streams.maxLen", 100000).Int()),
			DLQMaxLen:    int64(g.Cfg().MustGet(ctx, "mq.streams.dlqMaxLen", 100000).Int()),
			BlockTimeout: g.Cfg().MustGet(ctx, "mq.streams.blockTimeout", "5s").Duration(),
		})
	}
	return NewLocal()
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
		observ.Default().ObserveMQSend(msg.Topic)
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
				if err := h(l.ctx, msg); err != nil {
					observ.Default().ObserveMQConsume(msg.Topic, "error")
				} else {
					observ.Default().ObserveMQConsume(msg.Topic, "ok")
				}
			}
		case <-l.ctx.Done():
			return
		}
	}
}
