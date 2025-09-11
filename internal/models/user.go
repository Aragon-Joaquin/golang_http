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
	Username string `json:"username" validate:"required,max=20,min=3"`
	Email    string `json:"email" validate:"required,email,max=50,min=3"`
}

func (s *UserStore) Create(ctx context.Context, user *d.UserSchema) (*d.UserSchema, *ErrorsStruct) {
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

	// 	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	// defer cancel()

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&newUser.Id, &newUser.Username, &newUser.Email)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "email"`:
			return nil, &ErrorsStruct{Message: ErrDuplicateEmail}
		case err.Error() == `pq: duplicate key value violates unique constraint "username"`:
			return nil, &ErrorsStruct{Message: ErrDuplicateUsername}
		default:
			return nil, &ErrorsStruct{Message: err.Error()}
		}
	}

	return &d.UserSchema{Id: newUser.Id, Username: newUser.Username, Email: newUser.Email}, nil
}

//! get user
