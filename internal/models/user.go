package models

import (
	"context"
	"database/sql"
	d "golang-http/internal/dtos"

	"github.com/go-playground/validator/v10"
)

var (
	ErrDuplicateEmail    = "duplicated email"
	ErrDuplicateUsername = "duplicated username"
)

type UserStore struct {
	db  *sql.DB
	val *validator.Validate
}

// ! create user
type createUser struct {
	Username string `validate:"required,max=20,min=3" name:"username"`
	Email    string `validate:"required,email,max=50,min=3" name:"email"`
}

func (s *UserStore) Create(ctx context.Context, user *d.UserSchema) (*d.UserSchema, *ErrorsStruct) {
	ctx, cancel := context.WithTimeout(ctx, ContextMaxTimeout)
	defer cancel()

	userCreate := createUser{Username: user.Username, Email: user.Email}
	err := s.val.Struct(userCreate)

	if err != nil {
		return nil, ValidatorErrorParser(err)
	}

	query := `
		SELECT * FROM CreateUser($1, $2);
	`
	args := []any{userCreate.Username, userCreate.Email}
	var newUser d.UserSchema

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&newUser.Id, &newUser.Username, &newUser.Email)

	if err != nil {
		return nil, CheckForGenericErrors(err)
	}

	return &d.UserSchema{Id: newUser.Id, Username: newUser.Username, Email: newUser.Email}, nil
}

//! get user
