package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	d "golang-http/internal/dtos"

	"github.com/go-playground/validator/v10"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type ModelsStorageStruct struct {
	User interface {
		Create(context.Context, *d.UserSchema) (*d.UserSchema, *ErrorsStruct)
	}
}

type ErrorsStruct struct {
	Validations map[string]string `json:"validationsErrors,omitempty"`
	Message     string            `json:"message"`
}

func ModelsStorage(db *sql.DB) *ModelsStorageStruct {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &ModelsStorageStruct{
		User: &UserStore{db: db, val: validate},
	}
}
