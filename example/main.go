//nolint:mnd // ok for example
package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/n-r-w/ctxlog"
)

func main() {
	ctx := ctxlog.MustContext(
		context.Background(),
		ctxlog.WithName("myapp"),
		ctxlog.WithDevelopment(),
		ctxlog.WithLevel(slog.LevelDebug),
	)
	defer func() { _ = ctxlog.Sync(ctx) }()

	// Each output using the logger will add the path to the file and line number

	// Working with logger in context
	WorkingWithContextExample(ctx)

	// Extracting from context and direct usage of logger
	ExtractFromContext(ctx)

	// Logging errors with automatic stack trace
	ctxlog.Error(ctx, "failed to connect to database", "host", "localhost", "port", 5432)

	// CloseError close io.Closer and log any error
	ctxlog.CloseError(ctx, &SomeCloser{})
}

func WorkingWithContextExample(ctx context.Context) {
	// Basic logging
	ctxlog.Info(ctx, "starting application", "version", "1.0.0")

	// With slog.Attributes
	ctxlog.Debug(ctx, "connecting to database", slog.String("host", "localhost"), slog.Int("port", 5432))

	// Logging with groups
	ctx = ctxlog.WithGroup(ctx, "database")
	ctxlog.Debug(ctx, "connecting to database", "host", "localhost", "port", 5432)

	// Logging with additional context
	ctx = ctxlog.With(ctx,
		"user_id", "12345",
		"ip", "192.168.1.1",
	)
	ctxlog.Info(ctx, "user logged in")
}

func ExtractFromContext(ctx context.Context) {
	// Extracting from context
	logger := ctxlog.FromContext(ctx)
	logger.Info("direct usage of logger")
}

type SomeCloser struct{}

func (*SomeCloser) Close() error {
	return errors.New("some error")
}
