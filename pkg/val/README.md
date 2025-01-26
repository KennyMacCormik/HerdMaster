# Validator Package (`val`) Documentation

## Overview
The `val` package provides a thread-safe validation mechanism built on the `go-playground/validator/v10` library. It supports struct and tag-based validation, custom validation rules, and is designed as a singleton to ensure a single validator instance is used across the application.

## Features
- **Thread-safe validation:** Ensures a single instance of the validator is used application-wide.
- **Custom validation support:** Easily register and use custom validation rules.
- **Struct and tag-based validation:** Validate structs and variables using tags and rules.
- **Error handling:** Detailed error reporting for failed validations.

## Installation

### Dependencies Installation
Install the following dependency:
```bash
go get github.com/go-playground/validator/v10
```

### Package Installation
Install the `val` package:
```sh
go get github.com/KennyMacCormik/HerdMaster/pkg/val
```

## Usage

### Getting the Validator
The `GetValidator` function provides a thread-safe singleton instance of the validator.

```go
import "path/to/val"

validator := val.GetValidator()
```

### Validating a Struct

```go
import "path/to/val"

// Define a struct with validation tags
type User struct {
    Name  string `validate:"required"`
    Email string `validate:"required,email"`
}

// Validate the struct
user := User{Name: "John Doe", Email: "john.doe@example.com"}
validator := val.GetValidator()
err := validator.ValidateStruct(user)
if err != nil {
    fmt.Printf("Validation errors: %v\n", err)
}
```

### Validating a Variable

```go
import "path/to/val"

validator := val.GetValidator()
err := validator.ValidateWithTag("test@example.com", "email")
if err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}
```

### Registering a Custom Validation

```go
import "path/to/val"

validator := val.GetValidator()
err := validator.RegisterValidation("is-even", func(fl validator.FieldLevel) bool {
    return fl.Field().Int()%2 == 0
})
if err != nil {
    fmt.Printf("Failed to register custom validation: %v\n", err)
}
```

## API Documentation

### `func GetValidator() Validator`
Returns the singleton instance of the validator. Ensures thread-safe lazy initialization.

### `func (v *validatorStruct) ValidateStruct(s any) error`
Validates a struct based on its tags. Returns detailed errors if validation fails.

### `func (v *validatorStruct) ValidateWithTag(variable any, tag string) error`
Validates a variable using the provided tag. Returns an error if validation fails.

### `func (v *validatorStruct) RegisterValidation(tag string, fn validator.Func) error`
Registers a custom validation function for a specific tag.

## Type Descriptions

### `type Validator interface`
```go
interface {
    ValidateWithTag(variable any, tag string) error
    ValidateStruct(s any) error
    RegisterValidation(tag string, fn validator.Func) error
}
```
An interface that defines validation functionality.

### `type validatorStruct struct`
A struct that wraps the go-playground validator and provides additional functionality.

## Custom Registered Validation Prefixes

### `urlprefix`
Ensures that a string starts with `http://` or `https://`.

#### Example Usage

```go
import "path/to/val"

validator := val.GetValidator()
err := validator.ValidateWithTag("https://example.com", "urlprefix")
if err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}

// Invalid example
err = validator.ValidateWithTag("ftp://example.com", "urlprefix")
if err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}
```

## License
This package is licensed under the [MIT License](https://opensource.org/licenses/MIT).

## Thanks
Special thanks to the contributors and maintainers of the following libraries used in this package:
- [`go-playground/validator`](https://github.com/go-playground/validator): For providing the core validation functionality.
