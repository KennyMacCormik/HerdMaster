# val Package

## Overview
The `val` package provides a thread-safe validation mechanism using the `go-playground/validator/v10` library. It supports:
- Struct-based validation.
- Tag-based validation.
- Custom validation rules.

The package is designed as a singleton to ensure a single validator instance is used throughout the application, promoting efficiency and consistency.

## Features
- **Thread-Safe Singleton**: Ensures only one validator instance is active.
- **Struct Validation**: Validate structs based on predefined tags.
- **Tag Validation**: Validate individual variables using custom tags.
- **Custom Validation Rules**: Easily extend functionality by registering new validation rules.

## Installation

### Dependencies
This package depends on:
- `go-playground/validator/v10`
- `github.com/stretchr/testify` (for testing purposes)

Install dependencies using `go get`:
```bash
go get github.com/go-playground/validator/v10
go get github.com/stretchr/testify
```

### Package Installation
To install the `val` package:
```bash
go get github.com/KennyMacCormik/HerdMaster/pkg/val
```

## Usage
Below are examples of how to use the `val` package.

### Struct Validation
```go
package main

import (
	"fmt"
	"github.com/KennyMacCormik/HerdMaster/pkg/val"
)

type User struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func main() {
	validator := val.GetValidator()

	user := User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	if err := validator.ValidateStruct(user); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed!")
	}
}
```

### Tag Validation
```go
package main

import (
	"fmt"
	"github.com/KennyMacCormik/HerdMaster/pkg/val"
)

func main() {
	validator := val.GetValidator()

	email := "invalid-email"
	if err := validator.ValidateWithTag(email, "email"); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed!")
	}
}
```

### Custom Validation
```go
package main

import (
	"fmt"
	"github.com/KennyMacCormik/HerdMaster/pkg/val"
	"github.com/go-playground/validator/v10"
)

func main() {
	validator := val.GetValidator()

	// Register a custom validation
	err := validator.RegisterValidation("is-even", func(fl validator.FieldLevel) bool {
		return fl.Field().Int()%2 == 0
	})

	if err != nil {
		fmt.Println("Error registering custom validation:", err)
		return
	}

	type TestStruct struct {
		Number int `validate:"is-even"`
	}

	test := TestStruct{Number: 4}
	if err := validator.ValidateStruct(test); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed!")
	}
}
```

## API Documentation

### `Validator` Interface
The `Validator` interface defines the following methods:
```go
type Validator interface {
	ValidateWithTag(variable any, tag string) error
	ValidateStruct(s any) error
	RegisterValidation(tag string, fn validator.Func) error
}
```

### Functions

#### `GetValidator`
Returns the singleton instance of the `Validator` interface.
```go
func GetValidator() Validator
```

#### `RegisterValidation`
Registers a custom validation function for a specific tag.
```go
func (v *validatorStruct) RegisterValidation(tag string, fn validator.Func) error
```

#### `ValidateWithTag`
Validates a variable using the provided tag.
```go
func (v *validatorStruct) ValidateWithTag(variable any, tag string) error
```

#### `ValidateStruct`
Validates a struct based on its tags.
```go
func (v *validatorStruct) ValidateStruct(s any) error
```

## Type Description

### `validatorStruct`
The `validatorStruct` is a wrapper around the `go-playground/validator` and implements the `Validator` interface.

### `Validator`
An interface abstracting the validation functionality.

## License
This project is licensed under the MIT License. See the [LICENSE](https://opensource.org/licenses/MIT) file for details.

## Thanks
Special thanks to the contributors and maintainers of:
- [`go-playground/validator`](https://github.com/go-playground/validator)
- [`testify`](https://github.com/stretchr/testify)
