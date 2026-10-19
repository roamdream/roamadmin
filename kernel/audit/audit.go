// Package audit 操作审计：全局 HTTP 中间件采集 + 异步批量落库。
//
// 目标：把「谁在何时对系统做了什么、结果如何、耗时多少、请求内容」作为一条可追踪的轨迹
// 完整落库（表 ra_oper_log，租户隔离），供安全审计与问题排查使用。
//
// 所有 /api/v1/* 请求自动记录（默认只记写操作 + 全部失败请求，读请求可经 audit.logGet 开启），
// 业务侧仅在需要特殊语义时调用 Annotate 覆盖 module/action（如「重置密码」），
// 或在登录/登出等场景调用 Record 手动记录，无需在每个接口逐一手动埋点。
package audit

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/db"
)

// tableOperLog 操作日志表。
const tableOperLog = "ra_oper_log"

// defaultSensitiveKeys 请求体默认脱敏字段（大小写不敏感）。
var defaultSensitiveKeys = []string{
	"password", "pass", "token", "secret", "accesskey", "accesssecret",
	"accesskeysecret", "clientsecret", "oldpassword", "newpassword",
}

// Entry 一条操作日志。
type Entry struct {
	TraceID    string
	TenantID   int64
	UserID     int64
	Username   string
	Module     string
	Action     string
	Method     string
	Path       string
	Query      string
	IP         string
	UserAgent  string
	Referer    string
	Status     int    // 0 成功，1 失败
	RespCode   int    // HTTP 状态码
	ErrMsg     string
	DurationMs int64
	TargetType string
	TargetID   int64
	ReqBody    string
	Detail     string // 数据变更 diff（可选，JSON）
	CreatedAt  time.Time
}

// Manager 审计管理器（进程级单例）。
type Manager struct {
	enabled     bool
	logGet      bool
	bodyCapture bool
	bodyMaxLen  int
	sensitive   map[string]bool
	skipPaths   map[string]bool
	queueSize   int
	batchSize   int
	flushEvery  time.Duration
	retention   int

	ch     chan *Entry
	stop   chan struct{}
	once   sync.Once
	wg     sync.WaitGroup
	mu     sync.Mutex
	closed bool
}

var (
	defaultOnce sync.Once
	defaultMgr  *Manager
)

// Default 返回进程级唯一审计管理器（懒加载配置）。
func Default() *Manager {
	defaultOnce.Do(func() { defaultMgr = newManager() })
	return defaultMgr
}

func newManager() *Manager {
	ctx := context.Background()
	cfg := g.Cfg()
	// 敏感字段 = 内置默认 + 配置追加（配置只能增不能减，保证密码/密钥类字段始终脱敏）。
	sensitive := make(map[string]bool, len(defaultSensitiveKeys)+8)
	for _, k := range defaultSensitiveKeys {
		sensitive[strings.ToLower(k)] = true
	}
	for _, k := range cfg.MustGet(ctx, "audit.sensitiveKeys").Strings() {
		sensitive[strings.ToLower(k)] = true
	}
	flushEvery := cfg.MustGet(ctx, "audit.flushInterval", "500ms").Duration()
	if flushEvery <= 0 {
		flushEvery = 500 * time.Millisecond
	}
	return &Manager{
		enabled:     cfg.MustGet(ctx, "audit.enabled", true).Bool(),
		logGet:      cfg.MustGet(ctx, "audit.logGet", false).Bool(),
		bodyCapture: cfg.MustGet(ctx, "audit.bodyCapture", true).Bool(),
		bodyMaxLen:  cfg.MustGet(ctx, "audit.bodyMaxLen", 2048).Int(),
		sensitive:   sensitive,
		skipPaths: map[string]bool{
			// 登录/登出由 handler 显式调用 Record 记录（登录无 JWT，需从请求体取用户名）。
			"/api/v1/system/login":  true,
			"/api/v1/system/logout": true,
		},
		queueSize:  cfg.MustGet(ctx, "audit.queueSize", 4096).Int(),
		batchSize:  cfg.MustGet(ctx, "audit.batchSize", 100).Int(),
		flushEvery: flushEvery,
		retention:  cfg.MustGet(ctx, "audit.retentionDays", 90).Int(),
	}
}

// Start 启动后台写入与过期清理协程（幂等；enabled=false 时不启动）。
func (m *Manager) Start(ctx context.Context) {
	if !m.enabled {
		return
	}
	m.once.Do(func() {
		if err := m.ensureSchema(ctx); err != nil {
			g.Log().Errorf(ctx, "[ra][audit] ensure schema failed: %v", err)
		}
		m.ch = make(chan *Entry, m.queueSize)
		m.stop = make(chan struct{})
		m.wg.Add(2)
		go m.writer()
		go m.cleanupLoop()
	})
}

