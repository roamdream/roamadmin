// Package swaggerui 提供内嵌的 API 文档界面（Scalar，含在线调试）。
//
// 渲染引擎采用 Scalar API Reference（MIT，https://github.com/scalar/scalar），
// 界面现代美观、自带深色主题，并提供 "Test Request" 在线调试（可直接在浏览器里
// 发起请求）。静态资源编译期内嵌、运行时零 CDN 依赖、离线可用。
//
// API 文档内容本身由 GoFrame 依据各接口的 g.Meta 标签自动生成
// （server.openapiPath，如 /api.json），本包只负责渲染与提供在线调试界面。
package swaggerui

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

//go:embed all:dist
var distFS embed.FS

// specURLPlaceholder index.html 中 OpenAPI 规范地址的占位符。
const specURLPlaceholder = "{SPEC_URL}"

// Init 将内嵌的 API 文档界面注册到 /swagger/*。
// specURL 为 OpenAPI 规范地址（来自配置 server.openapiPath，如 /api.json）。
func Init(s *ghttp.Server, specURL string) {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	indexHTML := renderIndex(specURL)
	fileServer := http.FileServer(http.FS(sub))

	s.BindHandler("/swagger/*", func(r *ghttp.Request) {
		path := r.Request.URL.Path
		// 首页：GoFrame 会把 /swagger 与 /swagger/ 归一化为 /swagger。
		if path == "/swagger" || path == "/swagger/" || path == "/swagger/index.html" {
			r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
			r.Response.WriteString(indexHTML)
			return
		}
		// 静态资源直接走底层 http.FileServer（RawWriter 绕过统一 JSON 响应体）。
		http.StripPrefix("/swagger", fileServer).ServeHTTP(r.Response.RawWriter(), r.Request)
	})
}

// renderIndex 读取首页模板并注入 OpenAPI 规范地址。
func renderIndex(specURL string) string {
	raw, err := distFS.ReadFile("dist/index.html")
	if err != nil {
		panic(err)
	}
	return strings.ReplaceAll(string(raw), specURLPlaceholder, specURL)
}
