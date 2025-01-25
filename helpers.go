package ctxlog

import (
	"context"
	"io"
	"log/slog"
	"reflect"
)

// Debug logs a message at Debug level. Uses context to get the logger.
func Debug(ctx context.Context, msg string, attrs ...any) {
	logHelper(ctx, slog.LevelDebug, msg, attrs...)
}

// Info logs a message at Info level. Uses context to get the logger.
func Info(ctx context.Context, msg string, attrs ...any) {
	logHelper(ctx, slog.LevelInfo, msg, attrs...)
}

// Warn logs a message at Warn level. Uses context to get the logger.
func Warn(ctx context.Context, msg string, attrs ...any) {
	logHelper(ctx, slog.LevelWarn, msg, attrs...)
}

// Error logs a message at Error level. Uses context to get the logger.
func Error(ctx context.Context, msg string, attrs ...any) {
	logHelper(ctx, slog.LevelError, msg, attrs...)
}

// Log logs a message at the specified level. Uses context to get the logger.
func Log(ctx context.Context, level slog.Level, msg string, attrs ...any) {
	logHelper(ctx, level, msg, attrs...)
}

// With returns a logger that includes the specified attributes.
func With(ctx context.Context, attrs ...any) context.Context {
	return ToContext(ctx, FromContext(ctx).With(attrs...))
}

// WithGroup returns a logger that starts a group if name is not empty.
func WithGroup(ctx context.Context, name string) context.Context {
	return ToContext(ctx, FromContext(ctx).WithGroup(name))
}

// SetSkipCallStack sets the number of stack frames to skip when logging.
func SetSkipCallStack(ctx context.Context, skip int) context.Context {
	return context.WithValue(ctx, ctxCallStackSkipKey, skip)
}

// GetSkipCallStack returns the number of stack frames to skip when logging.
func GetSkipCallStack(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(ctxCallStackSkipKey).(int)
	return v, ok
}

func logHelper(ctx context.Context, level slog.Level, msg string, attrs ...any) {
	logger := FromContext(ctx)

	if _, ok := GetSkipCallStack(ctx); !ok {
		// skip the call stack
		const skip = 6
		ctx = SetSkipCallStack(ctx, skip)
	}

	logger.Log(ctx, level, msg, attrs...)
}

// Sync synchronizes logging in the context.
func Sync(ctx context.Context) error {
	return FromContext(ctx).Sync()
}

// CloseError closes c, logging any error that occurs.
func CloseError(ctx context.Context, c io.Closer) {
	log := FromContext(ctx)
	t := reflect.TypeOf(c).String()
	if err := c.Close(); err != nil {
		log.Error("failed to close", slog.String("type", t), slog.Any("error", err))
	}
}
