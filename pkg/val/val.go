// Package val provides a thread-safe validation mechanism using the go-playground/validator/v10 library.
// It supports struct and tag-based validation, as well as custom validation rules.
// The package is designed as a singleton to ensure a single validator instance is used across the application.
package val

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"sync"
)

// Validator defines the methods for validation functionality.
// It abstracts the underlying validator implementation and makes the package testable.
type Validator interface {
	ValidateWithTag(variable any, tag string) error
	ValidateStruct(s any) error
	RegisterValidation(tag string, fn validator.Func) error
}

// singleton holds the single instance of the validator.
// validatorOnce ensures the singleton is initialized only once, making it thread-safe.
var (
	singleton     validatorStruct
	validatorOnce sync.Once
)

// validatorStruct is a wrapper around the go-playground validator.
// It encapsulates the core validation logic and methods to interact with the validator.
type validatorStruct struct {
	validator *validator.Validate
}

// GetValidator returns the singleton instance of the Validator interface.
// It ensures the validator is lazily initialized and thread-safe using sync.Once.
func GetValidator() Validator {
	validatorOnce.Do(func() {
		singleton = newValidator()
	})
	return &singleton
}

// RegisterValidation registers a custom validation function for a specific tag.
// Example usage:
//
//	err := validator.RegisterValidation("is-even", func(fl validator.FieldLevel) bool {
//	    return fl.Field().Int()%2 == 0
//	})
func (v *validatorStruct) RegisterValidation(tag string, fn validator.Func) error {
	return v.validator.RegisterValidation(tag, fn)
}

// ValidateWithTag validates a variable using the provided tag.
// It returns an error if validation fails.
// Example usage:
//
//	err := validator.ValidateWithTag("test@example.com", "email")
func (v *validatorStruct) ValidateWithTag(variable any, tag string) error {
	if err := v.validator.Var(variable, tag); err != nil {
		return handleValidatorError(err)
	}
	return nil
}

// ValidateStruct validates a struct based on its tags.
// It ensures the input is not nil, not empty, and is of type struct.
// If validation fails, it returns a detailed error.
func (v *validatorStruct) ValidateStruct(s any) error {
	if err := validateStruct(s); err != nil {
		return err
	}

	if err := v.validator.Struct(s); err != nil {
		return handleValidatorError(err)
	}
	return nil
}

// TODO: add unit tests for new validations

// addCustomValidators registers custom validation rules with the validator instance.
//
// Custom Validators:
// - "urlprefix": Validates that a string starts with "http://" or "https://".
//
// Parameters:
// - v (*validator.Validate): The validator instance where custom validations will be registered.
//
// This function is intended to be called during validator initialization to
// ensure the custom rules are consistently available across the application.
func addCustomValidators(v *validator.Validate) {
	fn := func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
	}
	_ = v.RegisterValidation("urlprefix", fn)
}

// newValidator initializes and returns a new ValidatorStruct instance.
// It configures the validator with required struct validation enabled.
func newValidator() validatorStruct {
	v := validatorStruct{validator: validator.New(validator.WithRequiredStructEnabled())}
	addCustomValidators(v.validator)
	return v
}

// validateStruct ensures the input is valid for struct validation.
// It checks that the input is not nil, not an uninitialized pointer, and of type struct.
func validateStruct(s any) error {
	if s == nil {
		return fmt.Errorf("input is nil")
	}

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Errorf("input is a nil pointer")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("input is not a struct, got: %s", val.Kind())
	}

	return nil
}

// handleValidatorError processes and formats validation errors.
// It extracts detailed field-specific errors for structured reporting.
func handleValidatorError(err error) error {
	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		var detailedErrors []string
		for _, fe := range valErr {
			detailedErrors = append(detailedErrors, fmt.Sprintf("Field '%s': %s", fe.Field(), fe.Error()))
		}
		return errors.New(strings.Join(detailedErrors, ", "))
	}
	return fmt.Errorf("unexpected validation error: %w", err)
}
