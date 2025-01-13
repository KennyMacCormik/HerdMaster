
# `val` Package Documentation

The `val` package provides a centralized and encapsulated validation mechanism using the [`go-playground/validator`](https://pkg.go.dev/github.com/go-playground/validator/v10) library. It is designed to ensure data integrity across the application by validating structures and variables based on predefined rules and tags.

## Features

1. **Singleton Validator Instance**: Implements a singleton pattern for the validator to leverage caching and thread safety.
2. **Validation for Variables**: Supports validation of individual variables using custom tags.
3. **Validation for Structures**: Enables validation of Go structs based on `validator` tags.
4. **Error Formatting**: Formats validation errors for better debugging and readability.

## Initialization

The package initializes a singleton instance of the validator during package initialization (`init` function). This ensures the validator is ready for use throughout the application without additional setup.

## Usage

### Getting the Validator

Use the `GetValidator` function to obtain a pointer to the singleton `GlobalValidator` instance.

```go
import "your_project/val"

validator := val.GetValidator()
```

### Struct Validation

To validate a struct with `validator` tags:
1. Add validation tags to your struct fields.
2. Call `ValidateStruct` with the struct instance.

Example:

```go
type User struct {
Name  string `validate:"required"`
Email string `validate:"required,email"`
Age   int    `validate:"gte=18"`
}

func main() {
user := User{Name: "John Doe", Email: "invalid-email", Age: 15}
validator := val.GetValidator()

if err := validator.ValidateStruct(user); err != nil {
log.Fatalf("Validation failed: %v", err)
}
}
```

### Variable Validation

To validate an individual variable with a custom tag:
1. Use the `ValidateWithTag` method.
2. Pass the variable and the validation rule.

Example:

```go
func main() {
validator := val.GetValidator()

// Validate if a value is an email
email := "invalid-email"
if err := validator.ValidateWithTag(email, "email"); err != nil {
log.Printf("Invalid email: %v", err)
}
}
```

## Functions

### `GetValidator`

```go
func GetValidator() *GlobalValidator
```

Returns a pointer to the singleton `GlobalValidator` instance.

**Note**: The returned pointer should not be replaced or re-initialized by consumers.

---

### `ValidateStruct`

```go
func (v *GlobalValidator) ValidateStruct(s any) error
```

Validates a struct against the `validator` tags defined in its fields.

- **Parameters**:
  - `s` (`any`): The struct instance to validate.
- **Returns**:
  - `error`: A formatted error (`validator.ValidationErrors`) if validation fails; `nil` otherwise.

---

### `ValidateWithTag`

```go
func (v *GlobalValidator) ValidateWithTag(variable any, tag string) error
```

Validates a single variable against a custom validation tag.

- **Parameters**:
  - `variable` (`any`): The variable to validate.
  - `tag` (`string`): The validation rule to apply (e.g., `required`, `email`, `numeric`).
- **Returns**:
  - `error`: An error if validation fails; `nil` otherwise.

---

### Helper Functions

#### `handleValidatorError`

```go
func handleValidatorError(err error) error
```

Internal helper function to format validation errors into readable messages.

- **Parameters**:
  - `err` (`error`): The error returned by the `validator` package.
- **Returns**:
  - `error`: A formatted error message with details about the failed fields.

## Example Scenarios

### Struct Validation

```go
type Product struct {
    Name  string  `validate:"required"`
    Price float64 `validate:"gte=0"`
}

func main() {
    product := Product{Name: "", Price: -5.0}
    validator := val.GetValidator()

    if err := validator.ValidateStruct(product); err != nil {
        log.Fatalf("Validation failed: %v", err)
    }
}
```

---

### Variable Validation

```go
func main() {
    validator := val.GetValidator()

    // Validate a numeric variable
    value := 42
    if err := validator.ValidateWithTag(value, "gte=50"); err != nil {
        log.Printf("Validation failed: %v", err)
    }
}
```

---