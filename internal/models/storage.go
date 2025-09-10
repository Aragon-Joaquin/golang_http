package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type ModelsStorageStruct struct {
	User interface {
		Create(context.Context, *UserSchema) (*UserSchema, error)
	}
}

func ModelsStorage(db *sql.DB) *ModelsStorageStruct {
	return &ModelsStorageStruct{
		User: &UserStore{db: db},
	}
}
