// 主程序入口：装配全局中间件 → 启动应用（驱动模块生命周期与 HTTP Server）。
package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/roamdream/roamadmin/kernel/app"
	"github.com/roamdream/roamadmin/kernel/errcode"
	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/kernel/swaggerui"
)

const version = "0.1.0"

func main() {
	ctx := gctx.New()

	s := g.Server()

	// OpenAPI 文档元信息（英文）。swagger/接口描述统一英文，前端再按需做多语言。
	s.GetOpenApi().Info.Title = "RoamAdmin API"
	s.GetOpenApi().Info.Version = version
	s.GetOpenApi().Info.Description = "RoamAdmin backend REST API. All endpoints return a unified envelope and English messages; timestamps are RFC3339 UTC."

	// JWT Bearer 认证方案：API 文档界面会显示认证（bearerAuth）入口，
	// 填入 /auth/login 返回的 token 后即可在线调试受保护接口。
	oai := s.GetOpenApi()
	oai.Components.SecuritySchemes = goai.SecuritySchemes{
		"bearerAuth": goai.SecuritySchemeRef{
			Value: &goai.SecurityScheme{
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
				Description:  "JWT token from POST /auth/login",
			},
		},
	}
	oai.Security = &goai.SecurityRequirements{
		{"bearerAuth": {}},
	}

	// API 文档界面（Scalar，在线调试）：内嵌 @scalar/api-reference，零 CDN 依赖、离线可用。
	swaggerui.Init(s, g.Cfg().MustGet(ctx, "server.openapiPath", "/api.json").String())

	// 全局中间件（注册顺序即执行顺序，由外到内）：
	//   可观测（panic 隔离 + metrics + traceId）→ 认证（JWT 注入用户）→ 统一响应体。
	s.Use(ra.Observ().Middleware, ra.Auth().Middleware, errcode.Middleware)

	// Prometheus 指标抓取端点。
	s.BindHandler("GET:/metrics", func(r *ghttp.Request) {
		ra.Observ().Handler().ServeHTTP(r.Response.Writer, r.Request)
	})

	a := app.New("RoamAdmin", version)
	if err := a.Run(ctx); err != nil {
		g.Log().Fatal(ctx, err)
	}
}
