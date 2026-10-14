//go:build !mqtt

package mqtt

import "errors"

// Broker 内嵌 MQTT Broker 占位实现（未使用 -tags mqtt 编译，不引入 MQTT 依赖）。
//
// 保持与真实现相同的导出方法签名，使 ra.MQTT() 门面在两种构建下均可编译；
// 未启用 MQTT 时所有 MQTT 操作返回错误 / 空值。
type Broker struct {
	started bool
}

// New 创建占位 Broker。
func New() *Broker { return &Broker{} }

// Available MQTT 能力是否已编译进本二进制。
func (b *Broker) Available() bool { return false }

// Start 占位：未启用 MQTT，返回错误。
func (b *Broker) Start(addr string) error {
	return errors.New("mqtt is disabled: build with -tags mqtt to enable")
}

// Publish 占位：未启用 MQTT，返回错误。
func (b *Broker) Publish(topic string, payload []byte, qos byte) error {
	return errors.New("mqtt is disabled: build with -tags mqtt to enable")
}

// OnMessage 占位：忽略。
func (b *Broker) OnMessage(h func(topic string, payload []byte)) {}

// Addr 占位：返回空。
func (b *Broker) Addr() string { return "" }

// Running 占位：始终 false。
func (b *Broker) Running() bool { return b.started }

// Close 占位：无操作。
func (b *Broker) Close() error { b.started = false; return nil }
