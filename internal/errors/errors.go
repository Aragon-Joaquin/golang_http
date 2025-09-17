package er

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

type ErrorsStruct struct {
	Validations map[string]string `json:"validationsErrors,omitempty"`
	Message     any               `json:"message"`
}

// ! just a swagegr placeholder until i fix this somehow
type ReturnedError struct {
	Data struct {
		Validations map[string]string `json:"validationsErrors,omitempty"`
		Message     string            `json:"message"`
	} `json:"data"`
	Error bool `json:"error"`
}

func ValidatorErrorParser(err error) *ErrorsStruct {
	var validateErrs validator.ValidationErrors
	errMessages := make(map[string]string)

	if !(errors.As(err, &validateErrs)) {
		errMessages[string(UnknownField)] = err.Error()
		return &ErrorsStruct{Validations: errMessages, Message: OnValidations}
	}

	//TODO: improve this
	for _, e := range validateErrs {
		switch e.Tag() {
		case Validator_ErrMap[isRequired]:
			errMessages[e.Field()] = fmt.Sprintf("%s is required", e.Field())
		case Validator_ErrMap[isEmail]:
			errMessages[e.Field()] = fmt.Sprintf("%s must be a valid email address", e.Field())
		case Validator_ErrMap[isMax]:
			errMessages[e.Field()] = fmt.Sprintf("%s must be less than %s", e.Field(), e.Param())
		case Validator_ErrMap[isMin]:
			errMessages[e.Field()] = fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
		default:
			errMessages[e.Field()] = fmt.Sprintf("%s is not valid (%s)", e.Field(), e.Tag())
		}
	}
	return &ErrorsStruct{Validations: errMessages, Message: OnValidations}
}

//todo:
// make possible error checking with the ps.ErrorCode
// write all the generic return errors here

// ! generic errors
func CheckForGenericErrors(err error) *ErrorsStruct {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case string(PGErr_NOT_NULL):
			return &ErrorsStruct{Validations: map[string]string{pgErr.Where: PG_ErrorCodes[PGErr_NOT_NULL]}, Message: DBConflict}
		case string(PGErr_UNIQUE):
			columnName := strings.Split(pgErr.ConstraintName, "_")[1]
			return &ErrorsStruct{Validations: map[string]string{columnName: PG_ErrorCodes[PGErr_UNIQUE]}, Message: DBConflict}
		case string(PGErr_UNDEFINED_COLUMN):
			return &ErrorsStruct{Message: UndefinedCol}
		default:
			return &ErrorsStruct{Message: pgErr.Message}
		}

	}

	switch err.Error() {
	case "context deadline exceeded":
		return &ErrorsStruct{Message: QueryTimeout}
	case sql.ErrNoRows.Error():
		return &ErrorsStruct{Message: NotFound}
	default:
		return &ErrorsStruct{Message: Unknown}
	}

}

func MatchErrorCodes(message any) int {
	errMsg, ok := message.(errIdentifier)

	if !ok {
		//check if the type pgConnErrorCode its thrown here
		log.Printf("bad status in MatchErrorCodes: %s", message)
		return http.StatusTeapot
	}

	switch errMsg {
	case QueryTimeout:
		return http.StatusRequestTimeout
	case NotFound:
		return http.StatusNotFound
	case DBConflict:
		return http.StatusConflict
	case OnValidations:
		return http.StatusBadRequest
	case UndefinedCol:
		return http.StatusBadRequest
	default:
		return http.StatusNotImplemented
	}
}
