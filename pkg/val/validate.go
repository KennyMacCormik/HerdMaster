package val

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func init() {
	ValInstance = newValidator()
}

type ErrTagViolated struct {
	Val any
	Tag string
}

func (e ErrTagViolated) Error() string {
	return fmt.Sprintf("invalid value %q: violates tag %q", e.Val, e.Tag)
}

func NewErrTagViolated(Val any, Tag string) error {
	return ErrTagViolated{
		Val: Val,
		Tag: Tag,
	}
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
		return NewErrTagViolated(variable, tag)
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
	for _, v := range valErr {
		tag := reflectActualTag(s, v.StructField())
		if tag == "" {
			tag = "err reflect tag"
		}
		return NewErrTagViolated(v.Value(), tag)
	}
	return nil
}

func reflectActualTag(s any, sf string) string {
	ref := reflect.TypeOf(s)

	for i := 0; i < ref.NumField(); i++ {
		fieldName := ref.Field(i).Name
		field, _ := ref.FieldByName(fieldName)
		if field.Type.Name() != "bool" {
			for j := 0; j < field.Type.NumField(); j++ {
				intFieldName := field.Type.Field(j)
				if intFieldName.Name == sf {
					return intFieldName.Tag.Get("validate")
				}
			}
		}
	}

	return ""
}
