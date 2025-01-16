// Package cfg provides a flexible and extensible configuration management system
// that integrates with Viper.
// It allows clients to register their custom configuration structs, bind environment
// variables, and set default values dynamically.
//
// Example usage:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/KennyMacCormik/HerdMaster/pkg/cfg"
//	)
//
//	// LoggingConfig represents the configuration for logging.
//	type LoggingConfig struct {
//		Format string `mapstructure:"log_format" validate:"oneof=text json"`            // HM_LOG_FORMAT. Default text
//		Level  string `mapstructure:"log_level" validate:"oneof=debug info warn error"` // HM_LOG_LEVEL. Default info
//	}
//
//	func main() {
//		// Define a configuration entry for the logging system
//		entry := cfg.ConfigEntry{
//			Config: &LoggingConfig{}, // Struct for logging configuration
//			BindArray: []cfg.BindValue{
//				{
//					ValName: "log_format", // Environment variable name
//				},
//				{
//					ValName:    "log_level", // Environment variable name
//					DefaultVal: "info",      // Default value if environment variable is not set
//				},
//			},
//		}
//
//		// Register the logging configuration
//		if err := cfg.RegisterConfig("log", entry); err != nil {
//			fmt.Printf("Error registering config: %v\n", err)
//			return
//		}
//
//		// Initialize all registered configurations
//		if err := cfg.NewConfig(); err != nil {
//			fmt.Printf("Error initializing configs: %v\n", err)
//			return
//		}
//
//		// Access the initialized configuration
//		logConfig, ok := cfg.GetConfig("log")
//		if !ok {
//			fmt.Println("Failed to retrieve the logging config.")
//			return
//		}
//
//		// Type assert to the specific configuration struct
//		if typedConfig, ok := logConfig.(*LoggingConfig); ok {
//			fmt.Printf("Initialized Logging Config: %+v\n", typedConfig)
//		} else {
//			fmt.Println("Failed to type assert logging config.")
//		}
//	}
package cfg

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"sync"
)

var (
	configEntries = make(map[string]ConfigEntry)
	mtx           sync.RWMutex
)

// ConfigEntry represents a registered configuration entry.
// It includes the configuration struct
// and an array of environment variable bindings.
type ConfigEntry struct {
	Config    any
	BindArray []BindValue
}

// BindValue represents a binding of an environment variable to a configuration field
// with an optional default value.
type BindValue struct {
	ValName    string // Environment variable name
	DefaultVal any    // Default value for the environment variable
}

// ViperOption represents a functional option for configuring the behavior of Viper.
// It allows for customizing Viper's initialization, such as setting environment variable prefixes.
// This enables flexible configuration adjustments without altering the core logic.
type ViperOption func() error

// NewConfig initializes the configuration system, binds all registered configuration entries,
// and loads their values from environment variables.
func NewConfig(list ...ViperOption) error {
	for _, opt := range list {
		err := opt()
		if err != nil {
			return err
		}
	}

	viper.AutomaticEnv()

	return newConfigWithLock()
}

// WithSetEnvPrefix sets the environment variable prefix for Viper.
// It ensures that all environment variables bound to the configuration entries
// will use the specified prefix, enabling namespacing and avoiding conflicts with other environment variables.
// EnvPrefix can't be empty string.
//
// Example Usage:
//
//	err := cfg.NewConfig(cfg.WithSetEnvPrefix("myapp"))
//	if err != nil {
//	    fmt.Printf("Error initializing configs: %v\n", err)
//	}
//
// In this example, Viper will look for environment variables prefixed with "MYAPP_".
func WithSetEnvPrefix(EnvPrefix string) ViperOption {
	if EnvPrefix == "" {
		return func() error {
			return fmt.Errorf("incorrect env prefix: %s", EnvPrefix)
		}
	}
	return func() error {
		viper.SetEnvPrefix(EnvPrefix)
		return nil
	}
}

// newConfigWithLock acquires a write lock and iterates through all registered configuration entries
// to bind and load their values from the environment.
func newConfigWithLock() error {
	mtx.Lock()
	defer mtx.Unlock()
	for k, v := range configEntries {
		if err := bindActualValue(&v); err != nil {
			return fmt.Errorf("failed to bind config %s: %w", k, err)
		}
	}
	return nil
}

// RegisterConfig allows clients to register their custom configuration structs
// along with their environment variable bindings.
func RegisterConfig(name string, configStruct ConfigEntry) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if err := validateConfigStruct(configStruct); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	if err := validateBindArray(configStruct.BindArray); err != nil {
		return fmt.Errorf("bind array validation failed: %w", err)
	}

	storeConfigStructWithLock(name, configStruct)
	return nil
}

// ListConfigs returns a list of all registered configuration names.
func ListConfigs() []string {
	mtx.RLock()
	defer mtx.RUnlock()
	keys := make([]string, 0, len(configEntries))
	for k := range configEntries {
		keys = append(keys, k)
	}
	return keys
}

// GetConfig retrieves a registered configuration struct by name.
func GetConfig(name string) (any, bool) {
	v, ok := getConfigWithRLock(name)
	return v.Config, ok
}

// storeConfigStructWithLock safely stores a configuration struct in the registry with a write lock.
func storeConfigStructWithLock(name string, configStruct ConfigEntry) {
	mtx.Lock()
	defer mtx.Unlock()
	configEntries[name] = configStruct
}

// validateBindArray ensures that all BindValues in the array have non-empty ValName fields.
func validateBindArray(bindArray []BindValue) error {
	for _, bind := range bindArray {
		if bind.ValName == "" {
			return fmt.Errorf("BindValue.ValName cannot be empty")
		}
	}
	return nil
}

// getConfigWithRLock retrieves a registered configuration struct by name with a read lock.
func getConfigWithRLock(name string) (ConfigEntry, bool) {
	mtx.RLock()
	defer mtx.RUnlock()
	v, ok := configEntries[name]
	return v, ok
}

// isNotNullOrDefault checks if a value is not nil and not its default value.
func isNotNullOrDefault(value any) bool {
	// Check if the value is nil
	if value == nil {
		return false
	}

	val := reflect.ValueOf(value)

	// Check for the default value
	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return false
		}
		val = val.Elem()
	}

	// Compare with zero value
	zero := reflect.Zero(val.Type()).Interface()
	return !reflect.DeepEqual(val.Interface(), zero)
}

// bindToEnv binds environment variables and sets default values for a configuration entry.
func bindToEnv(entry *ConfigEntry) error {
	for _, e := range entry.BindArray {
		if err := viper.BindEnv(e.ValName); err != nil {
			return fmt.Errorf("failed to bind %s: %w", e.ValName, err)
		}
		if isNotNullOrDefault(e.DefaultVal) {
			viper.SetDefault(e.ValName, e.DefaultVal)
		}
	}
	return nil
}

// validateConfigStruct validates that the Config field of a ConfigEntry is a pointer to a struct.
func validateConfigStruct(config ConfigEntry) error {
	configValue := reflect.ValueOf(config.Config)
	if configValue.Kind() != reflect.Ptr || configValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to a struct, got %s", configValue.Kind())
	}
	return nil
}

// bindActualValue binds environment variables to a configuration entry and unmarshals
// the environment values into the Config field.
func bindActualValue(entry *ConfigEntry) error {
	if err := bindToEnv(entry); err != nil {
		return err
	}

	if err := viper.Unmarshal(entry.Config); err != nil {
		return fmt.Errorf("failed to unmarshal into config: %w", err)
	}

	return nil
}
