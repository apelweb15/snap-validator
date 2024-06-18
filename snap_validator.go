package snap_validator

import (
	"github.com/apelweb15/snap-validator/internal/validator"
)

var snapValidator *SnapValidator

func init() {
	snapValidator = New()
}

func New() *SnapValidator {
	v := new(SnapValidator)
	v.validator = validator.NewService[any]()
	return v
}

type SnapValidator struct {
	validator validator.Validator[any]
}

func ValidateStruct(data interface{}, serviceCode string) error {
	return snapValidator.ValidateStruct(data, serviceCode)
}
func (v *SnapValidator) ValidateStruct(data interface{}, serviceCode string) error {
	return v.validator.ValidateStructSnapServiceCode(data, serviceCode)
}
