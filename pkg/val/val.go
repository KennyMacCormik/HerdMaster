package val

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"sync"
)

type Validator interface {
	ValidateWithTag(variable any, tag string) error
	ValidateStruct(s any) error
	RegisterValidation(tag string, fn validator.Func) error
}

var (
	singleton     GlobalValidator
	validatorOnce sync.Once
)

// GlobalValidator is a wrapper around a validator object.
type GlobalValidator struct {
	validator *validator.Validate
}

// newValidator initialize validator singleton.
// See https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Singleton for details
func newValidator() GlobalValidator {
	return GlobalValidator{validator: validator.New(validator.WithRequiredStructEnabled())}
}

// GetValidator returns a Validator interface
func GetValidator() Validator {
	validatorOnce.Do(func() {
		singleton = newValidator()
	})
	return &singleton
}

// ValidateWithTag accepts validator tag according to https://pkg.go.dev/github.com/go-playground/validator/v10#readme-baked-in-validations
func (v *GlobalValidator) ValidateWithTag(variable any, tag string) error {
	if err := v.validator.Var(variable, tag); err != nil {
		return handleValidatorError(err)
	}
	return nil
}

// ValidateStruct validates structure based on structure's tags
func (v *GlobalValidator) ValidateStruct(s any) error {
	if err := v.validator.Struct(s); err != nil {
		return handleValidatorError(err)
	}
	return nil
}

// RegisterValidation adds a custom validation function for a given tag.
// Example:
// validator := GetValidator()
//
//	err := validator.RegisterValidation("is-even", func(fl validator.FieldLevel) bool {
//	    value := fl.Field().Int()
//	    return value%2 == 0
//	})
func (v *GlobalValidator) RegisterValidation(tag string, fn validator.Func) error {
	return v.validator.RegisterValidation(tag, fn)
}

// handleValidatorError used to format validator's errors
func handleValidatorError(err error) error {
	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		return valErr
	}
	return errors.New("unexpected validation error: " + err.Error())
}
