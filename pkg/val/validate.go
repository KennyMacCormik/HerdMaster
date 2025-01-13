package val

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

func init() {
	singleton = newValidator()
}

var singleton GlobalValidator

// GlobalValidator is a wrapper around a validator object.
type GlobalValidator struct {
	validator *validator.Validate
}

// newValidator initialize validator sinleton.
// See https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Singleton for details
func newValidator() GlobalValidator {
	return GlobalValidator{validator: validator.New(validator.WithRequiredStructEnabled())}
}

// GetValidator returns a pointer to the validator singleton.
// Used for DI compliance.
// Note: The returned pointer shouldn't be replaced or re-initialized by consumers.
func GetValidator() *GlobalValidator {
	return &singleton
}

// ValidateWithTag accepts validator tag according to https://pkg.go.dev/github.com/go-playground/validator/v10#readme-baked-in-validations
func (v *GlobalValidator) ValidateWithTag(variable any, tag string) error {
	if err := v.validator.Var(variable, tag); err != nil {
		return err
	}
	return nil
}

// ValidateStruct validates structure based on structure's tags
func (v *GlobalValidator) ValidateStruct(s any) error {
	if err := v.validator.Struct(s); err != nil {
		return handleValidatorError(s, err)
	}
	return nil
}

// handleValidatorError used to format validator's errors
func handleValidatorError(s any, err error) error {
	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		return valErr
	}
	return errors.New("unexpected validation error: " + err.Error())
}
