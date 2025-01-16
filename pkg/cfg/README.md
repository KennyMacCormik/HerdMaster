
# README.md

## Overview

The `cfg` package provides a robust and extensible configuration management system that seamlessly integrates with Viper.
This package allows applications to register and manage custom configuration structs, bind them to environment variables,
and set dynamic defaults. It is part of the HerdMaster repository and ensures centralized and reusable configuration management
across all microservices.

## Features

- **Custom Configuration Registration**: Register custom structs with specific environment variable bindings.
- **Dynamic Defaults**: Provides default values when environment variables are unset.
- **Thread-Safe Design**: Ensures safe concurrent access and updates to configurations.
- **Functional Options**: Customize Viper’s initialization via functional options.
- **Validation Ready**: Ensures all configurations adhere to validation rules using the `val` package.

## Installation

### Dependencies Installation

Ensure that the Viper library is installed:

```bash
go get github.com/spf13/viper
```

### Package Installation

Install the `cfg` package:

```bash
go get github.com/KennyMacCormik/HerdMaster/pkg/cfg
```

## Usage

Below is an example demonstrating the use of the `cfg` package:

```go
package main

import (
    "fmt"
    "github.com/KennyMacCormik/HerdMaster/pkg/cfg"
)

// LoggingConfig represents the configuration for logging.
type LoggingConfig struct {
    Format string `mapstructure:"log_format" validate:"oneof=text json"`            // LOG_FORMAT. Default text
    Level  string `mapstructure:"log_level" validate:"oneof=debug info warn error"` // LOG_LEVEL. Default info
}

func main() {
    // Define a configuration entry for the logging system
    entry := cfg.ConfigEntry{
        Config: &LoggingConfig{}, // Struct for logging configuration
        BindArray: []cfg.BindValue{
            {
                ValName: "log_format", // Environment variable name
            },
            {
                ValName:    "log_level", // Environment variable name
                DefaultVal: "info",      // Default value if environment variable is not set
            },
        },
    }

    // Register the logging configuration
    if err := cfg.RegisterConfig("log", entry); err != nil {
        fmt.Printf("Error registering config: %v\n", err)
        return
    }

    // Initialize all registered configurations
    if err := cfg.NewConfig(); err != nil {
        fmt.Printf("Error initializing configs: %v\n", err)
        return
    }

    // Access the initialized configuration
    logConfig, ok := cfg.GetConfig("log")
    if !ok {
        fmt.Println("Failed to retrieve the logging config.")
        return
    }

    // Type assert to the specific configuration struct
    if typedConfig, ok := logConfig.(*LoggingConfig); ok {
        fmt.Printf("Initialized Logging Config: %+v\n", typedConfig)
    } else {
        fmt.Println("Failed to type assert logging config.")
    }
}
```

## API Documentation

### Type Descriptions

#### `type ConfigEntry`
- Represents a registered configuration entry.
- **Fields**:
    - `Config any`: The configuration struct to be registered.
    - `BindArray []BindValue`: A list of environment variable bindings.

#### `type BindValue`
- Represents a binding of an environment variable to a configuration field.
- **Fields**:
    - `ValName string`: The environment variable name.
    - `DefaultVal any`: The default value for the environment variable.

#### `type ViperOption`
- A functional option for customizing Viper’s behavior.
- **Example**:
  ```go
  func WithSetEnvPrefix(prefix string) ViperOption {
      // Logic for setting prefix
  }
  ```

### Exported Functions

#### `func RegisterConfig(name string, configStruct ConfigEntry) error`
Registers a custom configuration struct.

#### `func ListConfigs() []string`
Returns a list of all registered configuration names.

#### `func GetConfig(name string) (any, bool)`
Retrieves a registered configuration struct by name.

#### `func NewConfig(list ...ViperOption) error`
Initializes the configuration system and applies functional options.

#### `func WithSetEnvPrefix(EnvPrefix string) ViperOption`
Sets the environment variable prefix for Viper.

## License

This package is licensed under the [MIT License](https://opensource.org/licenses/MIT).

## Thanks

Special thanks to:
- The [Viper library](https://github.com/spf13/viper) maintainers for providing an excellent configuration management tool.
