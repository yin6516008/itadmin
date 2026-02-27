package handler

import (
	"log/slog"
	"net/http"

	"backend/internal/infra"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CORS 处理跨域请求，开发阶段允许所有来源。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Trace-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// RequestLogger 为每个请求注入带 trace_id、method、path 的 logger 到 context。
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.NewString()
		}

		logger := slog.Default().With(
			slog.String("trace_id", traceID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)

		ctx := infra.WithLogger(c.Request.Context(), logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
