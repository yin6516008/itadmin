// Package errcode 定义统一的业务错误码，供 handler 和 service 层使用。
package errcode

// AppError 表示可直接返回给前端的业务错误，包含 HTTP 状态码和业务码。
type AppError struct {
	HTTPStatus int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

// Error 实现 error 接口。
func (e *AppError) Error() string { return e.Message }

var (
	// ErrInvalidParam 表示参数校验失败。
	ErrInvalidParam = &AppError{HTTPStatus: 400, Code: 40001, Message: "参数校验失败"}
	// ErrUnauthorized 表示用户未登录。
	ErrUnauthorized = &AppError{HTTPStatus: 401, Code: 40100, Message: "未登录"}
	// ErrForbidden 表示用户无权限。
	ErrForbidden = &AppError{HTTPStatus: 403, Code: 40300, Message: "无权限"}
	// ErrNotFound 表示请求的资源不存在。
	ErrNotFound = &AppError{HTTPStatus: 404, Code: 40400, Message: "资源不存在"}
	// ErrInternal 表示服务器内部错误。
	ErrInternal = &AppError{HTTPStatus: 500, Code: 50000, Message: "服务器内部错误"}
)
