package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

//TODO: this file is very messy, ton of inconsistencies.

// ! pgConn error code status
type PgConnErrorCode string

var (
	PGErr_NOT_NULL         PgConnErrorCode = "23502"
	PGErr_UNIQUE           PgConnErrorCode = "23505"
	PGErr_UNDEFINED_COLUMN PgConnErrorCode = "42703"
)

var PG_ErrorCodes = map[PgConnErrorCode]string{
	PGErr_NOT_NULL:         "NotNull constraint violation",
	PGErr_UNIQUE:           "Unique constraint violation",
	PGErr_UNDEFINED_COLUMN: "Undefined column",
}

// ! validator errors enum
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

// TODO: make map
var (
	UnknownField = "UnknownField"

	// psConn Errors
	ErrMsg_UndefinedCol  = "Undefined column."
	ErrMsg_OnValidations = "Error while validating the request JSON data, please check if the data provided is correct."
	ErrMsg_DBConflict    = "Data conflicting existing with previous registers. "

	// generic Errors
	ErrMsg_QueryTimeout = "Query timeout. Please try again later."
	ErrMsg_NotFound     = "Resource not found."
	ErrMsg_Conflict     = "Resource already exists."
	ErrMsg_Unknown      = "Unknown Error."

	ErrMsg_JSONReading = "Error while reading the JSON data. Exceeded the 1mb limit."
)

func ValidatorErrorParser(err error) *ErrorsStruct {
	var validateErrs validator.ValidationErrors
	errMessages := make(map[string]string)

	if !(errors.As(err, &validateErrs)) {
		errMessages[UnknownField] = err.Error()
		return &ErrorsStruct{Validations: errMessages, Message: ErrMsg_OnValidations}
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
	return &ErrorsStruct{Validations: errMessages, Message: ErrMsg_OnValidations}
}

//todo:
// make possible error checking with the ps.ErrorCode
// write all the generic return errors here

// ! generic errors
func CheckForGenericErrors(err error) *ErrorsStruct {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case string(PGErr_NOT_NULL):
			return &ErrorsStruct{Validations: map[string]string{pgErr.Where: PG_ErrorCodes[PGErr_NOT_NULL]}, Message: ErrMsg_DBConflict}
		case string(PGErr_UNIQUE):
			columnName := strings.Split(pgErr.ConstraintName, "_")[1]
			return &ErrorsStruct{Validations: map[string]string{columnName: PG_ErrorCodes[PGErr_UNIQUE]}, Message: ErrMsg_DBConflict}
		case string(PGErr_UNDEFINED_COLUMN):
			return &ErrorsStruct{Message: ErrMsg_UndefinedCol}
		default:
			return &ErrorsStruct{Message: pgErr.Message}
		}

	}

	switch err.Error() {
	case "context deadline exceeded":
		return &ErrorsStruct{Message: ErrMsg_QueryTimeout}
	case sql.ErrNoRows.Error():
		return &ErrorsStruct{Message: ErrMsg_NotFound}
	default:
		return &ErrorsStruct{Message: ErrMsg_Unknown}
	}

}
