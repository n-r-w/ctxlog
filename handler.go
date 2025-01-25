package ctxlog

import (
	"context"
	"log/slog"
	"runtime"
	"strings"

	"go.uber.org/zap/buffer"
)

type handler struct {
	slog.Handler

	level     slog.Leveler
	logSource bool
}

func newHandler(h slog.Handler, level slog.Leveler, logSource bool) slog.Handler {
	return handler{
		Handler:   h,
		level:     level,
		logSource: logSource,
	}
}

type callStackSkipKey struct{}

// ctxCallStackSkipKey is a context key that skips the call stack.
var ctxCallStackSkipKey = callStackSkipKey{} //nolint:gochecknoglobals // ok for context

// Handle adds attributes from context to the record and then calls the handler.
func (h handler) Handle(ctx context.Context, r slog.Record) error {
	if h.logSource && r.PC != 0 {
		const (
			maxCallers  = 100
			defaultSkip = 4
		)

		skip, ok := ctx.Value(ctxCallStackSkipKey).(int)
		if !ok {
			skip = defaultSkip
		}

		pc := make([]uintptr, maxCallers)
		nCallers := runtime.Callers(skip, pc)
		frame, _ := runtime.CallersFrames(pc[:nCallers]).Next()

		r.AddAttrs(slog.String("source", trimmedPath(frame.File, frame.Line)))
	}

	return h.Handler.Handle(ctx, r)
}

// taken from zapcore.EntryCaller.TrimmedPath.
var _pool = buffer.NewPool() //nolint:gochecknoglobals // singleton

func trimmedPath(file string, line int) string {
	idx := strings.LastIndexByte(file, '/')
	if idx == -1 {
		return file
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(file[:idx], '/')
	if idx == -1 {
		return file
	}
	buf := _pool.Get()
	// Keep everything after the penultimate separator.
	buf.AppendString(file[idx+1:])
	buf.AppendByte(':')
	buf.AppendInt(int64(line))
	caller := buf.String()
	buf.Free()
	return caller
}

// Enabled reports whether records at the specified level should be processed.
func (h handler) Enabled(_ context.Context, level slog.Level) bool {
	return h.level.Level() <= level
}

// WithAttrs returns a new handler that has attributes from both handlers.
func (h handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newHandler(
		h.Handler.WithAttrs(attrs),
		h.level,
		h.logSource)
}

// WithGroup returns a new handler with a group added to the handler.
func (h handler) WithGroup(group string) slog.Handler {
	return newHandler(
		h.Handler.WithGroup(group),
		h.level,
		h.logSource)
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
