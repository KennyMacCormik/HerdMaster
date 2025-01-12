package val

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

func init() {
	ValInstance = newValidator()
}

var ValInstance Validator

type Validator struct {
	validator *validator.Validate
}

func newValidator() Validator {
	return Validator{validator: validator.New(validator.WithRequiredStructEnabled())}
}

func (v *Validator) ValidateWithTag(variable any, tag string) error {
	if err := v.validator.Var(variable, tag); err != nil {
		return err
	}
	return nil
}

func (v *Validator) ValidateStruct(s any) error {
	if err := v.validator.Struct(s); err != nil {
		return handleValidatorError(s, err)
	}
	return nil
}

func handleValidatorError(s any, err error) error {
	var valErr validator.ValidationErrors
	errors.As(err, &valErr)
	return valErr
}
