package ctxlog

import (
	"context"
)

// NewWrapper returns a new wrapper that implements the ILogger interface.
func NewWrapper() ILogger {
	return &wrapper{}
}

type wrapper struct{}

func (w *wrapper) Debug(ctx context.Context, msg string, args ...any) {
	Debug(ctx, msg, args...)
}

func (w *wrapper) Info(ctx context.Context, msg string, args ...any) {
	Info(ctx, msg, args...)
}

func (w *wrapper) Warn(ctx context.Context, msg string, args ...any) {
	Warn(ctx, msg, args...)
}

func (w *wrapper) Error(ctx context.Context, msg string, args ...any) {
	Error(ctx, msg, args...)
}

type stubWrapper struct{}

// NewStubWrapper returns a new stub wrapper that implements the ILogger interface
// and does nothing.
func NewStubWrapper() ILogger {
	return &stubWrapper{}
}

func (w *stubWrapper) Debug(_ context.Context, _ string, _ ...any) {
}

func (w *stubWrapper) Info(_ context.Context, _ string, _ ...any) {
}

func (w *stubWrapper) Warn(_ context.Context, _ string, _ ...any) {
}

func (w *stubWrapper) Error(_ context.Context, _ string, _ ...any) {
}
