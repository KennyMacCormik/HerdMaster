package log

import (
	"log/slog"
	"os"
)

const defaultLogLevel = 0 // info

var logLevelMap = map[string]slog.Level{
	"debug": -4,
	"info":  0,
	"warn":  4,
	"error": 8,
}

func NewLogger(level, format string) *slog.Logger {
	if validateLoggingConf(level, format) {
		lg := loggerWithConf(level, format)
		return lg
	} else {
		lg := DefaultLogger()
		lg.Info("config validation failed, running logger with default values level=info format=text")
		return lg
	}
}

func DefaultLogger() *slog.Logger {
	var logLevel = new(slog.LevelVar)
	logLevel.Set(defaultLogLevel)

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}

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

func loggerWithConf(level, format string) *slog.Logger {
	var logLevel = new(slog.LevelVar)
	logLevel.Set(logLevelMap[level])

	if format == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}
