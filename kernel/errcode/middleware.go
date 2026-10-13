package errcode

import (
	"errors"
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// Middleware 统一响应体中间件：把 handler 的 (data, error) 包装成 Response。
//
// 约定：handler 返回 (data, error)；error 为 *Error 时按业务码返回，
// 否则按内部错误兜底。响应统一带 traceId，便于链路排查。
func Middleware(r *ghttp.Request) {
	r.Middleware.Next()

	// 已手动写响应（swagger、/metrics、文件流等）则跳过。
	if r.Response.BufferLength() > 0 || r.Response.BytesWritten() > 0 {
		return
	}

	var (
		err = r.GetError()
		res = r.GetHandlerResponse()
	)

	var resp *Response
	switch {
	case err != nil:
		var e *Error
		if errors.As(err, &e) {
			resp = FailErr(e)
		} else {
			// 非框架错误（如 gerror 校验/DB 错误）：映射业务码并透出原始文案。
			resp = Fail(codeFromGF(gerror.Code(err)), err.Error())
		}
	case r.Response.Status > 0 && r.Response.Status != http.StatusOK:
		resp = Fail(codeByStatus(r.Response.Status), "")
	default:
		resp = OK(res)
	}

	if tid := gctx.CtxId(r.Context()); tid != "" {
		resp.TraceID = tid
	}
	r.Response.WriteJson(resp)
}

func codeByStatus(status int) Code {
	switch status {
	case http.StatusUnauthorized:
		return CodeUnauthorized
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusTooManyRequests:
		return CodeTooMany
	default:
		return CodeInternal
	}
}

// codeFromGF 将 GoFrame 通用错误码映射为框架业务码。
func codeFromGF(c gcode.Code) Code {
	switch c {
	case gcode.CodeValidationFailed, gcode.CodeInvalidParameter,
		gcode.CodeMissingParameter, gcode.CodeInvalidRequest:
		return CodeParam
	case gcode.CodeDbOperationError:
		return CodeDB
	case gcode.CodeNotAuthorized:
		return CodeForbidden
	case gcode.CodeNotFound:
		return CodeNotFound
	case gcode.CodeServerBusy:
		return CodeTooMany
	default:
		return CodeInternal
	}
}
