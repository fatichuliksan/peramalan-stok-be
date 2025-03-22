package validator

import "github.com/go-playground/validator/v10"

type ValidatorHelper struct {
	validate *validator.Validate
}

// Interface ...
type Interface interface {
	Validate(i interface{}) error
	Validator() *ValidatorHelper
}

// NewValidator ...
func NewValidator() Interface {
	return &ValidatorHelper{
		validate: validator.New(),
	}
}

func (t *ValidatorHelper) Validate(i interface{}) error {
	return t.validate.Struct(i)
}

func (t *ValidatorHelper) Validator() *ValidatorHelper {
	return t
}
