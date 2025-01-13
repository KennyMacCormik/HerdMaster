
# `val` Package Documentation

The `val` package provides a centralized validation mechanism for the application using the `go-playground/validator/v10` library. It ensures thread-safe, reusable validation with support for custom rules and abstraction via an interface.

---

## Features

1. **Singleton Validator**: A thread-safe, globally available validator instance.
2. **Struct Validation**: Validates struct fields using `validate` tags.
3. **Tag-Based Validation**: Validates individual variables against custom tags.
4. **Custom Validation Rules**: Allows registering custom validation functions.
5. **Abstraction with Interface**: Exposes a `Validator` interface for decoupled and testable design.
6. **Error Formatting**: Formats validation errors consistently.

---

## Usage

### GetValidator

```go
func GetValidator() Validator
```

Returns the singleton instance of the `Validator` interface.

#### Example:

```go
validator := GetValidator()

type User struct {
    Name  string `validate:"required"`
    Email string `validate:"required,email"`
}

user := User{Name: "John Doe", Email: "john.doe@example.com"}

err := validator.ValidateStruct(user)
if err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}
```

---

### ValidateStruct

```go
func (v *GlobalValidator) ValidateStruct(s any) error
```

Validates a struct based on `validate` tags defined on its fields.

#### Example:

```go
type Product struct {
    Name  string  `validate:"required"`
    Price float64 `validate:"required,gt=0"`
}

product := Product{Name: "Laptop", Price: 1500.00}
err := validator.ValidateStruct(product)
if err != nil {
    fmt.Printf("Validation errors: %v\n", err)
}
```

---

### ValidateWithTag

```go
func (v *GlobalValidator) ValidateWithTag(variable any, tag string) error
```

Validates a variable against a specific validation tag.

#### Example:

```go
err := validator.ValidateWithTag("test@example.com", "email")
if err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}
```

---

### RegisterValidation

```go
func (v *GlobalValidator) RegisterValidation(tag string, fn validator.Func) error
```

Registers a custom validation function for a given tag.

#### Example:

```go
validator := GetValidator()

err := validator.RegisterValidation("is-even", func(fl validator.FieldLevel) bool {
    value := fl.Field().Int()
    return value%2 == 0
})

if err != nil {
    log.Fatalf("Failed to register custom validation: %v\n", err)
}

type Number struct {
    Value int `validate:"is-even"`
}

num := Number{Value: 3}
err = validator.ValidateStruct(num)
if err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}
```

---

### Error Handling

The `handleValidatorError` function formats validation errors.

#### Example:

```go
func handleValidatorError(err error) error {
    var valErr validator.ValidationErrors
    if errors.As(err, &valErr) {
        return valErr
    }
    return fmt.Errorf("unexpected validation error: %w", err)
}
```

---

## Integration

This package is designed to work as a singleton, ensuring efficient and consistent validation across the application. Use it wherever validation logic is required.

---

## Example

```go
package main

import (
    "fmt"
    "your_project/val"
)

func main() {
    validator := val.GetValidator()

    type User struct {
        Name  string `validate:"required"`
        Email string `validate:"required,email"`
    }

    user := User{Name: "", Email: "invalid-email"}

    err := validator.ValidateStruct(user)
    if err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    }
}
```

---

## Limitations

1. **Custom Rules**: Ensure custom rules are registered before use.
2. **Error Interpretation**: Validation errors are returned in raw form, requiring custom formatting if needed.

---

## Unit Tests

The `val` package includes tests for:
1. Singleton behavior (`GetValidator`)
2. Struct and tag-based validation
3. Custom validation rules
4. Error handling

Run the tests with:

```bash
go test ./... -v
```

---
