package utils

import "github.com/go-playground/validator/v10"

type Validator struct {
	Validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		Validate: validator.New(),
	}
}
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.Validate.Struct(s)
}
