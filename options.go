package ctxlog

import (
	"log/slog"
	"testing"
	"time"
)

// Option is a function for configuring the logger.
type Option func(*options)

// WithDevelopment (default) sets the logger to Development mode.
// Development uses a readable format with colored level display.
func WithDevelopment() Option {
	return func(o *options) {
		o.env = developmentEnv
	}
}

// WithProduction sets the logger to Production mode.
// Production uses a compact JSON format without colored level display.
func WithProduction() Option {
	return func(o *options) {
		o.env = productionEnv
	}
}

// WithLevel sets the minimum logging level.
func WithLevel(level slog.Leveler) Option {
	return func(o *options) {
		o.level = level
	}
}

// WithName sets the logger name.
func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// WithSource adds the file name and line number of the call to the log record.
func WithSource() Option {
	return func(o *options) {
		o.addSource = true
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
