package ctxlog

import (
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestLogger_OutputFormat verifies that the logger correctly formats messages
// based on environment type and log level.
func TestLogger_OutputFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string     // Test case name.
		env         EnvType    // Environment type (development/production).
		message     string     // Message to log.
		level       slog.Level // Log level.
		wantContain []string   // Strings that should be present in the output.
	}{
		{
			name:    "development format debug message",
			env:     EnvDevelopment,
			message: "test debug message",
			level:   slog.LevelDebug,
			wantContain: []string{
				"DEBUG",
				"test debug message",
			},
		},
		{
			name:    "production format info message",
			env:     EnvProduction,
			message: "test info message",
			level:   slog.LevelInfo,
			wantContain: []string{
				"INFO",
				"test info message",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a test buffer to capture output.
			buffer := &zaptest.Buffer{}

			// Create test logger with buffer and verify setup.
			logger, err := New(
				WithEnvType(tt.env),
				WithLevel(tt.level),
				WithTesting(t),
				WithTestBuffer(buffer),
			)
			require.NoError(t, err, "failed to create logger")
			require.NotNil(t, logger, "logger should not be nil")

			// Log a message using the appropriate level.
			ctx := context.Background()
			switch tt.level {
			case slog.LevelDebug:
				logger.Debug(ctx, tt.message)
			case slog.LevelInfo:
				logger.Info(ctx, tt.message)
			case slog.LevelWarn:
				logger.Warn(ctx, tt.message)
			case slog.LevelError:
				logger.Error(ctx, tt.message)
			}

			require.NoError(t, logger.Sync())

			// Get and verify the output.
			output := buffer.String()
			for _, want := range tt.wantContain {
				require.Contains(t, output, want,
					"log output should contain %q in environment %v at level %v",
					want, tt.env, tt.level)
			}
		})
	}
}

// TestLogger_WithAttributes tests Logger.With.
func TestLogger_WithAttributes(t *testing.T) {
	t.Parallel()

	// Create a test buffer to capture output
	buffer := &zaptest.Buffer{}

	// Create logger with test buffer
	logger, err := New(
		WithEnvType(EnvDevelopment),
		WithLevel(slog.LevelDebug),
		WithTesting(t),
		WithTestBuffer(buffer),
	)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Test with attributes
	withAttrs := logger.With(
		"string", "value",
		"int", 42,
		"bool", true,
	)

	ctx := context.Background()
	withAttrs.Info(ctx, "test message with attributes")

	require.NoError(t, logger.Sync())

	// Get and verify the output
	logOutput := buffer.String()
	require.Contains(t, logOutput, "value", "output should contain string attribute value")
	require.Contains(t, logOutput, "42", "output should contain int attribute value")
	require.Contains(t, logOutput, "true", "output should contain bool attribute value")
}

// TestLogger_WithGroup tests the WithGroup method behavior.
func TestLogger_WithGroup(t *testing.T) {
	t.Parallel()

	// Test cases
	tests := []struct {
		name      string
		groupName string
		message   string
		wantGroup bool // Whether the group should appear in output
	}{
		{
			name:      "with empty group name",
			groupName: "",
			message:   "message without group",
			wantGroup: false,
		},
		{
			name:      "with non-empty group name",
			groupName: "test_group",
			message:   "message with group",
			wantGroup: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a test buffer to capture output
			buffer := &zaptest.Buffer{}

			// Create base logger
			logger, err := New(
				WithEnvType(EnvDevelopment),
				WithLevel(slog.LevelDebug),
				WithTesting(t),
				WithTestBuffer(buffer),
			)
			require.NoError(t, err)
			require.NotNil(t, logger)

			// Create logger with group
			groupLogger := logger.WithGroup(tt.groupName)

			// Verify logger instance
			if tt.groupName == "" {
				require.Equal(t, logger, groupLogger, "empty group name should return same logger")
			} else {
				require.NotEqual(t, logger, groupLogger, "non-empty group name should return new logger")
			}

			// Log a message
			ctx := context.Background()
			groupLogger.Info(ctx, tt.message)

			require.NoError(t, groupLogger.Sync())

			// Verify output
			output := buffer.String()
			if tt.wantGroup {
				require.Contains(t, output, tt.groupName, "output should contain group name")
			}
			require.Contains(t, output, tt.message, "output should contain message")
		})
	}
}

// TestLogger_WithTimeLayout tests the WithTimeLayout option behavior.
func TestLogger_WithTimeLayout(t *testing.T) {
	t.Parallel()

	// Create a test buffer to capture output
	buffer := &zaptest.Buffer{}

	// Use a custom time layout that includes timezone
	customLayout := "2006-01-02T15:04:05.999-0700"

	// Create logger with custom time layout
	logger, err := New(
		WithEnvType(EnvDevelopment),
		WithLevel(slog.LevelDebug),
		WithTimeLayout(customLayout),
		WithTesting(t),
		WithTestBuffer(buffer),
	)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Log a message
	ctx := context.Background()
	logger.Info(ctx, "test message")

	require.NoError(t, logger.Sync())

	// Get the output and verify the time format
	output := buffer.String()

	// Extract timestamp from the beginning of the line up to the first tab
	timestamp := strings.Split(output, "\t")[0]

	// Parse the timestamp to verify it's in the expected format
	_, err = time.Parse(customLayout, timestamp)
	require.NoError(t, err, "timestamp should be in the custom format %q, got: %q", customLayout, timestamp)
}

