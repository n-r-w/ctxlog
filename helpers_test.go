package ctxlog

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantLevel slog.Level
		wantError bool
		errorMsg  string
	}{
		{
			name:      "debug uppercase",
			input:     "DEBUG",
			wantLevel: slog.LevelDebug,
			wantError: false,
			errorMsg:  "",
		},
		{
			name:      "debug lowercase",
			input:     "debug",
			wantLevel: slog.LevelDebug,
			wantError: false,
			errorMsg:  "",
		},
		{
			name:      "debug mixed case",
			input:     "DeBuG",
			wantLevel: slog.LevelDebug,
			wantError: false,
			errorMsg:  "",
		},
		{
			name:      "debug with spaces",
			input:     "  DEBUG  ",
			wantLevel: slog.LevelDebug,
			wantError: false,
			errorMsg:  "",
		},
		{
			name:      "info uppercase",
			input:     "INFO",
			wantLevel: slog.LevelInfo,
			wantError: false,
			errorMsg:  "",
		},
		{
			name:      "info lowercase",
			input:     "info",
			wantLevel: slog.LevelInfo,
			wantError: false,
			errorMsg:  "",
		},
		{
			name:      "warn uppercase",
			input:     "WARN",
			wantLevel: slog.LevelWarn,
			wantError: false,
			errorMsg:  "",
		},
		{
			name:      "warning uppercase",
			input:     "WARNING",
			wantLevel: slog.LevelWarn,
			wantError: false,
			errorMsg:  "",
		},
		{
			name:      "warning lowercase",
			input:     "warning",
			wantLevel: slog.LevelWarn,
			wantError: false,
			errorMsg:  "",
		},
		{
			name:      "error uppercase",
			input:     "ERROR",
			wantLevel: slog.LevelError,
			wantError: false,
			errorMsg:  "",
		},
		{
			name:      "error lowercase",
			input:     "error",
			wantLevel: slog.LevelError,
			wantError: false,
			errorMsg:  "",
		},
		{
			name:      "invalid level",
			input:     "INVALID",
			wantLevel: slog.LevelInfo,
			wantError: true,
			errorMsg:  "unknown log level",
		},
		{
			name:      "empty string",
			input:     "",
			wantLevel: slog.LevelInfo,
			wantError: true,
			errorMsg:  "unknown log level",
		},
		{
			name:      "only spaces",
			input:     "   ",
			wantLevel: slog.LevelInfo,
			wantError: true,
			errorMsg:  "unknown log level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotLevel, gotErr := ParseLogLevel(tt.input)

			if tt.wantError {
				require.Error(t, gotErr, "ParseLogLevel should return error for input %q", tt.input)
				assert.Contains(t, gotErr.Error(), tt.errorMsg,
					"error message should contain %q", tt.errorMsg)
			} else {
				require.NoError(t, gotErr, "ParseLogLevel should not return error for input %q", tt.input)
			}

			assert.Equal(t, tt.wantLevel, gotLevel,
				"ParseLogLevel(%q) level should match expected", tt.input)
		})
	}
}
