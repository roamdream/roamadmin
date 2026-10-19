package audit

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/roamdream/roamadmin/kernel/auth"
)

// annKey 注解的 ctx key（经 SetCtxVar 注入，供中间件在 Next 之后读取）。
type annKey struct{}

// annotation 业务侧对 module/action 的语义覆盖，以及数据变更 diff（detail）。
type annotation struct {
	module string
	action string
	detail string
}

// setAnnotation 读取并合并已有注解后写回（Annotate / AnnotateDetail 复用）。
func setAnnotation(ctx context.Context, fn func(*annotation)) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return
	}
	a := annotationOf(r)
	if a == nil {
		a = &annotation{}
	}
	fn(a)
	r.SetCtxVar(annKey{}, *a)
}

// Annotate 在业务 handler 内覆盖本条操作日志的 module/action 语义。
// 其余字段（用户、method/path/ip、耗时、结果、目标）仍由中间件自动采集。
func Annotate(ctx context.Context, module, action string) {
	setAnnotation(ctx, func(a *annotation) {
		if module != "" {
			a.module = module
		}
		if action != "" {
			a.action = action
		}
	})
}

// AnnotateDetail 附加数据变更 diff（JSON），写入 detail 字段；可与 Annotate 叠加。
func AnnotateDetail(ctx context.Context, detail string) {
	if detail == "" {
		return
	}
	setAnnotation(ctx, func(a *annotation) { a.detail = detail })
}

func annotationOf(r *ghttp.Request) *annotation {
	v := r.GetCtxVar(annKey{})
	if v == nil || v.IsNil() {
		return nil
	}
	if a, ok := v.Val().(annotation); ok {
		return &a
	}
	return nil
}

// Middleware 全局操作审计中间件：放在 auth 之后、errcode 之前。
func (m *Manager) Middleware(r *ghttp.Request) {
	if !m.enabled {
		r.Middleware.Next()
		return
	}
	path := r.URL.Path
	if !strings.HasPrefix(path, "/api/v1/") || m.skipPaths[path] {
		r.Middleware.Next()
		return
	}

	write := isWriteMethod(r.Method)
	capture := m.bodyCapture && (write || jsonContentType(r))
	var reqBody string
	if capture {
		reqBody = m.captureBody(r)
	}

	start := time.Now()
	r.Middleware.Next()

	failed := r.GetError() != nil || r.Response.Status >= 400
	if !write && !m.logGet && !failed {
		return
	}

	e := &Entry{
		Method:     r.Method,
		Path:       path,
		Query:      r.URL.RawQuery,
		IP:         r.GetClientIp(),
		UserAgent:  r.GetHeader("User-Agent"),
		Referer:    r.GetHeader("Referer"),
		TraceID:    gctx.CtxId(r.Context()),
		RespCode:   r.Response.Status,
		DurationMs: time.Since(start).Milliseconds(),
		ReqBody:    reqBody,
	}
	if failed {
		e.Status = 1
	}
	if err := r.GetError(); err != nil {
		e.ErrMsg = err.Error()
	}
	e.Module, e.Action = inferModuleAction(r.Method, path)
	if ann := annotationOf(r); ann != nil {
		if ann.module != "" {
			e.Module = ann.module
		}
		if ann.action != "" {
			e.Action = ann.action
		}
		if ann.detail != "" {
			e.Detail = ann.detail
		}
	}
	e.TargetType, e.TargetID = inferTarget(path)
	if c, ok := auth.UserFrom(r.Context()); ok {
		e.TenantID = c.TenantID
		e.UserID = c.UserID
		e.Username = c.Username
	}
	m.record(e)
}

// inferModuleAction 由路径与 HTTP 方法推断 module/action（无业务注解时的兜底）。
// 例：/api/v1/device/tag POST → module=device action=create；GET /api/v1/device/tag/1 → action=get。
func inferModuleAction(method, path string) (string, string) {
	module := pathSegment(path, 2)
	if module == "" {
		return "", ""
	}
	return module, inferAction(method, path)
}

func inferAction(method, path string) string {
	switch method {
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	case "GET":
		if trailingID(path) {
			return "get"
		}
		return "list"
	default:
		return strings.ToLower(method)
	}
}

// inferTarget 从路径推断目标资源（模块 + 数字 ID）。
// 例：/api/v1/system/user/7/password → target_type=system target_id=7。
func inferTarget(path string) (string, int64) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	module := pathSegment(path, 2)
	for i := 3; i < len(parts); i++ {
		if id, ok := parseID(parts[i]); ok {
			return module, id
		}
	}
	return "", 0
}

func pathSegment(path string, idx int) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if idx < len(parts) {
		return parts[idx]
	}
	return ""
}

func trailingID(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return false
	}
	_, ok := parseID(parts[len(parts)-1])
	return ok
}

func parseID(s string) (int64, bool) {
	if s == "" {
		return 0, false
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}

func isWriteMethod(m string) bool {
	switch m {
	case "POST", "PUT", "PATCH", "DELETE":
		return true
	}
	return false
}

func jsonContentType(r *ghttp.Request) bool {
	return strings.Contains(r.GetHeader("Content-Type"), "application/json")
}

// captureBody 读取请求体并脱敏敏感字段、按长度截断。
func (m *Manager) captureBody(r *ghttp.Request) string {
	body := r.GetBody()
	if len(body) == 0 {
		return ""
	}
	if len(body) > m.bodyMaxLen {
		body = body[:m.bodyMaxLen]
	}
	return m.maskSensitive(body)
}

// maskSensitive 对 JSON 请求体做敏感字段脱敏（key 大小写不敏感，递归）。
// 非 JSON 内容原样截断返回。
func (m *Manager) maskSensitive(body []byte) string {
	var v any
	if err := json.Unmarshal(body, &v); err != nil {
		return truncate(string(body), m.bodyMaxLen)
	}
	maskValue(v, m.sensitive)
	b, err := json.Marshal(v)
	if err != nil {
		return truncate(string(body), m.bodyMaxLen)
	}
	return truncate(string(b), m.bodyMaxLen)
}

func maskValue(v any, sensitive map[string]bool) {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			if sensitive[strings.ToLower(k)] {
				t[k] = "***"
				continue
			}
			maskValue(val, sensitive)
		}
	case []any:
		for _, item := range t {
			maskValue(item, sensitive)
		}
	}
}

// truncate 按 rune 安全截断（避免截断多字节字符产生乱码）。
func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max]) + "…(truncated)"
}
