// Package webui 提供内嵌的管理后台前端（Vue3 + Element Plus）。
//
// 构建产物（web/ 下执行 npm run build）输出到本包的 dist/ 目录，经 go:embed 编译进
// 二进制，运行时挂在 /admin 前缀，与后端同源、单二进制部署、零 CORS。
//
// 前端路由采用 history 模式（base /admin/）：/admin 与 /admin/* 命中 SPA 回退逻辑，
// 静态资源（/admin/assets/*）优先按真实文件返回，其余路径统一回 index.html。
package webui

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

//go:embed all:dist
var distFS embed.FS

// Init 将内嵌的管理后台注册到 /admin/*。
func Init(s *ghttp.Server) {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	indexHTML, err := fs.ReadFile(sub, "index.html")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(sub))

	s.BindHandler("/admin/*", func(r *ghttp.Request) {
		path := r.Request.URL.Path
		rel := strings.TrimPrefix(path, "/admin")

		// /admin 与 /admin/ 首页：直接回 index.html。
		if rel == "" || rel == "/" {
			writeHTML(r, indexHTML)
			return
		}

		// 存在真实静态文件（如 /admin/assets/xxx.js）则按文件返回。
		name := strings.TrimPrefix(rel, "/")
		if _, err := fs.Stat(sub, name); err == nil {
			http.StripPrefix("/admin", fileServer).ServeHTTP(r.Response.RawWriter(), r.Request)
			return
		}

		// 其余路径（前端 history 路由）回退到 index.html。
		writeHTML(r, indexHTML)
	})
}

func writeHTML(r *ghttp.Request, html []byte) {
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Response.WriteString(string(html))
}
