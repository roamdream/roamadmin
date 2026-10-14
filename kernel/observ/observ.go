// Package observ 可观测性：Prometheus metrics + 链路追踪埋点。
//
// 三件套（结构化日志 + metrics + 追踪）中的 metrics 与 tracing 在此；
// 结构化日志由 kernel/logger（glog）承担，三者通过 context 中的 traceId 串联。
package observ

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/roamdream/roamadmin/kernel/pool"
)

// Manager 可观测管理器。
type Manager struct {
	registry  *prometheus.Registry
	requests  *prometheus.CounterVec
	duration  *prometheus.HistogramVec
	mqSend    *prometheus.CounterVec
	mqConsume *prometheus.CounterVec
	mqDLQ     *prometheus.CounterVec
}

// New 创建可观测管理器，注册 HTTP / 协程池 / 消息总线指标与 Go 运行时指标（goroutine 数等）。
func New() *Manager {
	reg := prometheus.NewRegistry()
	m := &Manager{registry: reg}
	m.requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ra_http_requests_total",
			Help: "HTTP 请求总数（按方法/路径/状态码）。",
		},
		[]string{"method", "path", "code"},
	)
	m.duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ra_http_request_duration_seconds",
			Help:    "HTTP 请求延迟分布。",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	m.mqSend = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "ra_mq_send_total", Help: "消息总线投递总数（按主题）。"},
		[]string{"topic"},
	)
	m.mqConsume = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "ra_mq_consume_total", Help: "消息总线消费总数（按主题/结果）。"},
		[]string{"topic", "result"},
	)
	m.mqDLQ = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "ra_mq_dlq_total", Help: "消息总线死信总数（按主题）。"},
		[]string{"topic"},
	)
	// NewGoCollector 覆盖 goroutine 数、GC 等运行时指标；进程指标需平台收集器（此处省略）。
	reg.MustRegister(m.requests, m.duration, m.mqSend, m.mqConsume, m.mqDLQ, prometheus.NewGoCollector())
	// 协程池指标：用 GaugeFunc 在抓取时读取（天然懒加载，避免启动期初始化协程池）。
	reg.MustRegister(
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "ra_pool_cap", Help: "协程池容量。"}, func() float64 {
			return float64(pool.Default().Cap())
		}),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "ra_pool_running", Help: "协程池运行中的协程数。"}, func() float64 {
			return float64(pool.Default().Running())
		}),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "ra_pool_free", Help: "协程池空闲协程数。"}, func() float64 {
			return float64(pool.Default().Free())
		}),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "ra_pool_waiting", Help: "协程池等待队列长度。"}, func() float64 {
			return float64(pool.Default().Waiting())
		}),
	)
	return m
}

// defaultMgr 进程级唯一管理器：/metrics 与各内核组件（mq/pool）共用同一 registry。
var defaultMgr = New()

// Default 返回进程级唯一可观测管理器。
func Default() *Manager { return defaultMgr }

// ObserveHTTP 记录一次 HTTP 请求（QPS / 延迟 / 状态码）。
func (m *Manager) ObserveHTTP(method, path string, code int, dur time.Duration) {
	codeStr := strconv.Itoa(code)
	m.requests.WithLabelValues(method, path, codeStr).Inc()
	m.duration.WithLabelValues(method, path).Observe(dur.Seconds())
}

// ObserveMQSend 记录一次消息投递。
func (m *Manager) ObserveMQSend(topic string) { m.mqSend.WithLabelValues(topic).Inc() }

// ObserveMQConsume 记录一次消息消费结果（result: ok / error）。
func (m *Manager) ObserveMQConsume(topic, result string) { m.mqConsume.WithLabelValues(topic, result).Inc() }

// ObserveMQDLQ 记录一次死信。
func (m *Manager) ObserveMQDLQ(topic string) { m.mqDLQ.WithLabelValues(topic).Inc() }

// Handler 返回 Prometheus 抓取端点 handler（挂到 /metrics）。
func (m *Manager) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// NewTraceContext 返回带新 traceId 的 context（glog 会自动打印该 traceId）。
func NewTraceContext() context.Context { return gctx.New() }
