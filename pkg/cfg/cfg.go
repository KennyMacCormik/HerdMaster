package cfg

import (
	"github.com/KennyMacCormik/HerdMaster/pkg/val"
	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvPrefix("hm")
	setLoggingEnv()
	setGrpcEnv()
	viper.AutomaticEnv()
}

type DefaultConfig struct {
	Log  LoggingConfig `mapstructure:",squash"`
	Grpc GrpcConfig    `mapstructure:",squash"`
}

type LoggingConfig struct {
	// HM_LOG_FORMAT. Default text
	Format string `mapstructure:"log_format" validate:"oneof=text json"`
	// HM_LOG_LEVEL. Default info
	Level string `mapstructure:"log_level" validate:"oneof=debug info warn error"`
}

type GrpcConfig struct {
	// HM_GRPC_HOST. Default text
	Host string `mapstructure:"grpc_host" validate:"ip4_addr|hostname_rfc1123"`
	// HM_GRPC_PORT. Default info
	Port string `mapstructure:"grpc_port" validate:"numeric,gt=1024,lt=65536"`
}

func setGrpcEnv() {
	viper.SetDefault("grpc_host", "0.0.0.0")
	_ = viper.BindEnv("grpc_host")

	viper.SetDefault("grpc_port", "8080")
	_ = viper.BindEnv("grpc_port")
}

func setLoggingEnv() {
	viper.SetDefault("log_format", "text")
	_ = viper.BindEnv("log_format")

	viper.SetDefault("log_level", "info")
	_ = viper.BindEnv("log_level")
}

func NewDefaultConfig(conf *DefaultConfig, val *val.GlobalValidator) error {
	if err := viper.Unmarshal(conf); err != nil {
		return err
	}

	if err := val.ValidateStruct(*conf); err != nil {
		return err
	}

	return nil
}
