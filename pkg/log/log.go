package log

import (
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
	logLevelVar     = new(slog.LevelVar)
	singletonLogger *slog.Logger
	loggerOnce      sync.Once
)

// GetLogger provides a singleton logger instance with specified level and format.
func GetLogger(level, format string) *slog.Logger {
	loggerOnce.Do(func() {
		singletonLogger = newLogger(level, format)
	})
	return singletonLogger
}

// newLogger creates a new logger with the specified level and format.
// Falls back to the default logger if validation fails.
func newLogger(level, format string) *slog.Logger {
	if validateLoggingConf(level, format) {
		lg := newLoggerWithConf(level, format)
		return lg
	} else {
		lg := defaultLogger()
		lg.Info("config validation failed, running logger with default values", "level", defaultLogLevel, "format", defaultLogFormat)
		return lg
	}
}

// defaultLogger returns a logger with the default configuration (info level, JSON format).
func defaultLogger() *slog.Logger {
	return newLoggerWithConf(defaultLogLevel, defaultLogFormat)
}

// validateLoggingConf validates the provided level and format.
func validateLoggingConf(level, format string) bool {
	if level != "debug" &&
		level != "info" &&
		level != "warn" &&
		level != "error" {
		return false
	}
	if format != "text" &&
		format != "json" {
		return false
	}
	return true
}

// newLoggerWithConf creates a logger with the specified level and format.
func newLoggerWithConf(level, format string) *slog.Logger {
	logLevelVar.Set(logLevelMap[level])

	if format == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevelVar}))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevelVar}))
}
