package observ

import (
	"net/http"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gdebug"

	"github.com/roamdream/roamadmin/kernel/errcode"
)

// Middleware 全局 HTTP 中间件：panic 隔离 + 请求日志（含 traceId）+ metrics 埋点。
func (m *Manager) Middleware(r *ghttp.Request) {
	start := time.Now()
	defer func() {
		if rec := recover(); rec != nil {
			g.Log().Errorf(r.Context(), "[ra] panic recovered: %v\n%s", rec, gdebug.Stack())
			r.Response.WriteStatus(http.StatusInternalServerError)
			r.Response.WriteJson(errcode.Fail(errcode.CodeInternal, ""))
		}
		// 回写 traceId 头，方便排查。
		if tid := gctx.CtxId(r.Context()); tid != "" {
			r.Response.Header().Set("X-Trace-Id", tid)
		}
		m.ObserveHTTP(r.Method, r.URL.Path, r.Response.Status, time.Since(start))
	}()
	r.Middleware.Next()
}
