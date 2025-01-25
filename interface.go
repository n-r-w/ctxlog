package ctxlog

import "context"

//go:generate mockgen -source interface.go -destination interface_mock.go -package ctxlog

// ILogger is a logger interface to extract the logger implementation
// from the project to external dependencies.
// This package contains two implementations of the ILogger interface:
// 1. ctxlog.NewWrapper - a real implementation that allows wrapping ctxlog.Logger in this interface.
// 2. ctxlog.NewStubWrapper - a real implementation that allows logging in the context, but does not do anything.
// WARNING: args must be even in the form of a key-value pair for structured logging.
type ILogger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}
