package ctxlog

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
)

// EnvType is a logger environment mode.
type EnvType int

const (
	// EnvDevelopment (default) sets the logger to Development mode.
	// Uses a readable format with colored level display.
	EnvDevelopment EnvType = iota

	// EnvProduction sets the logger to Production mode.
	// Uses a compact JSON format without colored level display.
	EnvProduction
)

// EnvTypeFromString returns the EnvType for a string.
// Accepts "DEV" and "DEVELOPMENT" for Development mode,
// and "PROD" and "PRODUCTION" for Production mode.
func EnvTypeFromString(s string) (EnvType, error) {
	switch s {
	case "DEV", "DEVELOPMENT":
		return EnvDevelopment, nil
	case "PROD", "PRODUCTION":
		return EnvProduction, nil
	}

	return EnvDevelopment, fmt.Errorf("unknown environment type: %s", s)
}

// Option is a function for configuring the logger.
type Option func(*options)

// WithEnvType sets the logger environment mode.
// default: EnvProduction.
func WithEnvType(env EnvType) Option {
	return func(o *options) {
		o.env = env
	}
}

// WithLevel sets the minimum logging level.
// default: LevelDebug.
func WithLevel(level slog.Leveler) Option {
	return func(o *options) {
		o.level = level
	}
}

// WithName adds a name to logger output.
// default: empty string.
func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// WithSource adds the file name and line number of the call to the log record.
// default: true.
func WithSource(b bool) Option {
	return func(o *options) {
		o.addSource = b
	}
}

// WithSampler sets the sampler for the logger.
func WithSampler(tick time.Duration, first, thereafter int) Option {
	return func(o *options) {
		o.samplingTick = tick
		o.samplingFirst = first
		o.samplingThereafter = thereafter
	}
}

// WithOtelTracing sets up the logger to use OpenTelemetry.
// default: false.
func WithOtelTracing() Option {
	return func(o *options) {
		o.otel = true
	}
}

// WithTesting sets up the logger for use in tests.
func WithTesting(t testing.TB) Option {
	return func(o *options) {
		o.testTB = t
	}
}

// WithTimeLayout sets the time layout for the logger.
// default: time.RFC3339Nano.
func WithTimeLayout(layout string) Option {
	return func(o *options) {
		o.timeLayout = layout
	}
}

// WithTestBuffer sets a buffer for capturing test output.
// This is useful for testing log output.
// Requires the WithTesting option to be set.
func WithTestBuffer(buffer *zaptest.Buffer) Option {
	return func(o *options) {
		o.testBuffer = buffer
	}
}
