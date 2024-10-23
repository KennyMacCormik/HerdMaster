package config

import (
	"errors"
	"fmt"
	"time"
)

type Network struct {
	// Address to listen. Default 0.0.0.0
	Host string `mapstructure:"net_host" validate:"ip4_addr|hostname_rfc1123" env:"HM_NET_HOST"`
	// Port to listen. Default 8080
	Port int `mapstructure:"net_port" validate:"numeric,gt=1024,lt=65536" env:"HM_NET_PORT"`
	// Maximum simultaneously working connections. Default runtime.NumCPU()
	MaxConn int `mapstructure:"net_max_conn" validate:"numeric,gte=0" env:"HM_NET_MAX_CONN"`
	// Idle connection timeout. Min 100 ms. Default 1 s
	Timeout time.Duration `mapstructure:"net_timeout" validate:"min=100ms" env:"HM_NET_TIMEOUT"`
}

type Logging struct {
	// Log format. Default text
	Format string `mapstructure:"log_format" validate:"oneof=text json" env:"HM_LOG_FORMAT"`
	// Log level. Default info
	Level string `mapstructure:"log_level" validate:"oneof=debug info warn error" env:"HM_LOG_LEVEL"`
}

type Database struct {
	// DB uri. Example postgresql://user:secret@host:5432/repos
	ConnString string `mapstructure:"db_uri" validate:"required,uri" env:"HM_DB_URI"`
	// Whether apply auto migrations or not. Default true
	AutoMigrate bool `mapstructure:"auto_migrate" validate:"boolean" env:"HM_DB_AUTO_MIGRATE"`
	// Whether fill dictionary tables or not. Default true
	AutoFillDict bool `mapstructure:"auto_fill_dict" validate:"boolean" env:"HM_DB_AUTO_FILL_DICT"`
}

type Config struct {
	Net Network  `mapstructure:",squash"`
	Log Logging  `mapstructure:",squash"`
	DB  Database `mapstructure:",squash"`
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
