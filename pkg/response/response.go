package response

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 响应状态码常量
const (
	CodeSuccess      = 200
	CodeError        = 500
	CodeInvalidParam = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
)

// Success 成功响应
func Success(ctx context.Context, c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(ctx context.Context, c *app.RequestContext, message string, data interface{}) {
	c.JSON(consts.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(ctx context.Context, c *app.RequestContext, message string) {
	c.JSON(consts.StatusOK, Response{
		Code:    CodeError,
		Message: message,
	})
}

// ErrorWithCode 带状态码的错误响应
func ErrorWithCode(ctx context.Context, c *app.RequestContext, code int, message string) {
	var httpStatus int
	switch code {
	case CodeInvalidParam:
		httpStatus = http.StatusBadRequest
	case CodeUnauthorized:
		httpStatus = http.StatusUnauthorized
	case CodeForbidden:
		httpStatus = http.StatusForbidden
	case CodeNotFound:
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

// InvalidParam 参数错误响应
func InvalidParam(ctx context.Context, c *app.RequestContext, message string) {
	ErrorWithCode(ctx, c, CodeInvalidParam, message)
}

// Unauthorized 未授权响应
func Unauthorized(ctx context.Context, c *app.RequestContext, message string) {
	ErrorWithCode(ctx, c, CodeUnauthorized, message)
}

// Forbidden 禁止访问响应
func Forbidden(ctx context.Context, c *app.RequestContext, message string) {
	ErrorWithCode(ctx, c, CodeForbidden, message)
}

// NotFound 未找到响应
func NotFound(ctx context.Context, c *app.RequestContext, message string) {
	ErrorWithCode(ctx, c, CodeNotFound, message)
}
