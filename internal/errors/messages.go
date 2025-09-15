package er

// ! pgConn error code status
type pgConnErrorCode string

var (
	PGErr_NOT_NULL         pgConnErrorCode = "23502"
	PGErr_UNIQUE           pgConnErrorCode = "23505"
	PGErr_UNDEFINED_COLUMN pgConnErrorCode = "42703"
)

var PG_ErrorCodes = map[pgConnErrorCode]string{
	PGErr_NOT_NULL:         "NotNull constraint violation",
	PGErr_UNIQUE:           "Unique constraint violation",
	PGErr_UNDEFINED_COLUMN: "Undefined column",
}

// ! validator errors enum
type valErrEnum int

const (
	isRequired valErrEnum = iota
	isEmail
	isMax
	isMin
)

var Validator_ErrMap = map[valErrEnum]string{
	isRequired: "required",
	isEmail:    "email",
	isMax:      "max",
	isMin:      "min",
}

//! error identifier messages - used for matching error strings

type errIdentifier string

const (
	UnknownField errIdentifier = "UnknownField"

	// psConn Errors
	UndefinedCol  errIdentifier = "Undefined column."
	OnValidations errIdentifier = "Error while validating the request JSON data, please check if the data provided is correct."
	DBConflict    errIdentifier = "Data conflicting existing with previous registers. "

	// generic Errors
	QueryTimeout errIdentifier = "Query timeout. Please try again later."
	NotFound     errIdentifier = "Resource not found."
	Conflict     errIdentifier = "Resource already exists."
	Unknown      errIdentifier = "Unknown Error."

	JSONReading errIdentifier = "Error while reading the JSON data. Exceeded the 1mb limit."
)
