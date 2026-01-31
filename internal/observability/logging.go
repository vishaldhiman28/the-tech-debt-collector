package observability

import (
	"context"
	"log/slog"
	"os"
)

// Logger wraps structured logging
type Logger struct {
	inner *slog.Logger
}

// NewLogger creates logger
func NewLogger(verbose bool) *Logger {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return &Logger{
		inner: slog.New(handler),
	}
}

// Info logs info message
func (l *Logger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.inner.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

// Error logs error message
func (l *Logger) Error(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	attrs = append(attrs, slog.String("error", err.Error()))
	l.inner.LogAttrs(ctx, slog.LevelError, msg, attrs...)
}

// Debug logs debug message
func (l *Logger) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.inner.LogAttrs(ctx, slog.LevelDebug, msg, attrs...)
}
