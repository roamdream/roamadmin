// Package api system 模块对外契约（请求/响应 DTO）。
package api

import "github.com/gogf/gf/v2/frame/g"

// LoginReq 登录请求。
type LoginReq struct {
	g.Meta   `path:"/system/login" method:"post" tags:"系统" summary:"用户登录"`
	Username string `json:"username" v:"required|length:1,64#用户名不能为空|用户名长度1-64" dc:"用户名"`
	Password string `json:"password" v:"required|length:1,64#密码不能为空|密码长度1-64" dc:"密码"`
}

// LoginRes 登录响应。
type LoginRes struct {
	Token  string   `json:"token"  dc:"访问令牌"`
	UserID int64    `json:"userId" dc:"用户 ID"`
	Name   string   `json:"name"   dc:"昵称"`
	Roles  []string `json:"roles"  dc:"角色编码"`
	Perms  []string `json:"perms"  dc:"权限码"`
}

// UserInfoReq 当前用户信息请求。
type UserInfoReq struct {
	g.Meta `path:"/system/user/info" method:"get" tags:"系统" summary:"当前用户信息"`
}

// UserInfoRes 当前用户信息响应。
type UserInfoRes struct {
	UserID   int64    `json:"userId"   dc:"用户 ID"`
	Username string   `json:"username" dc:"用户名"`
	Nickname string   `json:"nickname" dc:"昵称"`
	Email    string   `json:"email"    dc:"邮箱"`
	Phone    string   `json:"phone"    dc:"手机号"`
	Avatar   string   `json:"avatar"   dc:"头像"`
	Roles    []string `json:"roles"    dc:"角色编码"`
	Perms    []string `json:"perms"    dc:"权限码"`
}

// LogoutReq 登出请求。
type LogoutReq struct {
	g.Meta `path:"/system/logout" method:"post" tags:"系统" summary:"用户登出"`
}

// LogoutRes 登出响应。
type LogoutRes struct {
	Logout bool `json:"logout" dc:"是否已登出"`
}
