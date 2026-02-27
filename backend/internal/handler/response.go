// Package handler 实现 HTTP 路由注册、参数校验和统一响应。
package handler

import (
	"errors"
	"net/http"

	"backend/internal/errcode"

	"github.com/gin-gonic/gin"
)

type responseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success 返回统一的成功 JSON 响应，code=0。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, responseBody{Code: 0, Message: "ok", Data: data})
}

// Error 返回统一的错误 JSON 响应，自动识别 AppError 提取状态码和业务码。
func Error(c *gin.Context, err error) {
	var appErr *errcode.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, responseBody{Code: appErr.Code, Message: appErr.Message})
	} else {
		c.JSON(http.StatusInternalServerError, responseBody{Code: 50000, Message: "服务器内部错误"})
	}
	c.Abort()
}
