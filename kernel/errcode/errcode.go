// Package errcode 定义全框架统一的错误码与响应体。
//
// 约定：错误码 0 表示成功；1xxxx 为框架通用错误；2xxxx 起由各业务模块自行扩展。
// 所有对外 HTTP 响应统一走 Response 结构，由内核中间件统一写回。
package errcode

import "fmt"

// Code 业务错误码。
type Code int

const (
	// CodeOK 成功。
	CodeOK Code = 0

	// —— 1xxxx 框架通用错误 ——
	CodeInternal     Code = 10000 // 内部错误
	CodeParam        Code = 10001 // 参数错误
	CodeUnauthorized Code = 10002 // 未认证/认证失败
	CodeForbidden    Code = 10003 // 无权限
	CodeNotFound     Code = 10004 // 资源不存在
	CodeConflict     Code = 10005 // 冲突（唯一键等）
	CodeTooMany      Code = 10006 // 请求过于频繁（被限流）
	CodeDB           Code = 10010 // 数据库错误
	CodeRedis        Code = 10011 // 缓存/Redis 错误
	CodeDownstream   Code = 10012 // 下游服务不可用（已熔断）

	// —— 2xxxx 业务错误（模块可继续扩展） ——
	CodeLoginFailed   Code = 20001 // 登录失败（账号或密码错误）
	CodeUserDisabled  Code = 20002 // 用户被禁用
	CodeTenantMissing Code = 20003 // 缺失租户上下文
	CodeDuplicate     Code = 20004 // 数据重复
)

var messages = map[Code]string{
	CodeOK:             "ok",
	CodeInternal:       "internal server error",
	CodeParam:          "invalid parameter",
	CodeUnauthorized:   "unauthorized",
	CodeForbidden:      "forbidden",
	CodeNotFound:       "not found",
	CodeConflict:       "conflict",
	CodeTooMany:        "too many requests",
	CodeDB:             "database error",
	CodeRedis:          "cache error",
	CodeDownstream:     "downstream unavailable",
	CodeLoginFailed:    "login failed",
	CodeUserDisabled:   "user disabled",
	CodeTenantMissing:  "tenant missing",
	CodeDuplicate:      "duplicate data",
}

// Message 返回错误码的默认文案。
func (c Code) Message() string {
	if m, ok := messages[c]; ok {
		return m
	}
	return "unknown error"
}

// Error 携带业务码的错误，可被 errors.Is / Unwrap 穿透。
type Error struct {
	Code Code
	Msg  string
	Err  error
}

func (e *Error) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Code.Message()
}

// Unwrap 支持 errors.Is / errors.As 透传原始错误。
func (e *Error) Unwrap() error { return e.Err }

// New 用业务码 + 文案构造错误。
func New(code Code, msg string) *Error { return &Error{Code: code, Msg: msg} }

// Wrap 用业务码包裹底层错误。
func Wrap(code Code, err error) *Error { return &Error{Code: code, Err: err} }

// Wrapf 用业务码 + 格式化文案包裹底层错误。
func Wrapf(code Code, err error, format string, args ...any) *Error {
	return &Error{Code: code, Msg: fmt.Sprintf(format, args...), Err: err}
}
