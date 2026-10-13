package errcode

// Response 统一对外响应体，由内核中间件统一写回 HTTP。
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	TraceID string `json:"traceId,omitempty"`
}

// OK 构造成功响应。
func OK(data any) *Response {
	return &Response{Code: int(CodeOK), Message: CodeOK.Message(), Data: data}
}

// Fail 构造失败响应。
func Fail(code Code, msg string) *Response {
	if msg == "" {
		msg = code.Message()
	}
	return &Response{Code: int(code), Message: msg}
}

// FailErr 从 *Error 构造失败响应。
func FailErr(e *Error) *Response {
	if e == nil {
		return Fail(CodeInternal, "")
	}
	return Fail(e.Code, e.Msg)
}

// WithTraceID 设置链路追踪 ID。
func (r *Response) WithTraceID(id string) *Response {
	r.TraceID = id
	return r
}
