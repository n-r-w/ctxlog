package ctxlog

import (
	"context"
	"testing"
)

type logFieldsKey struct{}

var ctxLogFieldsKey = logFieldsKey{} //nolint:gochecknoglobals // ok for context

// ToContext adds a logger to the context.
// If the context already has a logger, it will be replaced.
func ToContext(ctx context.Context, logger *Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, ctxLogFieldsKey, logger)
}

// ToContextFromContext extracts the logger from ctxWithLogger and adds it to targetCtx.
func ToContextFromContext(targetCtx, ctxWithLogger context.Context) context.Context {
	return ToContext(targetCtx, FromContext(ctxWithLogger))
}

// FromContext returns the logger from the context.
// If the context does not have a logger, it panics.
func FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		panic("context is nil")
	}

	if l, ok := TryFromContext(ctx); ok {
		return l
	}

	panic("context does not have a logger")
}

// TryFromContext attempts to get a logger from the context.
func TryFromContext(ctx context.Context) (*Logger, bool) {
	if ctx == nil {
		return nil, false
	}

	if l, ok := ctx.Value(ctxLogFieldsKey).(*Logger); ok {
		return l, true
	}

	return nil, false
}

// InContext returns true if the context has a logger.
func InContext(ctx context.Context) bool {
	_, ok := TryFromContext(ctx)
	return ok
}

// NewContext returns a new context with the logger.
func NewContext(ctx context.Context, opts ...Option) (context.Context, error) {
	l, err := New(opts...)
	if err != nil {
		return nil, err
	}
	return ToContext(ctx, l), nil
}

// MustContext returns a new context with the logger. Panics if the logger cannot be created.
func MustContext(ctx context.Context, opts ...Option) context.Context {
	l := Must(opts...)
	return ToContext(ctx, l)
}

// ToTestContext returns a context where logs will be written to testing.TB.
func ToTestContext(ctx context.Context, t testing.TB) context.Context {
	t.Helper()
	return ToContext(ctx, Must(WithTesting(t)))
}
