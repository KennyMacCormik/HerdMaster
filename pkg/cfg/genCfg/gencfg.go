// Package genCfg contains predefined configuration types that streamline the use of the cfg package.
// These structs are designed to work seamlessly with the val package, ensuring compatibility with
// validation rules and enabling robust configuration management across HerdMaster microservices.
//
// Key Features:
// - Predefined structs for commonly used configurations like gRPC and logging.
// - Compatibility with the val package for validation out-of-the-box.
// - Integration with the cfg package for centralized and reusable configuration management.
//
// Example Usage (Cross-Package Integration):
//
// 1. Define a configuration entry and register it using the cfg package:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/KennyMacCormik/HerdMaster/pkg/cfg"
//	    "github.com/KennyMacCormik/HerdMaster/pkg/genCfg"
//	    "github.com/KennyMacCormik/HerdMaster/pkg/val"
//	)
//
//	func main() {
//	    // Register the gRPC configuration
//	    grpcEntry := cfg.ConfigEntry{
//	        Config: &genCfg.GrpcConfig{}, // Use predefined GrpcConfig from genCfg
//	        BindArray: []cfg.BindValue{
//	            {ValName: "grpc_host", DefaultVal: "127.0.0.1"},
//	            {ValName: "grpc_port", DefaultVal: "50051"},
//	        },
//	    }
//
//	    // Register with cfg
//	    err := cfg.RegisterConfig("grpc", grpcEntry)
//	    if err != nil {
//	        fmt.Printf("Error registering gRPC config: %v\n", err)
//	        return
//	    }
//
//	    // Initialize configurations
//	    err = cfg.NewConfig()
//	    if err != nil {
//	        fmt.Printf("Error initializing configurations: %v\n", err)
//	        return
//	    }
//
//	    // Retrieve the configuration
//	    conf, ok := cfg.GetConfig("grpc")
//	    if !ok {
//	        fmt.Println("Failed to retrieve gRPC config")
//	        return
//	    }
//
//	    // Type assert to the specific configuration struct
//	    grpcConf, ok := conf.(*genCfg.GrpcConfig)
//	    if !ok {
//	        fmt.Println("Type assertion for gRPC config failed")
//	        return
//	    }
//
//	    // Validate the typed configuration struct
//	    validator := val.GetValidator()
//	    err = validator.ValidateStruct(grpcConf)
//	    if err != nil {
//	        fmt.Printf("Validation error: %v\n", err)
//	        return
//	    }
//
//	    // Access and print the gRPC configuration
//	    fmt.Printf("gRPC Config - Host: %s, Port: %s\n", grpcConf.Host, grpcConf.Port)
//	}
//
// This example demonstrates how to leverage genCfg, cfg, and val packages together to register,
// validate, and use a predefined configuration struct for gRPC settings.
package genCfg

// GrpcConfig represents the configuration for gRPC servers.
// Host: The IP address or hostname of the gRPC server.
// Validates as IPv4 or hostname (RFC1123).
// Port: The port number for the gRPC server.
// Validates as a numeric value between 1025 and 65 535 (exclusive).
type GrpcConfig struct {
	Host string `mapstructure:"grpc_host" validate:"ip4_addr|hostname_rfc1123,required"`
	Port string `mapstructure:"grpc_port" validate:"numeric,gt=1024,lt=65536,required"`
}

// LoggingConfig represents the configuration for logging systems.
// In particular, log package from this repo.
// Format: Specifies the log format, either "text" or "json".
// Level: Specifies the log level, which must be one of "debug", "info", "warn", or "error".
type LoggingConfig struct {
	Format string `mapstructure:"log_format" validate:"oneof=text json"`
	Level  string `mapstructure:"log_level" validate:"oneof=debug info warn error"`
}
