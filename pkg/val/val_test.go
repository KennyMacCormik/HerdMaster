package val

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestGetValidatorSingleton verifies that GetValidator always returns the same singleton instance.
func TestGetValidatorSingleton(t *testing.T) {
	validator1 := GetValidator()
	validator2 := GetValidator()

	assert.Equal(t, validator1, validator2, "expected GetValidator to return the same singleton instance")
}

// TestValidateStruct_Valid verifies that ValidateStruct successfully validates a valid struct.
func TestValidateStruct_Valid(t *testing.T) {
	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	validatorVar := GetValidator()

	// Valid struct
	testObj := TestStruct{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	err := validatorVar.ValidateStruct(testObj)
	assert.NoError(t, err, "expected no validation errors for valid struct")
}

// TestValidateStruct_Invalid verifies that ValidateStruct returns detailed validation errors for invalid structs.
func TestValidateStruct_Invalid(t *testing.T) {
	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	validatorVar := GetValidator()

	// Invalid struct
	testObj := TestStruct{
		Name:  "",
		Email: "invalid-email",
	}

	err := validatorVar.ValidateStruct(testObj)
	require.Error(t, err, "expected validation errors for invalid struct")
	assert.Contains(t, err.Error(), "Name", "expected validation error for 'Name'")
	assert.Contains(t, err.Error(), "Email", "expected validation error for 'Email'")
}

// TestValidateWithTag_Valid ensures that ValidateWithTag correctly validates variables against valid tags.
func TestValidateWithTag_Valid(t *testing.T) {
	validatorVar := GetValidator()

	err := validatorVar.ValidateWithTag("test@example.com", "email")
	assert.NoError(t, err, "expected no validation errors for valid email")
}

// TestValidateWithTag_Invalid ensures that ValidateWithTag returns errors for invalid variable and tag combinations.
func TestValidateWithTag_Invalid(t *testing.T) {
	validatorVar := GetValidator()

	err := validatorVar.ValidateWithTag("invalid-email", "email")
	assert.Error(t, err, "expected validation errors for invalid email")
	assert.Contains(t, err.Error(), "email", "expected validation error for 'email' tag")
}

// TestRegisterValidation verifies the functionality of registering and using custom validation rules.
func TestRegisterValidation(t *testing.T) {
	val := GetValidator()

	// Register custom validation
	err := val.RegisterValidation("is-even", func(fl validator.FieldLevel) bool {
		return fl.Field().Int()%2 == 0
	})
	assert.NoError(t, err, "expected no error when registering custom validation")

	// Define a struct to test the custom validation
	type TestStruct struct {
		Value int `validate:"is-even"`
	}

	// Valid value
	validObj := TestStruct{Value: 4}
	err = val.ValidateStruct(validObj)
	assert.NoError(t, err, "expected no validation errors for even value")

	// Invalid value
	invalidObj := TestStruct{Value: 3}
	err = val.ValidateStruct(invalidObj)
	require.Error(t, err, "expected validation error for odd value")
	assert.Contains(t, err.Error(), "is-even", "expected custom validation error for 'is-even'")
}

// TestHandleUnexpectedValidatorError verifies the handling of unexpected validation errors.
func TestHandleUnexpectedValidatorError(t *testing.T) {
	// Test unexpected error
	err := handleValidatorError(errors.New("unexpected"))
	assert.EqualError(t, err, "unexpected validation error: unexpected", "expected formatted error message for unexpected error")

	// Test validation errors
	valErr := validator.ValidationErrors{}
	err = handleValidatorError(valErr)
	assert.NotEqual(t, valErr, err, "unexpected validation errors returned is expected to be not of type ValidationErrors")
}

// TestRegisterValidation_NilCallback ensures that attempting to register a nil validation function results in an error.
func TestRegisterValidation_NilCallback(t *testing.T) {
	val := GetValidator()

	// Attempt to register a custom validation with a nil function
	err := val.RegisterValidation("nil-validation", nil)
	assert.Error(t, err, "expected an error when registering a nil validation function")
}

// TestRegisterValidation_PanicInCallback ensures that panics in custom validation callbacks are handled as expected.
func TestRegisterValidation_PanicInCallback(t *testing.T) {
	val := GetValidator()

	// Register a custom validation that panics
	err := val.RegisterValidation("panic-validation", func(fl validator.FieldLevel) bool {
		panic("intentional panic in validation")
	})
	assert.NoError(t, err, "expected no error when registering a valid custom validation")

	// Define a struct to test the custom validation
	type TestStruct struct {
		Value int `validate:"panic-validation"`
	}

	// Test struct that triggers the panic
	testObj := TestStruct{Value: 5}

	// Validate the struct
	assert.Panics(t, func() {
		_ = val.ValidateStruct(testObj)
	}, "expected validation to panic due to callback panic")
}

// TestValidateStruct_UnsupportedType verifies that ValidateStruct returns an error for unsupported input types.
func TestValidateStruct_UnsupportedType(t *testing.T) {
	val := GetValidator()

	// Test unsupported type (e.g., a basic int instead of a struct)
	unsupported := 123

	err := val.ValidateStruct(unsupported)
	assert.Error(t, err, "expected an error when validating an unsupported type")
	assert.Contains(t, err.Error(), "input is not a struct", "expected error message indicating type mismatch")
}

// TestValidateStruct_UninitializedPointer ensures that uninitialized pointers to structs result in an appropriate error.
func TestValidateStruct_UninitializedPointer(t *testing.T) {
	val := GetValidator()

	// Test uninitialized pointer to a struct
	var uninitialized *struct {
		Field1 string `validate:"required"`
	}

	err := val.ValidateStruct(uninitialized)
	assert.Error(t, err, "expected an error when validating an uninitialized pointer")
	assert.Contains(t, err.Error(), "nil pointer", "expected error message indicating nil pointer")
}
