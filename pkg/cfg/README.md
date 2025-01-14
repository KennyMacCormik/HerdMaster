# cfg Package

## Overview

The `cfg` package is a flexible and extensible configuration management system that integrates with the Viper library. It provides an easy-to-use API for registering configuration structs, binding environment variables, and setting default values dynamically.

## Features

- Seamless integration with the [Viper](https://github.com/spf13/viper) library for configuration management.
- Support for dynamic registration of configuration structs.
- Environment variable binding and default value management.
- Thread-safe operations using synchronization primitives.
- Support for custom validation of configuration fields.

## Installation

### Dependencies Installation

To use this package, you need to install its dependencies:

```bash
go get github.com/spf13/viper
go get github.com/stretchr/testify
```

### Package Installation

Install the `cfg` package directly in your Go project:

```bash
go get github.com/KennyMacCormik/HerdMaster/pkg/cfg
```

## Usage

Below is an example of how to use the `cfg` package to manage configurations in your project:

```go
package main

import (
	"fmt"
	"github.com/KennyMacCormik/HerdMaster/pkg/cfg"
)

type LoggingConfig struct {
	Format string `mapstructure:"log_format" validate:"oneof=text json"`
	Level  string `mapstructure:"log_level" validate:"oneof=debug info warn error"`
}

func main() {
	entry := cfg.ConfigEntry{
		Config: &LoggingConfig{},
		BindArray: []cfg.BindValue{
			{ValName: "log_format"},
			{ValName: "log_level", DefaultVal: "info"},
		},
	}

	if err := cfg.RegisterConfig("log", entry); err != nil {
		fmt.Printf("Error registering config: %v
", err)
		return
	}

	if err := cfg.NewConfig(); err != nil {
		fmt.Printf("Error initializing configs: %v
", err)
		return
	}

	logConfig, ok := cfg.GetConfig("log")
	if !ok {
		fmt.Println("Failed to retrieve the logging config.")
		return
	}

	if typedConfig, ok := logConfig.(*LoggingConfig); ok {
		fmt.Printf("Initialized Logging Config: %+v
", typedConfig)
	} else {
		fmt.Println("Failed to type assert logging config.")
	}
}
```

## API Documentation

### `RegisterConfig`

Registers a configuration struct and its associated bindings.

```go
func RegisterConfig(name string, configStruct ConfigEntry) error
```

### `NewConfig`

Initializes the configuration system, binds all registered entries, and loads their values from the environment.

```go
func NewConfig() error
```

### `ListConfigs`

Returns a list of all registered configuration names.

```go
func ListConfigs() []string
```

### `GetConfig`

Retrieves a registered configuration struct by name.

```go
func GetConfig(name string) (any, bool)
```

### Internal Functions

#### `validateConfigStruct`

Validates that the `Config` field of a `ConfigEntry` is a pointer to a struct.

#### `bindActualValue`

Binds environment variables and unmarshals values into a configuration struct.

## Type Descriptions

### `ConfigEntry`

Represents a configuration entry, including the struct and associated bindings.

```go
type ConfigEntry struct {
	Config    any
	BindArray []BindValue
}
```

### `BindValue`

Represents a single binding of an environment variable to a configuration field, with an optional default value.

```go
type BindValue struct {
	ValName    string
	DefaultVal any
}
```

## var Descriptions

### `configEntries`

Holds all registered configuration entries.

```go
var configEntries = make(map[string]ConfigEntry)
```

### `mtx`

Provides a mutex for thread-safe operations.

```go
var mtx sync.RWMutex
```

## License

This project is licensed under the MIT License. For more details, see [MIT License](https://opensource.org/licenses/MIT).

## Thanks

Special thanks to:
- [Viper](https://github.com/spf13/viper) for providing a robust configuration management library.
- [Testify](https://github.com/stretchr/testify) for simplifying unit testing.
