
# `log` Package Documentation

The `log` package provides a flexible, singleton-based logging system using the `log/slog` package. It supports configurable log levels and formats, with sensible defaults for fallback.

## Features

1. **Singleton Logger**: Ensures a single, globally accessible logger instance throughout the application.
2. **Configurable Levels and Formats**:
    - Supported levels: `debug`, `info`, `warn`, `error`
    - Supported formats: `text`, `json`
3. **Default Fallback**: Provides a default logger with `info` level and `json` format when configuration validation fails.
4. **Validation**: Validates the provided level and format to prevent misconfiguration.

---

## Usage

### GetLogger

```go
func GetLogger(level, format string) *slog.Logger
```

Returns a singleton logger instance with the specified log level and format. If called multiple times, the same instance is returned.

#### Example:

```go
logger := log.GetLogger("debug", "json")
logger.Debug("Debug message")
```

---

### DefaultLogger

```go
func defaultLogger() *slog.Logger
```

Returns a logger with the default configuration (`info` level, `json` format).

#### Example:

```go
logger := defaultLogger()
logger.Info("Default logger message")
```

---

### Configuration Validation

#### validateLoggingConf

```go
func validateLoggingConf(level, format string) bool
```

Validates the specified log level and format. Returns `true` if valid, `false` otherwise.

#### Example:

```go
valid := validateLoggingConf("info", "json")
fmt.Println(valid) // Output: true
```

---

### Logger Creation

#### newLogger

```go
func newLogger(level, format string) *slog.Logger
```

Creates a logger with the specified configuration. Falls back to the default logger if validation fails.

#### Example:

```go
logger := newLogger("warn", "text")
logger.Warn("Warning message")
```

#### newLoggerWithConf

```go
func newLoggerWithConf(level, format string) *slog.Logger
```

Creates a logger with validated configuration for the specified level and format.

---

## Supported Log Levels

- `debug`
- `info`
- `warn`
- `error`

---

## Supported Formats

- `text`: Outputs logs in human-readable text format.
- `json`: Outputs logs in JSON format.

---

## Example Usage

```go
package main

import (
	"log"
)

func main() {
	logger := log.GetLogger("info", "json")
	logger.Info("Application started")
}
```

---

## Unit Tests

The `log` package includes comprehensive unit tests to validate:
1. Singleton behavior (`GetLogger`)
2. Valid and invalid configurations (`validateLoggingConf`)
3. Log output formats (`text`, `json`)
4. Default fallback behavior

Run the tests with:

```bash
go test ./... -v
```

---

## Limitations

1. The logger configuration cannot be updated after the first call to `GetLogger`.
2. The package does not support log rotation or external logging services (e.g., Elasticsearch).

---
