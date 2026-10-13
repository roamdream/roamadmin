package auth

import (
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/roamdream/roamadmin/kernel/errcode"
)

// Middleware 解析 Authorization: Bearer <token> 并注入用户上下文。
// 未带 token 时不报错，仅不注入用户，由 RequireAuth / RequirePermission 决定是否拒绝。
func (m *Manager) Middleware(r *ghttp.Request) {
	if token := extractToken(r); token != "" {
		if claims, err := m.ParseToken(token); err == nil {
			r.SetCtx(WithUser(r.Context(), claims))
		}
	}
	r.Middleware.Next()
}

// RequireAuth 要求已登录。
func RequireAuth(r *ghttp.Request) {
	if _, ok := UserFrom(r.Context()); !ok {
		writeError(r, http.StatusUnauthorized, errcode.CodeUnauthorized)
		return
	}
	r.Middleware.Next()
}

// RequireRole 要求具备指定角色。
func RequireRole(roles ...string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		c, ok := UserFrom(r.Context())
		if !ok {
			writeError(r, http.StatusUnauthorized, errcode.CodeUnauthorized)
			return
		}
		for _, role := range roles {
			if c.HasRole(role) || c.IsSuper() {
				r.Middleware.Next()
				return
			}
		}
		writeError(r, http.StatusForbidden, errcode.CodeForbidden)
	}
}

// RequirePermission 要求具备指定权限码。
func RequirePermission(perms ...string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		c, ok := UserFrom(r.Context())
		if !ok {
			writeError(r, http.StatusUnauthorized, errcode.CodeUnauthorized)
			return
		}
		for _, perm := range perms {
			if c.HasPermission(perm) {
				r.Middleware.Next()
				return
			}
		}
		writeError(r, http.StatusForbidden, errcode.CodeForbidden)
	}
}

func extractToken(r *ghttp.Request) string {
	header := r.GetHeader("Authorization")
	if header == "" {
		return ""
	}
	parts := strings.Fields(header)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return header
}

func writeError(r *ghttp.Request, status int, code errcode.Code) {
	r.Response.WriteHeader(status)
	r.Response.WriteJson(errcode.Fail(code, ""))
	r.ExitAll()
}
