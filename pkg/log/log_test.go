package log

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"os"
	"testing"
)

func resetLoggerConf() {
	singletonLogger = loggerConfig{
		logLevelVar: new(slog.LevelVar),
		level:       defaultLogLevel,
		format:      defaultLogFormat,
		logger:      nil,
		output:      nil,
	}
}

func TestGetLogger_DefaultConfiguration(t *testing.T) {
	defer resetLoggerConf()
	logger, err := GetLogger()
	require.NoError(t, err, "expected no error for default logger initialization")
	assert.NotNil(t, logger, "expected logger instance to be non-nil")

	output := &bytes.Buffer{}
	logger, err = ConfigureLogger(WithOutput(output))
	require.NoError(t, err, "expected no error for configuring logger output")

	logger.Info("Test message")
	assert.Contains(t, output.String(), "Test message", "expected log message to be written to the output")
}

func TestConfigureLogger_WithValidConfig(t *testing.T) {
	defer resetLoggerConf()
	output := &bytes.Buffer{}
	logger, err := ConfigureLogger(WithConfig("debug", "text"), WithOutput(output))
	require.NoError(t, err, "expected no error for valid logger configuration")
	assert.NotNil(t, logger, "expected logger instance to be non-nil")

	logger.Debug("Debug message")
	assert.Contains(t, output.String(), "Debug message", "expected debug log message to be written to the output")
}

func TestConfigureLogger_WithInvalidLevel(t *testing.T) {
	defer resetLoggerConf()
	output := &bytes.Buffer{}
	logger, err := ConfigureLogger(WithConfig("invalid", "json"), WithOutput(output))
	require.Error(t, err, "expected an error for invalid log level")
	assert.Contains(t, err.Error(), "invalid log level", "expected error message to indicate invalid log level")

	logger.Info("Fallback message")
	assert.Contains(t, output.String(), "Fallback message", "expected fallback log message to be written to the output")
}

func TestConfigureLogger_WithInvalidFormat(t *testing.T) {
	defer resetLoggerConf()
	output := &bytes.Buffer{}
	logger, err := ConfigureLogger(WithConfig("info", "invalid"), WithOutput(output))
	require.Error(t, err, "expected an error for invalid log format")
	assert.Contains(t, err.Error(), "invalid log format", "expected error message to indicate invalid log format")

	logger.Info("Fallback message")
	assert.Contains(t, output.String(), "Fallback message", "expected fallback log message to be written to the output")
}

func TestWithOutput_NilFallbackToStdout(t *testing.T) {
	defer resetLoggerConf()
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	defer func() {
		os.Stdout = oldStdout
		r.Close()
		w.Close()
	}()

	os.Stdout = w

	logger, err := ConfigureLogger(WithOutput(nil))
	require.NoError(t, err, "expected no error for nil output configuration")
	assert.NotNil(t, logger, "expected logger instance to be non-nil")

	logger.Info("Stdout message")
	w.Close()

	out := &bytes.Buffer{}
	io.Copy(out, r)
	assert.Contains(t, out.String(), "Stdout message", "expected log message to be written to stdout")
}

func TestValidateAndNormalizeLoggingConf(t *testing.T) {
	defer resetLoggerConf()
	tests := []struct {
		name           string
		level          string
		format         string
		expectedLevel  string
		expectedFormat string
		expectedError  string
	}{
		{"ValidConfig", "info", "json", "info", "json", ""},
		{"InvalidLevel", "invalid", "json", "info", "json", "invalid log level"},
		{"InvalidFormat", "info", "invalid", "info", "json", "invalid log format"},
		{"InvalidBoth", "invalid", "invalid", "info", "json", "invalid log level invalid\ninvalid log format invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, format, err := validateAndNormalizeLoggingConf(tt.level, tt.format)
			assert.Equal(t, tt.expectedLevel, level, "expected normalized level")
			assert.Equal(t, tt.expectedFormat, format, "expected normalized format")

			if tt.expectedError != "" {
				require.Error(t, err, "expected an error")
				assert.Contains(t, err.Error(), tt.expectedError, "expected specific error message")
			} else {
				assert.NoError(t, err, "expected no error")
			}
		})
	}
}

func TestInitLoggerWithMultipleOptions(t *testing.T) {
	defer resetLoggerConf()
	output := &bytes.Buffer{}
	logger, err := ConfigureLogger(WithConfig("debug", "json"), WithOutput(output))
	require.NoError(t, err, "expected no error for configuring logger with multiple options")

	logger.Debug("Debug message")
	assert.Contains(t, output.String(), "Debug message", "expected debug log message in JSON format")
	assert.Contains(t, output.String(), `"level":"DEBUG"`, "expected debug level in JSON format")
}

func TestConfigureLogger_MultipleInstances(t *testing.T) {
	defer resetLoggerConf()
	output1 := &bytes.Buffer{}
	output2 := &bytes.Buffer{}

	logger1, err := ConfigureLogger(WithConfig("info", "json"), WithOutput(output1))
	require.NoError(t, err, "expected no error for first logger configuration")
	logger1.Info("First logger message")

	logger2, err := ConfigureLogger(WithConfig("debug", "text"), WithOutput(output2))
	require.NoError(t, err, "expected no error for second logger configuration")
	logger2.Debug("Second logger message")

	assert.Contains(t, output1.String(), "First logger message", "expected message from the first logger")
	assert.Contains(t, output2.String(), "Second logger message", "expected message from the second logger")
}
