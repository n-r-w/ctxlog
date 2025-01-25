// Package ctxlog provides a slog logger with zap/otelzap backend.
package ctxlog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"syscall"
	"testing"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

// Logger is a slog with zap backend.
type Logger struct {
	*slog.Logger

	opts       options
	zapLogger  *zap.Logger
	otelLogger *otelzap.Logger
}

type options struct {
	env                EnvType
	level              slog.Leveler
	addSource          bool
	name               string
	testTB             testing.TB
	samplingTick       time.Duration
	samplingFirst      int
	samplingThereafter int
	otel               bool
	timeLayout         string
}

// New creates a new logger.
// Default: Development mode, Debug level, with call source display.
func New(opts ...Option) (*Logger, error) {
	o := options{
		env:        EnvProduction,
		level:      slog.LevelDebug,
		addSource:  true,
		timeLayout: time.RFC3339Nano,
	}
	for _, opt := range opts {
		opt(&o)
	}

	zLevel := zapLevel(o.level)

	var zapConf zap.Config
	switch o.env {
	case EnvDevelopment:
		zapConf = zap.NewDevelopmentConfig()
		zapConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case EnvProduction:
		zapConf = zap.NewProductionConfig()
	}
	zapConf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(o.timeLayout)
	zapConf.Level = zap.NewAtomicLevelAt(zLevel)

	zapLogger, err := zapConf.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create zap logger: %w", err)
	}

	if o.testTB != nil {
		var zo []zap.Option
		if o.env == EnvDevelopment {
			zo = append(zo, zap.Development())
		}
		zapLogger = zaptest.NewLogger(o.testTB, zaptest.Level(zLevel), zaptest.WrapOptions(zo...))
	}

	return newLoggerHelper(zapLogger, o), nil
}

func newLoggerHelper(zapLogger *zap.Logger, opts options) *Logger {
	core := zapLogger.Core()
	if opts.samplingTick != 0 {
		core = zapcore.NewSamplerWithOptions(core, opts.samplingTick, opts.samplingFirst, opts.samplingThereafter)
	}

	slogLogger := slog.New(
		newHandler(
			zapslog.NewHandler(core,
				zapslog.WithName(opts.name),
				zapslog.WithCaller(false),
			),
			opts.level,
			opts.addSource,
		),
	)

	l := &Logger{
		Logger:    slogLogger,
		opts:      opts,
		zapLogger: zapLogger,
	}

	if opts.otel {
		l.otelLogger = otelzap.New(
			zapLogger,
			otelzap.WithStackTrace(true),
			otelzap.WithTraceIDField(true),
		)
	}

	return l
}

// Must panics if the logger cannot be created.
func Must(opts ...Option) *Logger {
	l, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return l
}

// Sync flushes buffered log entries. Use it in defer.
func (l *Logger) Sync() error {
	if err := l.zapLogger.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
		return fmt.Errorf("failed to sync logger: %w", err)
	}

	return nil
}

// With returns a logger that includes the specified attributes.
func (l *Logger) With(args ...any) *Logger {
	if len(args) == 0 {
		return l
	}

	return &Logger{
		Logger:    l.Logger.With(args...),
		zapLogger: l.zapLogger,
	}
}

// WithGroup returns a logger that starts a group if name is not empty.
func (l *Logger) WithGroup(name string) *Logger {
	if name == "" {
		return l
	}

	return &Logger{
		Logger:    l.Logger.WithGroup(name),
		zapLogger: l.zapLogger,
	}
}

// Debug is implement ILogger interface.
func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.LogWithLevel(ctx, slog.LevelDebug, msg, defaultSkipCallStack, args...)
}

// Info is implement ILogger interface.
func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.LogWithLevel(ctx, slog.LevelInfo, msg, defaultSkipCallStack, args...)
}

// Warn is implement ILogger interface.
func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.LogWithLevel(ctx, slog.LevelWarn, msg, defaultSkipCallStack, args...)
}

// Error is implement ILogger interface.
func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.LogWithLevel(ctx, slog.LevelError, msg, defaultSkipCallStack, args...)
}

// LogWithLevel logs a message at the specified level and stack frame skip.
func (l *Logger) LogWithLevel(ctx context.Context, level slog.Level, msg string, skip int, attrs ...any) {
	if _, ok := GetSkipCallStack(ctx); !ok {
		// skip the call stack
		ctx = SetSkipCallStack(ctx, skip)
	}

	l.Log(ctx, level, msg, attrs...)
}

func zapLevel(level slog.Leveler) zapcore.Level {
	switch level {
	case slog.LevelDebug:
		return zap.DebugLevel
	case slog.LevelInfo:
		return zap.InfoLevel
	case slog.LevelWarn:
		return zap.WarnLevel
	case slog.LevelError:
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
