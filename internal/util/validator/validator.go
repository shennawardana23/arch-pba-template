package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/mochammadshenna/arch-pba-template/internal/model/api"
	"github.com/mochammadshenna/arch-pba-template/internal/util/exceptioncode"
)

var Validator *validator.Validate

func New() *validator.Validate {
	Validator = validator.New()

	return Validator
}

func Validate(e interface{}) error {
	err := Validator.Struct(e)
	if err == nil {
		return err
	}
	errors := []api.ErrorValidate{}
	if err != nil {
		for _, er := range err.(validator.ValidationErrors) {
			errors = append(errors, api.ErrorValidate{
				Key:     er.Field(),
				Code:    "VALIDATION",
				Message: er.Error(),
			})
		}
	}
	return api.ErrorResponse{
		Code:    exceptioncode.CodeInvalidValidation,
		Message: "validation error",
		Errors:  errors,
	}
}
