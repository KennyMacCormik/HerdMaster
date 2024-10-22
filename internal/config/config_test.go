package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

type testCase struct {
	env map[string]string
	err string
}

func setEnv(list map[string]string) {
	for k, v := range list {
		if v != "" {
			_ = os.Setenv(k, v)
		}
	}
}

func unsetEnv(list map[string]string) {
	for k, v := range list {
		if v != "" {
			_ = os.Unsetenv(k)
		}
	}
}

func TestConfig_Positive_AllPresent(t *testing.T) {
	test := testCase{
		env: map[string]string{
			// type Network struct
			"HM_NET_HOST":     "127.0.0.1",
			"HM_NET_PORT":     "8081",
			"HM_NET_MAX_CONN": "101",
			"HM_NET_TIMEOUT":  "1m",
			// type Logging struct
			"HM_LOG_FORMAT": "json",
			"HM_LOG_LEVEL":  "warn",
		},
	}

	setEnv(test.env)
	defer unsetEnv(test.env)

	conf, err := New()
	assert.NoError(t, err)
	// NET
	assert.Equal(t, test.env["HM_NET_HOST"], conf.Net.Host)

	i, err := strconv.Atoi(os.Getenv("HM_NET_PORT"))
	assert.NoError(t, err)
	assert.Equal(t, i, conf.Net.Port)

	i, err = strconv.Atoi(os.Getenv("HM_NET_MAX_CONN"))
	assert.NoError(t, err)
	assert.Equal(t, i, conf.Net.MaxConn)

	d, err := time.ParseDuration(test.env["HM_NET_TIMEOUT"])
	assert.NoError(t, err)
	assert.Equal(t, d, conf.Net.Timeout)
	// LOG
	assert.Equal(t, test.env["HM_LOG_FORMAT"], conf.Log.Format)
	assert.Equal(t, test.env["HM_LOG_LEVEL"], conf.Log.Level)
}

func TestConfig_Positive_AllDefault(t *testing.T) {
	test := testCase{
		env: map[string]string{
			// type Network struct
			"HM_NET_HOST":     "",
			"HM_NET_PORT":     "",
			"HM_NET_MAX_CONN": "",
			"HM_NET_TIMEOUT":  "",
			// type Logging struct
			"HM_LOG_FORMAT": "",
			"HM_LOG_LEVEL":  "",
		},
	}

	setEnv(test.env)
	defer unsetEnv(test.env)

	conf, err := New()
	assert.NoError(t, err)
	// NET
	assert.Equal(t, "0.0.0.0", conf.Net.Host)
	assert.Equal(t, 8080, conf.Net.Port)
	assert.Equal(t, runtime.NumCPU(), conf.Net.MaxConn)
	assert.Equal(t, 1*time.Second, conf.Net.Timeout)
	// LOG
	assert.Equal(t, "text", conf.Log.Format)
	assert.Equal(t, "info", conf.Log.Level)
}

// type Network struct
func TestConfig_Positive_HM_NET_HOST_AllValid_IP(t *testing.T) {

	tests := []testCase{
		{
			env: map[string]string{
				"HM_NET_HOST": "0.0.0.0",
			},
		},
		{
			env: map[string]string{
				"HM_NET_HOST": "127.0.0.1",
			},
		},
		{
			env: map[string]string{
				"HM_NET_HOST": "192.168.12.2",
			},
		},
		{
			env: map[string]string{
				"HM_NET_HOST": "hm.local",
			},
		},
		{
			env: map[string]string{
				"HM_NET_HOST": "server",
			},
		},
	}

	for _, test := range tests {
		setEnv(test.env)
		conf, err := New()
		assert.NoError(t, err)
		assert.Equal(t, test.env["HM_NET_HOST"], conf.Net.Host)
		unsetEnv(test.env)
	}
}

