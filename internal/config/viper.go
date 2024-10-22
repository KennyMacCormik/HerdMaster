package config

import (
	"github.com/spf13/viper"
	"runtime"
	"strconv"
)

func setNetworkEnv() {
	viper.SetDefault("net_host", "0.0.0.0")
	_ = viper.BindEnv("net_host")

	viper.SetDefault("net_port", "8080")
	_ = viper.BindEnv("net_port")

	viper.SetDefault("net_max_conn", strconv.Itoa(runtime.NumCPU()))
	_ = viper.BindEnv("net_max_conn")

	viper.SetDefault("net_timeout", "1s")
	_ = viper.BindEnv("net_timeout")
}

func setLoggingEnv() {
	viper.SetDefault("log_format", "text")
	_ = viper.BindEnv("format")

	viper.SetDefault("log_level", "info")
	_ = viper.BindEnv("level")
}

func loadEnv(c *Config) error {
	viper.SetEnvPrefix("HM")

	setNetworkEnv()
	setLoggingEnv()

	viper.AutomaticEnv()
	return viper.Unmarshal(c)
}
