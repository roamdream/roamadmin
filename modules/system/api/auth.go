// Package api system 模块对外契约（请求/响应 DTO）。
package api

import "github.com/gogf/gf/v2/frame/g"

// LoginReq 登录请求。
type LoginReq struct {
	g.Meta   `path:"/system/login" method:"post" tags:"system" summary:"Login"`
	Username string `json:"username" v:"required|length:1,64#username is required|username length must be 1-64" dc:"Username"`
	Password string `json:"password" v:"required|length:1,64#password is required|password length must be 1-64" dc:"Password"`
}

// LoginRes 登录响应。
type LoginRes struct {
	Token  string   `json:"token"  dc:"Access token"`
	UserID int64    `json:"userId" dc:"User ID"`
	Name   string   `json:"name"   dc:"Nickname"`
	Roles  []string `json:"roles"  dc:"Role codes"`
	Perms  []string `json:"perms"  dc:"Permission codes"`
}

// UserInfoReq 当前用户信息请求。
type UserInfoReq struct {
	g.Meta `path:"/system/user/info" method:"get" tags:"system" summary:"Get current user info"`
}

// UserInfoRes 当前用户信息响应。
type UserInfoRes struct {
	UserID   int64    `json:"userId"   dc:"User ID"`
	Username string   `json:"username" dc:"Username"`
	Nickname string   `json:"nickname" dc:"Nickname"`
	Email    string   `json:"email"    dc:"Email"`
	Phone    string   `json:"phone"    dc:"Phone"`
	Avatar   string   `json:"avatar"   dc:"Avatar"`
	Roles    []string `json:"roles"    dc:"Role codes"`
	Perms    []string `json:"perms"    dc:"Permission codes"`
}

// LogoutReq 登出请求。
type LogoutReq struct {
	g.Meta `path:"/system/logout" method:"post" tags:"system" summary:"Logout"`
}

// LogoutRes 登出响应。
type LogoutRes struct {
	Logout bool `json:"logout" dc:"Whether logged out"`
}
