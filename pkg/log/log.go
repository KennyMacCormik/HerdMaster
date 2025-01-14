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
	loggerMutex     sync.Mutex
	singletonLogger loggerConfig = loggerConfig{
		logLevelVar: new(slog.LevelVar),
		level:       defaultLogLevel,
		format:      defaultLogFormat,
	}
)

type LoggerOption func(*loggerConfig)

type loggerConfig struct {
	logLevelVar *slog.LevelVar
	logger      *slog.Logger
	level       string
	format      string
	output      io.Writer
}

// LoggingConfig is a structure ready for viper and validator
type LoggingConfig struct {
	Format string `mapstructure:"log_format" validate:"oneof=text json"`            // LOG_FORMAT. Default text
	Level  string `mapstructure:"log_level" validate:"oneof=debug info warn error"` // LOG_LEVEL. Default info
}

// GetLogger returns current logger.
// If no call of SetLogger was invoked before GetLogger, GetLogger will init logger with default settings.
func GetLogger() (*slog.Logger, error) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	if singletonLogger.logger == nil {
		return initLogger()
	}
	return singletonLogger.logger, nil
}

// ConfigureLogger sets config to logger.
func ConfigureLogger(options ...LoggerOption) (*slog.Logger, error) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	return initLogger(options...)
}

// initLoggerWithLock always returns working *slog.Logger.
// Error indicates issues with supplied conf.
// In any case, there is a fallback mechanism so you can start using logger right away
func initLogger(option ...LoggerOption) (*slog.Logger, error) {
	for _, opt := range option {
		opt(&singletonLogger)
	}
	l, f, err := validateAndNormalizeLoggingConf(singletonLogger.level, singletonLogger.format)
	singletonLogger.logger = newLoggerWithConf(l, f)
	return singletonLogger.logger, err
}

// WithDefault sets the logger to use the default configuration.
func WithDefault() LoggerOption {
	return func(cfg *loggerConfig) {
		cfg.level = defaultLogLevel
		cfg.format = defaultLogFormat
	}
}

// WithConfig allows specifying the level and format for the logger.
func WithConfig(level, format string) LoggerOption {
	return func(cfg *loggerConfig) {
		cfg.level = level
		cfg.format = format
	}
}

// WithOutput allows you to change log output. In case output == nil falls back to os.Stdout.
func WithOutput(output io.Writer) LoggerOption {
	return func(cfg *loggerConfig) {
		if output != nil {
			cfg.output = output
		}
	}
}

// validateAndNormalizeLoggingConf validates the provided level and format.
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

// newLoggerWithConf creates a logger with the specified level and format.
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