// Logging
func TestConfig_Positive_HM_LOG_LEVEL_AllValid(t *testing.T) {
	tests := []testCase{
		{
			env: map[string]string{
				"HM_LOG_LEVEL": "debug",
			},
		},
		{
			env: map[string]string{
				"HM_LOG_LEVEL": "error",
			},
		},
		{
			env: map[string]string{
				"HM_LOG_LEVEL": "warn",
			},
		},
		{
			env: map[string]string{
				"HM_LOG_LEVEL": "info",
			},
		},
	}

	for _, test := range tests {
		setEnv(test.env)
		conf, err := New()
		assert.NoError(t, err)
		assert.Equal(t, test.env["HM_LOG_LEVEL"], conf.Log.Level)
		unsetEnv(test.env)
	}
}

func TestConfig_Positive_HM_LOG_FORMAT_AllValid(t *testing.T) {
	tests := []testCase{
		{
			env: map[string]string{
				"HM_LOG_FORMAT": "text",
			},
		},
		{
			env: map[string]string{
				"HM_LOG_FORMAT": "json",
			},
		},
	}

	for _, test := range tests {
		setEnv(test.env)
		conf, err := New()
		assert.NoError(t, err)
		assert.Equal(t, test.env["HM_LOG_FORMAT"], conf.Log.Format)
		unsetEnv(test.env)
	}
}

func TestConfig_Negative_HM_LOG_FORMAT_BogusArgs(t *testing.T) {
	tests := []testCase{
		{
			env: map[string]string{
				"HM_LOG_FORMAT": "textt",
			},
			err: "config validation error: env 'HM_LOG_FORMAT' value 'textt' invalid, 'oneof=text json' expected;",
		},
		{
			env: map[string]string{
				"HM_LOG_FORMAT": "123",
			},
			err: "config validation error: env 'HM_LOG_FORMAT' value '123' invalid, 'oneof=text json' expected;",
		},
		{
			env: map[string]string{
				"HM_LOG_FORMAT": "qwe%rty",
			},
			err: "config validation error: env 'HM_LOG_FORMAT' value 'qwe%rty' invalid, 'oneof=text json' expected;",
		},
		{
			env: map[string]string{
				"HM_LOG_FORMAT": "js0n",
			},
			err: "config validation error: env 'HM_LOG_FORMAT' value 'js0n' invalid, 'oneof=text json' expected;",
		},
		{
			env: map[string]string{
				"HM_LOG_FORMAT": "jsOn",
			},
			err: "config validation error: env 'HM_LOG_FORMAT' value 'jsOn' invalid, 'oneof=text json' expected;",
		},
	}

	for _, test := range tests {
		setEnv(test.env)
		conf, err := New()
		assert.Equal(t, Config{}, conf)
		assert.EqualError(t, err, test.err)
		unsetEnv(test.env)
	}
}

func TestConfig_Negative_HM_LOG_LEVEL_BogusArgs(t *testing.T) {
	tests := []testCase{
		{
			env: map[string]string{
				"HM_LOG_LEVEL": "debugg",
			},
			err: "config validation error: env 'HM_LOG_LEVEL' value 'debugg' invalid, 'oneof=debug info warn error' expected;",
		},
		{
			env: map[string]string{
				"HM_LOG_LEVEL": "123",
			},
			err: "config validation error: env 'HM_LOG_LEVEL' value '123' invalid, 'oneof=debug info warn error' expected;",
		},
		{
			env: map[string]string{
				"HM_LOG_LEVEL": "qwe%rty",
			},
			err: "config validation error: env 'HM_LOG_LEVEL' value 'qwe%rty' invalid, 'oneof=debug info warn error' expected;",
		},
		{
			env: map[string]string{
				"HM_LOG_LEVEL": "inf0",
			},
			err: "config validation error: env 'HM_LOG_LEVEL' value 'inf0' invalid, 'oneof=debug info warn error' expected;",
		},
		{
			env: map[string]string{
				"HM_LOG_LEVEL": "infO",
			},
			err: "config validation error: env 'HM_LOG_LEVEL' value 'infO' invalid, 'oneof=debug info warn error' expected;",
		},
	}

	for _, test := range tests {
		setEnv(test.env)
		conf, err := New()
		assert.Equal(t, Config{}, conf)
		assert.EqualError(t, err, test.err)
		unsetEnv(test.env)
	}
}
