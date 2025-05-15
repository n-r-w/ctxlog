package ctxlog

import (
	"context"
	"fmt"
)

// GolangMigrateLogger is a logger for golang-migrate.
type GolangMigrateLogger struct {
	l ILogger
}

// NewGolangMigrateLogger creates a new golang-migrate logger.
func NewGolangMigrateLogger(l ILogger) *GolangMigrateLogger {
	return &GolangMigrateLogger{l: l}
}

// Printf logs a message.
func (g *GolangMigrateLogger) Printf(format string, v ...any) {
	g.l.Info(context.Background(), fmt.Sprintf(format, v...))
}

// Verbose returns true.
func (g *GolangMigrateLogger) Verbose() bool {
	return true
}