// TestLogger_Sampling tests Logger.WithSampler.
func TestLogger_Sampling(t *testing.T) {
	t.Parallel()

	// Create a buffer to capture log output
	buf := &zaptest.Buffer{}

	logger, err := New(
		WithEnvType(EnvDevelopment),
		WithLevel(slog.LevelDebug),
		WithTesting(t),
		WithTestBuffer(buf),                     // Use buffer to capture output
		WithSampler(100*time.Millisecond, 1, 0), // Allow 1 initial message, then drop all messages in the same tick
	)
	require.NoError(t, err)
	require.NotNil(t, logger)

	ctx := context.Background()

	// Log messages rapidly
	for i := range 5 {
		logger.Info(ctx, "test message", "iteration", i)
	}

	// Ensure logs are flushed
	require.NoError(t, logger.Sync())

	// Get captured output
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Verify sampling behavior:
	// 1. Should have exactly 1 message (initial allowed message)
	require.Len(t, lines, 1, "should allow only initial message")

	// 2. Verify the content of the first message
	require.Contains(t, lines[0], `"iteration": 0`,
		"message should contain correct iteration number")
}

// TestLogger_WithName tests the WithName option behavior.
func TestLogger_WithName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		loggerName string
		message    string
	}{
		{
			name:       "with empty name",
			loggerName: "",
			message:    "message without logger name",
		},
		{
			name:       "with non-empty name",
			loggerName: "test_logger",
			message:    "message with logger name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a test buffer to capture output.
			buffer := &zaptest.Buffer{}

			// Create logger with test buffer and name.
			logger, err := New(
				WithEnvType(EnvDevelopment),
				WithLevel(slog.LevelDebug),
				WithTesting(t),
				WithTestBuffer(buffer),
				WithName(tt.loggerName),
			)
			require.NoError(t, err)
			require.NotNil(t, logger)

			// Log a message.
			ctx := context.Background()
			logger.Info(ctx, tt.message)

			require.NoError(t, logger.Sync())

			// Verify output.
			output := buffer.String()
			if tt.loggerName != "" {
				require.Contains(t, output, tt.loggerName, "output should contain logger name")
			}
			require.Contains(t, output, tt.message, "output should contain message")
		})
	}
}

// TestLogger_WithSource tests the WithSource option behavior.
func TestLogger_WithSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		addSource bool
		message   string
	}{
		{
			name:      "with source enabled",
			addSource: true,
			message:   "message with source",
		},
		{
			name:      "with source disabled",
			addSource: false,
			message:   "message without source",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a test buffer to capture output.
			buffer := &zaptest.Buffer{}

			// Create logger with test buffer and source option.
			logger, err := New(
				WithEnvType(EnvDevelopment),
				WithLevel(slog.LevelDebug),
				WithTesting(t),
				WithTestBuffer(buffer),
				WithSource(tt.addSource),
			)
			require.NoError(t, err)
			require.NotNil(t, logger)

			// Log a message.
			ctx := context.Background()
			logger.Info(ctx, tt.message)

			require.NoError(t, logger.Sync())

			// Verify output.
			output := buffer.String()
			if tt.addSource {
				require.Contains(t, output, "log_test.go", "output should contain source file name")
			} else {
				require.NotContains(t, output, "log_test.go", "output should not contain source file name")
			}
			require.Contains(t, output, tt.message, "output should contain message")
		})
	}
}

// TestLogger_Context tests the context-related logger methods.
func TestLogger_Context(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := Must(WithTesting(t))

	// Test ToContext and FromContext
	t.Run("ToContext and FromContext", func(t *testing.T) {
		t.Parallel()

		ctxWithLogger := ToContext(ctx, logger)
		require.True(t, InContext(ctxWithLogger), "logger should be in context")

		retrievedLogger := FromContext(ctxWithLogger)
		require.Equal(t, logger, retrievedLogger, "retrieved logger should match original")
	})

	// Test TryFromContext
	t.Run("TryFromContext", func(t *testing.T) {
		t.Parallel()

		ctxWithLogger := ToContext(ctx, logger)
		retrievedLogger, ok := TryFromContext(ctxWithLogger)
		require.True(t, ok, "TryFromContext should return true")
		require.Equal(t, logger, retrievedLogger, "retrieved logger should match original")

		_, ok = TryFromContext(ctx)
		require.False(t, ok, "TryFromContext should return false for context without logger")
	})

	// Test InContext
	t.Run("InContext", func(t *testing.T) {
		t.Parallel()

		ctxWithLogger := ToContext(ctx, logger)
		require.True(t, InContext(ctxWithLogger), "InContext should return true for context with logger")
		require.False(t, InContext(ctx), "InContext should return false for context without logger")
	})
}
