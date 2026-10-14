package mq

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gredis"
)

// newTestStreams 构造指向本地 Redis 的 Streams；Redis 不可用时跳过（集成测试）。
func newTestStreams(t *testing.T) *Streams {
	t.Helper()
	r, err := gredis.New(&gredis.Config{Address: "127.0.0.1:6379", Db: 0})
	if err != nil {
		t.Skipf("redis adapter init failed: %v", err)
	}
	c := rawClient(r)
	if c == nil {
		t.Skip("redis adapter client unavailable")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := c.Ping(ctx).Err(); err != nil {
		t.Skipf("redis unreachable: %v", err)
	}
	return NewStreams(StreamsConfig{
		Redis:        r,
		Group:        fmt.Sprintf("test-grp-%d", time.Now().UnixNano()),
		Consumer:     "test-consumer",
		MaxLen:       100,
		DLQMaxLen:    100,
		BlockTimeout: time.Second,
	})
}

// TestStreamsRoundTrip 验证 XADD → XREADGROUP → 回调 → XACK 全链路。
func TestStreamsRoundTrip(t *testing.T) {
	s := newTestStreams(t)
	defer s.Close()
	topic := fmt.Sprintf("test.rt.%d", time.Now().UnixNano())
	key := streamPrefix + topic

	got := make(chan *Message, 1)
	if err := s.Subscribe(topic, func(_ context.Context, m *Message) error {
		got <- m
		return nil
	}); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := s.Send(context.Background(), &Message{
		Topic:   topic,
		Key:     "k1",
		Payload: []byte("hello"),
		Headers: map[string]string{"a": "b"},
	}); err != nil {
		t.Fatalf("send: %v", err)
	}

	select {
	case m := <-got:
		if string(m.Payload) != "hello" || m.Key != "k1" || m.Headers["a"] != "b" {
			t.Fatalf("unexpected message: topic=%q key=%q payload=%q headers=%v", m.Topic, m.Key, m.Payload, m.Headers)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for consume")
	}

	// 消费成功后应 ack，pending 归零。
	time.Sleep(300 * time.Millisecond)
	pending, err := s.client.XPending(context.Background(), key, s.group).Result()
	if err != nil {
		t.Fatalf("xpending: %v", err)
	}
	if pending.Count != 0 {
		t.Fatalf("expected 0 pending after ack, got %d", pending.Count)
	}
}

// TestStreamsDLQ 验证消费失败转死信 + ack。
func TestStreamsDLQ(t *testing.T) {
	s := newTestStreams(t)
	defer s.Close()
	topic := fmt.Sprintf("test.dlq.%d", time.Now().UnixNano())
	key := streamPrefix + topic
	dlqKey := key + streamDLQSuffix

	got := make(chan *Message, 1)
	if err := s.Subscribe(topic, func(_ context.Context, m *Message) error {
		got <- m
		return errors.New("boom")
	}); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := s.Send(context.Background(), &Message{Topic: topic, Payload: []byte("bad")}); err != nil {
		t.Fatalf("send: %v", err)
	}
	select {
	case <-got:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for consume")
	}

	// 失败消息进死信流，且原消息已 ack（pending 归零）。
	time.Sleep(300 * time.Millisecond)
	n, err := s.client.XLen(context.Background(), dlqKey).Result()
	if err != nil {
		t.Fatalf("xlen dlq: %v", err)
	}
	if n < 1 {
		t.Fatalf("expected >=1 dlq entry, got %d", n)
	}
	pending, err := s.client.XPending(context.Background(), key, s.group).Result()
	if err != nil {
		t.Fatalf("xpending: %v", err)
	}
	if pending.Count != 0 {
		t.Fatalf("expected 0 pending after dlq, got %d", pending.Count)
	}
}
