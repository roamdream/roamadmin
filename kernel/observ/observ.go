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
)

// Manager 可观测管理器。
type Manager struct {
	registry *prometheus.Registry
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

// New 创建可观测管理器，注册 HTTP 指标与 Go 运行时指标（goroutine 数等）。
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
	// NewGoCollector 覆盖 goroutine 数、GC 等运行时指标；进程指标需平台收集器（此处省略）。
	reg.MustRegister(m.requests, m.duration, prometheus.NewGoCollector())
	return m
}

// ObserveHTTP 记录一次 HTTP 请求（QPS / 延迟 / 状态码）。
func (m *Manager) ObserveHTTP(method, path string, code int, dur time.Duration) {
	codeStr := strconv.Itoa(code)
	m.requests.WithLabelValues(method, path, codeStr).Inc()
	m.duration.WithLabelValues(method, path).Observe(dur.Seconds())
}

// Handler 返回 Prometheus 抓取端点 handler（挂到 /metrics）。
func (m *Manager) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// NewTraceContext 返回带新 traceId 的 context（glog 会自动打印该 traceId）。
func NewTraceContext() context.Context { return gctx.New() }
