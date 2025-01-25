# CTXLog

A context-aware structured logging package for Go that provides a robust wrapper around `log/slog` with enhanced functionality and configuration options. Under the hood, it uses `go.uber.org/zap` as the logging backend for optimal performance and reliability.

## Package Overview

CTXLog is designed to simplify structured logging in Go applications by providing context-aware logging capabilities with support for both development and production environments. It seamlessly integrates with Go's standard `log/slog` package while adding powerful features like context propagation, sampling, and OpenTelemetry integration. The package leverages `go.uber.org/zap` as its logging backend, combining the convenience of `slog`'s interface with Zap's high-performance logging capabilities.

## Key Features

- Context-aware logging with automatic context propagation
- Support for both development (human-readable) and production (JSON) logging formats
- Automatic source code location (file and line number) tracking
- Log sampling capabilities for high-throughput applications
- OpenTelemetry integration
- Built-in test logging support
- Error logging with automatic stack traces
- Convenient helper functions for common logging patterns
- Group-based logging for better organization
- High-performance logging through Zap backend

## Configuration Options

The logger can be configured using the following options:

### Environment Mode

- `WithDevelopment()`: (Default) Sets readable format with colored level display
- `WithProduction()`: Sets compact JSON format without color

### Log Level and Source

- `WithLevel(level slog.Leveler)`: Sets the minimum logging level
- `WithSource()`: Adds file name and line number to log records

### Identification

- `WithName(name string)`: Sets the logger name

### Sampling

- `WithSampler(tick time.Duration, first, thereafter int)`: Configures log sampling
  - `tick`: Sampling interval
  - `first`: Number of entries to log during the interval
  - `thereafter`: Number of entries to log after the initial entries

### Integration

- `WithOtelTracing()`: Enables OpenTelemetry integration
- `WithTesting(t testing.TB)`: Configures logger for use in tests

## Installation

```bash
go get github.com/n-r-w/ctxlog
```

## Usage Examples

See [example/main.go](example/main.go) for comprehensive usage examples, including:

- Basic logger setup with production mode and debug level
- Context propagation and extraction
- Structured logging with various field types
- Group-based logging
- Error handling with automatic stack traces
- Working with context-bound loggers

The example demonstrates all major features and provides a practical reference for integrating CTXLog into your application.
