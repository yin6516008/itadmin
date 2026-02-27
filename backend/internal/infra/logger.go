// Package infra 提供基础设施初始化（数据库、日志）和日志上下文工具。
package infra

import (
	"context"
	"log/slog"
	"os"
)

// InitLogger 根据运行环境初始化全局 logger。
// production 环境使用 JSON 格式、Info 级别；其他环境使用 Text 格式、Debug 级别。
func InitLogger(env string) {
	opts := &slog.HandlerOptions{AddSource: true}
	if env == "production" {
		opts.Level = slog.LevelInfo
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
	} else {
		opts.Level = slog.LevelDebug
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
	}
}

type logCtxKey struct{}

// WithLogger 将带请求级字段的 logger 存入 context。
func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, logCtxKey{}, l)
}

// Logger 从 context 取出 logger，不存在时返回 slog.Default()。
func Logger(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(logCtxKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
