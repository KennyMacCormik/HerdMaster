package log

import (
	"bytes"
	"github.com/KennyMacCormik/HerdMaster/pkg/cfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
		_ = r.Close()
	}()

	os.Stdout = w

	f()

	_ = w.Close()

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	return buf.String()
}

func TestGetLogger_SingletonBehavior(t *testing.T) {
	logger1 := GetLogger("info", "json")
	logger2 := GetLogger("debug", "text")

	assert.Equal(t, logger1, logger2, "expected logger1 and logger2 to be the same instance")
}

func TestNewLogger_ValidConfiguration(t *testing.T) {
	output := captureOutput(func() {
		logger := newLogger("debug", "text")
		logger.Debug("Test debug log")
	})

	assert.NotEmpty(t, output, "expected debug log output, got empty output")
}

func TestNewLogger_InvalidConfiguration(t *testing.T) {
	const fallbackLogMessage = "Fallback logger log"

	output := captureOutput(func() {
		logger := newLogger("invalid", "invalid")
		logger.Info(fallbackLogMessage)
	})

	assert.NotEmpty(t, output, "expected fallback log output, got empty output")
	assert.Contains(t, output, fallbackLogMessage, "expected fallback log message to be in the output")
}

func TestDefaultLogger(t *testing.T) {
	const defaultLogMessage = "Default logger log"

	output := captureOutput(func() {
		logger := DefaultLogger()
		logger.Info(defaultLogMessage)
	})

	assert.NotEmpty(t, output, "expected default logger output, got empty output")
	assert.Contains(t, output, defaultLogMessage, "expected default log message to be in the output")
}

func TestValidateLoggingConf_Valid(t *testing.T) {
	validCases := []struct {
		level  string
		format string
	}{
		{"debug", "text"},
		{"info", "json"},
		{"warn", "text"},
		{"error", "json"},
	}

	for _, tc := range validCases {
		assert.True(t, validateLoggingConf(tc.level, tc.format), "expected valid configuration for level=%s, format=%s", tc.level, tc.format)
	}
}

func TestValidateLoggingConf_Invalid(t *testing.T) {
	invalidCases := []struct {
		level  string
		format string
	}{
		{"invalid", "json"},
		{"info", "invalid"},
		{"", "text"},
		{"debug", ""},
	}

	for _, tc := range invalidCases {
		assert.False(t, validateLoggingConf(tc.level, tc.format), "expected invalid configuration for level=%s, format=%s", tc.level, tc.format)
	}
}

func TestNewLoggerWithConf_Text(t *testing.T) {
	const testLogMessage = "Test info log"

	output := captureOutput(func() {
		logger := newLoggerWithConf("info", "text")
		logger.Info(testLogMessage)
	})

	require.NotEmpty(t, output, "expected info log output in text format, got empty output")
	assert.Contains(t, output, testLogMessage, "expected log message to be in the output")
}

func TestNewLoggerWithConf_JSON(t *testing.T) {
	const testLogMessage = "Test info log"

	output := captureOutput(func() {
		logger := newLoggerWithConf("info", "json")
		logger.Info(testLogMessage)
	})

	require.NotEmpty(t, output, "expected info log output in JSON format, got empty output")
	assert.Contains(t, output, `"msg":"Test info log"`, "expected JSON-formatted log message to be in the output")
}

func TestLoggingConfigIntegration_DefaultValues(t *testing.T) {
	entry := cfg.ConfigEntry{
		Config: &LoggingConfig{},
		BindArray: []cfg.BindValue{
			{ValName: "log_format", DefaultVal: "text"},
			{ValName: "log_level", DefaultVal: "info"},
		},
	}

	err := cfg.RegisterConfig("log", entry)
	assert.NoError(t, err, "expected no error when registering logging configuration")

	err = cfg.NewConfig()
	assert.NoError(t, err, "expected no error when initializing configuration")

	config, ok := cfg.GetConfig("log")
	assert.True(t, ok, "expected to find registered logging configuration")

	typedConfig, ok := config.(*LoggingConfig)
	assert.True(t, ok, "expected config to be of type LoggingConfig")
	assert.Equal(t, "text", typedConfig.Format, "expected default format to be 'text'")
	assert.Equal(t, "info", typedConfig.Level, "expected default level to be 'info'")
}

func TestLoggingConfigIntegration_EnvironmentOverrides(t *testing.T) {
	entry := cfg.ConfigEntry{
		Config: &LoggingConfig{},
		BindArray: []cfg.BindValue{
			{ValName: "log_format", DefaultVal: "text"},
			{ValName: "log_level", DefaultVal: "info"},
		},
	}

	err := cfg.RegisterConfig("log", entry)
	assert.NoError(t, err, "expected no error when registering logging configuration")

	t.Setenv("HM_LOG_FORMAT", "json")
	t.Setenv("HM_LOG_LEVEL", "debug")

	err = cfg.NewConfig()
	assert.NoError(t, err, "expected no error when initializing configuration")

	config, ok := cfg.GetConfig("log")
	assert.True(t, ok, "expected to find registered logging configuration")

	typedConfig, ok := config.(*LoggingConfig)
	assert.True(t, ok, "expected config to be of type LoggingConfig")
	assert.Equal(t, "json", typedConfig.Format, "expected format to be overridden to 'json'")
	assert.Equal(t, "debug", typedConfig.Level, "expected level to be overridden to 'debug'")
}

func TestLoggingConfigIntegration_SingletonBehavior(t *testing.T) {
	entry := cfg.ConfigEntry{
		Config: &LoggingConfig{},
		BindArray: []cfg.BindValue{
			{ValName: "log_format", DefaultVal: "text"},
			{ValName: "log_level", DefaultVal: "info"},
		},
	}

	err := cfg.RegisterConfig("log", entry)
	assert.NoError(t, err, "expected no error when registering logging configuration")

	err = cfg.NewConfig()
	assert.NoError(t, err, "expected no error when initializing configuration")

	config, ok := cfg.GetConfig("log")
	assert.True(t, ok, "expected to find registered logging configuration")

	typedConfig, ok := config.(*LoggingConfig)
	assert.True(t, ok, "expected config to be of type LoggingConfig")

	logger1 := GetLogger(typedConfig.Level, typedConfig.Format)
	logger2 := GetLogger(typedConfig.Level, typedConfig.Format)

	assert.Equal(t, logger1, logger2, "expected singleton logger instance")
}
