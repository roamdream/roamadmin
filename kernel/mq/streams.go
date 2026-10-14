package mq

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/roamdream/roamadmin/kernel/observ"
)

const (
	streamPrefix    = "ra:mq:"
	streamDLQSuffix = ":dlq"
)

// StreamsConfig Redis Streams 消息管理器配置。
type StreamsConfig struct {
	Redis        *gredis.Redis // 底层 gredis（用于取出 go-redis 客户端）；nil 时用 g.Redis()
	Group        string        // 消费者组名
	Consumer     string        // 消费者名（多实例需唯一）
	MaxLen       int64         // 主流 MAXLEN 近似上限
	DLQMaxLen    int64         // 死信流 MAXLEN 近似上限
	BlockTimeout time.Duration // XREADGROUP BLOCK 阻塞时长
}

// Streams 基于 Redis Streams 的持久化消息管理器：消费者组 + ack + 死信。
//
// gredis 未直接封装 X* 命令（且其 Do 会把嵌套数组拍平为字符串，无法可靠解析），
// 故从 gredis 适配器取出底层 go-redis UniversalClient 直用类型化 Streams API。
type Streams struct {
	client    goredis.UniversalClient
	group     string
	consumer  string
	maxLen    int64
	dlqMaxLen int64
	block     time.Duration

	ctx    context.Context
	cancel context.CancelFunc

	mu   sync.Mutex
	subs map[string]context.CancelFunc
}

// NewStreams 创建 Redis Streams 消息管理器。
func NewStreams(cfg StreamsConfig) *Streams {
	if cfg.Redis == nil {
		cfg.Redis = g.Redis()
	}
	if cfg.Group == "" {
		cfg.Group = "ra"
	}
	if cfg.Consumer == "" {
		cfg.Consumer = "default"
	}
	if cfg.MaxLen <= 0 {
		cfg.MaxLen = 100000
	}
	if cfg.DLQMaxLen <= 0 {
		cfg.DLQMaxLen = 100000
	}
	if cfg.BlockTimeout <= 0 {
		cfg.BlockTimeout = 5 * time.Second
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Streams{
		client:    rawClient(cfg.Redis),
		group:     cfg.Group,
		consumer:  cfg.Consumer,
		maxLen:    cfg.MaxLen,
		dlqMaxLen: cfg.DLQMaxLen,
		block:     cfg.BlockTimeout,
		ctx:       ctx,
		cancel:    cancel,
		subs:      make(map[string]context.CancelFunc),
	}
}

// rawClient 从 gredis 适配器取出底层 go-redis UniversalClient（共享同一连接池）。
func rawClient(r *gredis.Redis) goredis.UniversalClient {
	if c, ok := r.GetAdapter().Client().(goredis.UniversalClient); ok {
		return c
	}
	return nil
}

// Send 投递消息：XADD 写入流（MAXLEN 近似截断控制长度）。
func (s *Streams) Send(ctx context.Context, msg *Message) error {
	if msg == nil {
		return nil
	}
	if s.client == nil {
		return errors.New("mq streams: redis client unavailable")
	}
	headers := ""
	if len(msg.Headers) > 0 {
		if b, err := json.Marshal(msg.Headers); err == nil {
			headers = string(b)
		}
	}
	_, err := s.client.XAdd(ctx, &goredis.XAddArgs{
		Stream: streamPrefix + msg.Topic,
		MaxLen: s.maxLen,
		Approx: true,
		ID:     "*",
		Values: map[string]any{
			"payload": string(msg.Payload),
			"key":     msg.Key,
			"headers": headers,
		},
	}).Result()
	if err != nil {
		return err
	}
	observ.Default().ObserveMQSend(msg.Topic)
	return nil
}

// Subscribe 订阅主题：建消费者组（已存在则忽略）并起一个消费 goroutine。
func (s *Streams) Subscribe(topic string, h ConsumerHandler) error {
	if s.client == nil {
		return errors.New("mq streams: redis client unavailable")
	}
	key := streamPrefix + topic
	if err := s.client.XGroupCreateMkStream(context.Background(), key, s.group, "$").Err(); err != nil {
		if !strings.Contains(err.Error(), "BUSYGROUP") {
			return err
		}
	}
	s.mu.Lock()
	if _, dup := s.subs[topic]; dup {
		s.mu.Unlock()
		return nil
	}
	ctx, cancel := context.WithCancel(s.ctx)
	s.subs[topic] = cancel
	s.mu.Unlock()
	go s.consume(ctx, topic, h)
	return nil
}

// consume 消费循环：XREADGROUP 阻塞读新消息，成功 ack，失败转死信后 ack。
func (s *Streams) consume(ctx context.Context, topic string, h ConsumerHandler) {
	key := streamPrefix + topic
	for {
		if ctx.Err() != nil {
			return
		}
		streams, err := s.client.XReadGroup(ctx, &goredis.XReadGroupArgs{
			Group:    s.group,
			Consumer: s.consumer,
			Streams:  []string{key, ">"},
			Count:    100,
			Block:    s.block,
		}).Result()
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			if errors.Is(err, goredis.Nil) {
				continue // BLOCK 超时无新消息，继续阻塞等待。
			}
			g.Log().Warningf(ctx, "[ra][mq] xreadgroup %s error: %v", topic, err)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		for _, stream := range streams {
			for _, m := range stream.Messages {
				msg := &Message{
					Topic:   topic,
					Key:     stringOf(m.Values["key"]),
					Payload: []byte(stringOf(m.Values["payload"])),
					Headers: parseHeaders(stringOf(m.Values["headers"])),
				}
				if err := h(ctx, msg); err != nil {
					observ.Default().ObserveMQConsume(topic, "error")
					observ.Default().ObserveMQDLQ(topic)
					s.toDLQ(ctx, topic, m.ID, msg, err)
				} else {
					observ.Default().ObserveMQConsume(topic, "ok")
				}
				// 成功/失败均 ack：失败已转死信，不再重投。
				_ = s.client.XAck(ctx, key, s.group, m.ID).Err()
			}
		}
	}
}

// toDLQ 将消费失败的消息写入死信流。
func (s *Streams) toDLQ(ctx context.Context, topic, id string, msg *Message, consumeErr error) {
	if _, err := s.client.XAdd(ctx, &goredis.XAddArgs{
		Stream: streamPrefix + topic + streamDLQSuffix,
		MaxLen: s.dlqMaxLen,
		Approx: true,
		ID:     "*",
		Values: map[string]any{
			"payload": string(msg.Payload),
			"key":     msg.Key,
			"error":   consumeErr.Error(),
			"origin":  id,
		},
	}).Result(); err != nil {
		g.Log().Errorf(ctx, "[ra][mq] dlq xadd %s error: %v", topic, err)
	}
}

// Close 停止所有消费 goroutine 并取消基上下文。
func (s *Streams) Close() error {
	s.mu.Lock()
	for _, cancel := range s.subs {
		cancel()
	}
	s.subs = make(map[string]context.CancelFunc)
	s.mu.Unlock()
	s.cancel()
	return nil
}

func stringOf(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return ""
	}
}

func parseHeaders(s string) map[string]string {
	if s == "" {
		return nil
	}
	var h map[string]string
	if err := json.Unmarshal([]byte(s), &h); err != nil {
		return nil
	}
	return h
}
