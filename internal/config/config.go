package config

import (
	"errors"
	"fmt"
	"time"
)

type Network struct {
	// Env HM_NET_HOST. Address to listen. Default 0.0.0.0
	Host string `mapstructure:"net_host" validate:"ip4_addr"`
	// Env HM_NET_PORT. Port to listen. Default 8080
	Port int `mapstructure:"net_port" validate:"numeric,gt=1024,lt=65536"`
	// Env HM_NET_MAX_CONN. Maximum simultaneously working connections. Default runtime.NumCPU()
	MaxConn int `mapstructure:"net_max_conn" validate:"numeric,gte=0"`
	// Env HM_NET_TIMEOUT. Idle connection timeout. Min 100 ms. Default 1 s
	Timeout time.Duration `mapstructure:"net_timeout" validate:"min=100ms"`
}

type Logging struct {
	// Env HM_LOG_FORMAT. Log format. Default text
	Format string `mapstructure:"log_format" validate:"oneof=text json"`
	// Env HM_LOG_LEVEL. Log level. Default info
	Level string `mapstructure:"log_level" validate:"oneof=debug info warn error"`
}

type Config struct {
	Net Network `mapstructure:",squash"`
	Log Logging `mapstructure:",squash"`
}

func New() (Config, error) {
	c := Config{}

	err := loadEnv(&c)
	if err != nil {
		return Config{}, fmt.Errorf("config unmarshalling error: %w", err)
	}

	err = validate(c)
	if err != nil {
		return Config{}, fmt.Errorf("config validation error: %w", errors.New(handleValidatorError(c, err)))
	}

	return c, nil
}
