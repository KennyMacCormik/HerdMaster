// Package log provides a simple, thread-safe singleton logger based on the Go slog package.
// It allows configuring the logger with customizable logging levels, formats, and output destinations.
// The logger supports default configurations and enables users to override them with options.
package log

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
)

const (
	defaultLogLevel  = "info"
	defaultLogFormat = "json"
)

var (
	logLevelMap = map[string]slog.Level{
		"debug": -4,
		"info":  0,
		"warn":  4,
		"error": 8,
	}
	loggerMutex sync.Mutex
	// singletonLogger holds the global logger configuration and state.
	singletonLogger = loggerConfig{
		logLevelVar: new(slog.LevelVar),
		level:       defaultLogLevel,
		format:      defaultLogFormat,
	}
)

// LoggerOption represents a functional option for configuring the logger.
type LoggerOption func(*loggerConfig)

// loggerConfig holds the internal state of the singleton logger, including configuration values
// such as logging level, format, output destination, and the slog.Logger instance.
type loggerConfig struct {
	logLevelVar *slog.LevelVar
	logger      *slog.Logger
	level       string
	format      string
	output      io.Writer
}

// LoggingConfig is a structure designed for integration with viper and validation libraries.
// It provides the fields necessary for external configuration.
type LoggingConfig struct {
	Format string `mapstructure:"log_format" validate:"oneof=text json"`            // LOG_FORMAT. Default text
	Level  string `mapstructure:"log_level" validate:"oneof=debug info warn error"` // LOG_LEVEL. Default info
}

// GetLogger returns the current logger instance. If the logger is uninitialized,
// it creates a logger with the default settings.
//
// Note: This method ensures the logger is always usable even without prior configuration.
func GetLogger() (*slog.Logger, error) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if singletonLogger.logger == nil {
		return initLogger()
	}
	return singletonLogger.logger, nil
}

// ConfigureLogger applies the provided options to configure the logger.
// It returns the updated logger instance and any error encountered during the process.
//
// Note: The caller must save the returned *slog.Logger instance to ensure the
// updated logger is used. Example:
//
//	logger, err := log.ConfigureLogger(log.WithConfig("debug", "text"))
//	if err != nil {
//		// Handle error
//	}
func ConfigureLogger(options ...LoggerOption) (*slog.Logger, error) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	return initLogger(options...)
}

// initLogger applies the provided options and initializes the logger configuration.
// It always returns a working *slog.Logger instance.
// Any errors are related to
// the supplied configuration but don't prevent the logger from functioning.
func initLogger(option ...LoggerOption) (*slog.Logger, error) {
	for _, opt := range option {
		opt(&singletonLogger)
	}

	l, f, err := validateAndNormalizeLoggingConf(singletonLogger.level, singletonLogger.format)

	singletonLogger.logger = newLoggerWithConf(l, f)
	return singletonLogger.logger, err
}

// WithDefault sets the logger to use the default configuration (info level, JSON format).
func WithDefault() LoggerOption {
	return func(cfg *loggerConfig) {
		cfg.level = defaultLogLevel
		cfg.format = defaultLogFormat
	}
}

// WithConfig allows specifying the log level and format for the logger.
// Example usage:
//
//	logger, err := log.ConfigureLogger(log.WithConfig("debug", "text"))
func WithConfig(level, format string) LoggerOption {
	return func(cfg *loggerConfig) {
		cfg.level = level
		cfg.format = format
	}
}

// WithOutput allows specifying a custom output destination for the logger.
// If output is nil, it defaults to os.Stdout.
func WithOutput(output io.Writer) LoggerOption {
	return func(cfg *loggerConfig) {
		if output != nil {
			cfg.output = output
		}
	}
}

// validateAndNormalizeLoggingConf validates and normalizes the logging level and format.
// It returns the resulting level, format, and any associated validation errors.
func validateAndNormalizeLoggingConf(level, format string) (string, string, error) {
	var errLevel, errFormat error
	resultLevel, resultFormat := level, format
	if _, ok := logLevelMap[resultLevel]; !ok {
		errLevel = fmt.Errorf("invalid log level %s", level)
		resultLevel = defaultLogLevel
	}
	if format != "text" && format != "json" {
		errFormat = fmt.Errorf("invalid log format %s", format)
		resultFormat = defaultLogFormat
	}
	return resultLevel, resultFormat, errors.Join(errLevel, errFormat)
}

// newLoggerWithConf creates a logger instance with the specified level and format.
// It falls back to os.Stdout if no output destination is specified.
func newLoggerWithConf(level, format string) *slog.Logger {
	output := singletonLogger.output
	if output == nil {
		output = os.Stdout // Default to stdout
	}

	singletonLogger.logLevelVar.Set(logLevelMap[level])

	if format == "text" {
		return slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{Level: singletonLogger.logLevelVar}))
	}

	return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: singletonLogger.logLevelVar}))
}
