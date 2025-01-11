package cfg

import (
	"github.com/KennyMacCormik/HerdMaster/pkg/val"
	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvPrefix("hm")
	setLoggingEnv()

	viper.AutomaticEnv()
}

type DefaultConfig struct {
	Log LoggingConfig `mapstructure:",squash"`
}

type LoggingConfig struct {
	// HM_LOG_FORMAT. Default text
	Format string `mapstructure:"log_format" validate:"oneof=text json"`
	// HM_LOG_LEVEL. Default info
	Level string `mapstructure:"log_level" validate:"oneof=debug info warn error"`
}

func setLoggingEnv() {
	viper.SetDefault("log_format", "text")
	_ = viper.BindEnv("log_format")

	viper.SetDefault("log_level", "info")
	_ = viper.BindEnv("log_level")
}

func NewDefaultConfig(conf *DefaultConfig) error {
	if err := viper.Unmarshal(conf); err != nil {
		return err
	}

	if err := val.ValInstance.ValidateStruct(*conf); err != nil {
		return err
	}

	return nil
}

func NewCustomConfig(conf any) error {
	if err := viper.Unmarshal(conf); err != nil {
		return err
	}

	if err := val.ValInstance.ValidateStruct(&conf); err != nil {
		return err
	}

	return nil
}
