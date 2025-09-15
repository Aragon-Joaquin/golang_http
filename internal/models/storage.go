package models

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	d "golang-http/internal/dtos"

	"github.com/go-playground/validator/v10"
)

var (
	ContextMaxTimeout = 3 * time.Second
)

type ModelsStorageStruct struct {
	User interface {
		Create(context.Context, *d.UserSchema) (*d.UserSchema, *ErrorsStruct)
		Get(context.Context, int64) (*d.UserSchema, *ErrorsStruct)
	}
}

type ErrorsStruct struct {
	Validations map[string]string `json:"validationsErrors,omitempty"`
	Message     string            `json:"message"`
}

func ModelsStorage(db *sql.DB) *ModelsStorageStruct {
	validate := validator.New(validator.WithRequiredStructEnabled())

	// uses the "name:''" tag in the struct as the name field
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("name")
	})

	return &ModelsStorageStruct{
		User: &UserStore{db: db, val: validate},
	}
}