// Close 优雅关闭：停止后台协程并 flush 队列中剩余日志。
func (m *Manager) Close(ctx context.Context) {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return
	}
	m.closed = true
	m.mu.Unlock()
	if m.ch == nil {
		return
	}
	close(m.stop)
	m.wg.Wait()
}

// Record 手动记录一条操作日志（登录/登出/数据变更等中间件无法自动捕获语义的场景）。
// 会自动补齐 traceId / method / path / ip / UA 与用户上下文（仅当 Entry 未显式给出时）。
func Record(ctx context.Context, e *Entry) {
	if e == nil {
		return
	}
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		if e.Method == "" {
			e.Method = r.Method
		}
		if e.Path == "" {
			e.Path = r.URL.Path
		}
		if e.Query == "" {
			e.Query = r.URL.RawQuery
		}
		if e.IP == "" {
			e.IP = r.GetClientIp()
		}
		if e.UserAgent == "" {
			e.UserAgent = r.GetHeader("User-Agent")
		}
		if e.Referer == "" {
			e.Referer = r.GetHeader("Referer")
		}
		if e.TraceID == "" {
			e.TraceID = gctx.CtxId(r.Context())
		}
	}
	if c, ok := auth.UserFrom(ctx); ok {
		if e.UserID == 0 {
			e.UserID = c.UserID
		}
		if e.TenantID == 0 {
			e.TenantID = c.TenantID
		}
		if e.Username == "" {
			e.Username = c.Username
		}
	}
	Default().record(e)
}

// record 入队一条日志（异步，非阻塞；队列满时丢弃并告警）。
func (m *Manager) record(e *Entry) {
	if !m.enabled || m.ch == nil {
		return
	}
	m.mu.Lock()
	closed := m.closed
	m.mu.Unlock()
	if closed {
		return
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now().UTC()
	}
	select {
	case m.ch <- e:
	default:
		g.Log().Warningf(context.Background(), "[ra][audit] queue full, drop entry: %s %s", e.Method, e.Path)
	}
}

// writer 后台批量写入协程。
func (m *Manager) writer() {
	defer m.wg.Done()
	buf := make([]*Entry, 0, m.batchSize)
	ticker := time.NewTicker(m.flushEvery)
	defer ticker.Stop()
	flush := func() {
		if len(buf) == 0 {
			return
		}
		m.flush(buf)
		buf = buf[:0]
	}
	for {
		select {
		case e := <-m.ch:
			buf = append(buf, e)
			if len(buf) >= m.batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-m.stop:
			// 停止信号：先尽力消费队列中剩余条目再退出，避免丢日志。
			for {
				select {
				case e := <-m.ch:
					buf = append(buf, e)
					if len(buf) >= m.batchSize {
						flush()
					}
				default:
					flush()
					return
				}
			}
		}
	}
}

// flush 批量落库（多行 INSERT）。
func (m *Manager) flush(entries []*Entry) {
	ctx := context.Background()
	rows := make([]g.Map, 0, len(entries))
	for _, e := range entries {
		if e == nil {
			continue
		}
		rows = append(rows, g.Map{
			"id":          db.Default().NextID(),
			"tenant_id":   e.TenantID,
			"user_id":     e.UserID,
			"username":    e.Username,
			"module":      e.Module,
			"action":      e.Action,
			"method":      e.Method,
			"path":        e.Path,
			"query":       e.Query,
			"ip":          e.IP,
			"user_agent":  e.UserAgent,
			"referer":     e.Referer,
			"status":      e.Status,
			"resp_code":   e.RespCode,
			"err_msg":     e.ErrMsg,
			"duration_ms": e.DurationMs,
			"target_type": e.TargetType,
			"target_id":   e.TargetID,
			"trace_id":    e.TraceID,
			"req_body":    e.ReqBody,
			"detail":      e.Detail,
			"created_at":  gtime.New(e.CreatedAt),
		})
	}
	if len(rows) == 0 {
		return
	}
	if _, err := db.Default().Model(tableOperLog).Data(rows).Insert(); err != nil {
		g.Log().Errorf(ctx, "[ra][audit] batch insert %d entries failed: %v", len(rows), err)
	}
}

// cleanupLoop 定期清理过期日志。
func (m *Manager) cleanupLoop() {
	defer m.wg.Done()
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.cleanup()
		case <-m.stop:
			return
		}
	}
}

func (m *Manager) cleanup() {
	if m.retention <= 0 {
		return
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -m.retention)
	ctx := context.Background()
	if _, err := db.Default().Model(tableOperLog).Where("created_at < ?", cutoff).Delete(); err != nil {
		g.Log().Errorf(ctx, "[ra][audit] cleanup old logs failed: %v", err)
	}
}
