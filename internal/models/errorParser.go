package models

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ErrorEnum int

const (
	isRequired ErrorEnum = iota
	isEmail
	isMax
	isMin
)

var stateName = map[ErrorEnum]string{
	isRequired: "required",
	isEmail:    "email",
	isMax:      "max",
	isMin:      "min",
}

var (
	ErrUnknownField  = "UnknownField"
	ErrOnValidations = "error while validating the request JSON data, please check if the data provided is correct"
)

func ValidatorErrorParser(err error) *ErrorsStruct {
	var validateErrs validator.ValidationErrors
	errMessages := make(map[string]string)

	if !(errors.As(err, &validateErrs)) {
		errMessages[ErrUnknownField] = err.Error()
		return &ErrorsStruct{Validations: errMessages, Message: ErrOnValidations}
	}

	//TODO: improve this
	for _, e := range validateErrs {
		switch e.Tag() {
		case stateName[isRequired]:
			errMessages[e.Field()] = fmt.Sprintf("%s is required", e.Field())
		case stateName[isEmail]:
			errMessages[e.Field()] = fmt.Sprintf("%s must be a valid email address", e.Field())
		case stateName[isMax]:
			errMessages[e.Field()] = fmt.Sprintf("%s must be less than %s", e.Field(), e.Param())
		case stateName[isMin]:
			errMessages[e.Field()] = fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
		default:
			errMessages[e.Field()] = fmt.Sprintf("%s is not valid (%s)", e.Field(), e.Tag())
		}
	}
	return &ErrorsStruct{Validations: errMessages, Message: ErrOnValidations}
}
