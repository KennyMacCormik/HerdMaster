
# Log Package

## Overview

The `log` package provides a simple, thread-safe singleton logger built on Go's `log/slog` package. It supports configurable logging levels, formats, and output destinations, ensuring flexibility for various use cases. By default, the logger is initialized with sensible defaults but allows customization through functional options.

## Features

- Thread-safe, singleton logger implementation.
- Customizable logging levels (`debug`, `info`, `warn`, `error`).
- Configurable logging formats (`json`, `text`).
- Support for custom output destinations.
- Default configurations ensuring a usable logger without prior setup.

## Installation

### Dependencies Installation

This package relies on Go's standard library. Ensure you are using **Go 1.23** or later.

### Package Installation

Install the package using the following command:

```sh
go get github.com/KennyMacCormik/HerdMaster/pkg/log
```

## Usage

### Example: Default Logger

```go
package main

import (
	"log/slog"
	"github.com/KennyMacCormik/HerdMaster/pkg/log"
)

func main() {
	logger, err := log.GetLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("This is a default logger message")
}
```

### Example: Configuring Logger

```go
package main

import (
	"bytes"
	"log/slog"
	"github.com/KennyMacCormik/HerdMaster/pkg/log"
)

func main() {
	output := &bytes.Buffer{}
	logger, err := log.ConfigureLogger(
		log.WithConfig("debug", "text"),
		log.WithOutput(output),
	)
	if err != nil {
		panic(err)
	}
	logger.Debug("This is a debug message")
}
```

### Note

The logger returned by `ConfigureLogger` **must be saved** and used for logging. This behavior is dictated by the immutability of `*slog.Logger` in the `slog` package.

## API Documentation

### Exported Functions

#### `GetLogger()`

```go
func GetLogger() (*slog.Logger, error)
```

Returns the current logger instance. If the logger is uninitialized, it creates one with the default settings.

#### `ConfigureLogger(options ...LoggerOption) (*slog.Logger, error)`

Configures the logger with the provided options. Returns the updated logger instance.

### Logger Configuration Options

#### `WithDefault()`

```go
func WithDefault() LoggerOption
```

Sets the logger to use the default configuration (`info` level, `json` format).

#### `WithConfig(level, format string)`

```go
func WithConfig(level, format string) LoggerOption
```

Allows specifying the logging level and format.

#### `WithOutput(output io.Writer)`

```go
func WithOutput(output io.Writer) LoggerOption
```

Sets a custom output destination for the logger. Defaults to `os.Stdout` if `output` is `nil`.

## Type Description

### `LoggerOption`

Represents a functional option for configuring the logger.

### `loggerConfig`

Internal configuration structure for the logger. Not exportable.

### `LoggingConfig`

```go
type LoggingConfig struct {
	Format string `mapstructure:"log_format" validate:"oneof=text json"`
	Level  string `mapstructure:"log_level" validate:"oneof=debug info warn error"`
}
```

A structure designed for integration with configuration tools like Viper and validation libraries.

## License

This project is licensed under the MIT License. See the [LICENSE](https://opensource.org/licenses/MIT) for details.

## Thanks

Special thanks to the Go team for the `log/slog` package. For more details, visit:
- [Go Documentation](https://pkg.go.dev/log/slog)
