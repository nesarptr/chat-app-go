package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	FailedField string `json:"failedField"`
	Message     string `json:"message"`
}

var validate = validator.New()

func ValidateStruct(uncheckedStruct interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(uncheckedStruct)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = strings.Split(err.StructNamespace(), ".")[1]
			element.Message = element.FailedField + " is not valid"
			errors = append(errors, &element)
		}
	}
	return errors
}
