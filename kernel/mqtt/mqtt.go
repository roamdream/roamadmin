// Package mqtt 内嵌 MQTT Broker 封装（基于 mochi-mqtt v2）。
//
// 该能力为可选模块：默认（不加 build tag）编译进占位实现，不引入 MQTT 依赖；
// 使用 `go build -tags mqtt` 才编译真正的 Broker（引入 mochi-mqtt 依赖）。
//
// 内核只暴露最小能力：启动监听、发布下行、订阅全部消息。
// 业务模块（gateway）负责把上行消息经事件总线分发，device 负责消化。
package mqtt

// EventMessage 事件总线主题：gateway 收到上行 MQTT 消息后发布，device 订阅消化。
//
// 放在内核层作为网关与设备模块间的解耦契约，双方只依赖内核，互不 import。
const EventMessage = "mqtt.message"
