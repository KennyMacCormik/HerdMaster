package val

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetValidatorSingleton(t *testing.T) {
	validator1 := GetValidator()
	validator2 := GetValidator()

	assert.Equal(t, validator1, validator2, "expected GetValidator to return the same singleton instance")
}

func TestValidateStruct_Valid(t *testing.T) {
	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	validator := GetValidator()

	// Valid struct
	testObj := TestStruct{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	err := validator.ValidateStruct(testObj)
	assert.NoError(t, err, "expected no validation errors for valid struct")
}

func TestValidateStruct_Invalid(t *testing.T) {
	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	validator := GetValidator()

	// Invalid struct
	testObj := TestStruct{
		Name:  "",
		Email: "invalid-email",
	}

	err := validator.ValidateStruct(testObj)
	require.Error(t, err, "expected validation errors for invalid struct")
	assert.Contains(t, err.Error(), "Name", "expected validation error for 'Name'")
	assert.Contains(t, err.Error(), "Email", "expected validation error for 'Email'")
}

func TestValidateWithTag_Valid(t *testing.T) {
	validator := GetValidator()

	err := validator.ValidateWithTag("test@example.com", "email")
	assert.NoError(t, err, "expected no validation errors for valid email")
}

func TestValidateWithTag_Invalid(t *testing.T) {
	validator := GetValidator()

	err := validator.ValidateWithTag("invalid-email", "email")
	assert.Error(t, err, "expected validation errors for invalid email")
	assert.Contains(t, err.Error(), "email", "expected validation error for 'email' tag")
}

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

func TestHandleValidatorError(t *testing.T) {
	// Test unexpected error
	err := handleValidatorError(errors.New("unexpected"))
	assert.EqualError(t, err, "unexpected validation error: unexpected", "expected formatted error message for unexpected error")

	// Test validation errors
	valErr := validator.ValidationErrors{}
	err = handleValidatorError(valErr)
	assert.Equal(t, valErr, err, "expected validation errors to be returned as-is")
}
