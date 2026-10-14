//go:build mqtt

package mqtt

import (
	"context"
	"errors"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	mqttserver "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
)

// Broker 内嵌 MQTT Broker。
type Broker struct {
	mu      sync.RWMutex
	server  *mqttserver.Server
	handler func(topic string, payload []byte)
	addr    string
	started bool
}

// New 创建 Broker（懒启动，Start 时绑定端口）。
func New() *Broker { return &Broker{} }

// Available MQTT 能力是否已编译进本二进制（仅 -tags mqtt 构建为 true）。
func (b *Broker) Available() bool { return true }

// Start 启动 broker，绑定监听地址；已启动则忽略。
func (b *Broker) Start(addr string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.started {
		return nil
	}
	opts := &mqttserver.Options{
		InlineClient: true,
		// P3 阶段开放接入（AllowAll）；设备鉴权（DeviceSecret 校验）后续接入 auth 自定义 Hook。
		Hooks: []mqttserver.HookLoadConfig{
			{Hook: &auth.AllowHook{}},
		},
		Listeners: []listeners.Config{
			{Type: "tcp", ID: "tcp1", Address: addr},
		},
	}
	s := mqttserver.New(opts)
	// 内联订阅 #，捕获所有发布（含客户端上行与自身下行）。
	s.Subscribe("#", 1, func(_ *mqttserver.Client, _ packets.Subscription, pk packets.Packet) {
		b.mu.RLock()
		h := b.handler
		b.mu.RUnlock()
		if h != nil {
			h(pk.TopicName, pk.Payload)
		}
	})
	b.server = s
	b.addr = addr
	b.started = true
	go func() {
		if err := s.Serve(); err != nil {
			g.Log().Errorf(context.Background(), "[ra][mqtt] broker serve error: %v", err)
		}
	}()
	g.Log().Infof(context.Background(), "[ra][mqtt] broker listening on %s", addr)
	return nil
}

// Publish 向 broker 发布一条消息（下行）。
func (b *Broker) Publish(topic string, payload []byte, qos byte) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.server == nil {
		return errors.New("mqtt broker not started")
	}
	return b.server.Publish(topic, payload, false, qos)
}

// OnMessage 注册全部消息回调（由 gateway 在启动时注入）。
func (b *Broker) OnMessage(h func(topic string, payload []byte)) {
	b.mu.Lock()
	b.handler = h
	b.mu.Unlock()
}

// Addr 当前监听地址。
func (b *Broker) Addr() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.addr
}

// Running broker 是否已启动。
func (b *Broker) Running() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.started
}

// Close 停止 broker。
func (b *Broker) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.server == nil {
		return nil
	}
	err := b.server.Close()
	b.server = nil
	b.started = false
	return err
}
