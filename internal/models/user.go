package models

import (
	"context"
	"database/sql"
	d "golang-http/internal/dtos"
	e "golang-http/internal/errors"

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
type CreateUser struct {
	Username string `validate:"required,max=20,min=3" name:"username"`
	Email    string `validate:"required,email,max=50,min=3" name:"email"`
}

func (s *UserStore) Create(ctx context.Context, user *d.UserSchema) (*d.UserSchema, *e.ErrorsStruct) {
	ctx, cancel := context.WithTimeout(ctx, ContextMaxTimeout)
	defer cancel()

	userCreate := CreateUser{Username: user.Username, Email: user.Email}

	if err := s.val.Struct(userCreate); err != nil {
		return nil, e.ValidatorErrorParser(err)
	}

	query := `
		SELECT id, username, email, created_at FROM CreateUser($1, $2);
	`
	args := []any{userCreate.Username, userCreate.Email}
	var newUser d.UserSchema

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&newUser.Id, &newUser.Username, &newUser.Email, &newUser.Created_at)

	if err != nil {
		return nil, e.CheckForGenericErrors(err)
	}

	return &d.UserSchema{Id: newUser.Id, Username: newUser.Username, Email: newUser.Email}, nil
}

// ! get user
type GetUser struct {
	Id int64 `validate:"required,number,gt=0" name:"id"`
}

func (s *UserStore) Get(ctx context.Context, id int64) (*d.UserSchema, *e.ErrorsStruct) {
	ctx, cancel := context.WithTimeout(ctx, ContextMaxTimeout)
	defer cancel()

	userGet := GetUser{Id: id}

	if err := s.val.Struct(userGet); err != nil {
		return nil, e.ValidatorErrorParser(err)
	}

	query := `
		SELECT id, username, email, created_at FROM users WHERE id = $1;
	`
	args := []any{userGet.Id}
	var newUser d.UserSchema
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&newUser.Id, &newUser.Username, &newUser.Email, &newUser.Created_at)

	if err != nil {
		return nil, e.CheckForGenericErrors(err)
	}

	return &d.UserSchema{Id: newUser.Id, Username: newUser.Username, Email: newUser.Email, Created_at: newUser.Created_at}, nil
}

// ! get by credentials
func (s *UserStore) GetByCredentials(ctx context.Context, user *d.UserSchema) (*d.UserSchema, *e.ErrorsStruct) {
	ctx, cancel := context.WithTimeout(ctx, ContextMaxTimeout)
	defer cancel()

	userCredential := CreateUser{Username: user.Username, Email: user.Email}

	if err := s.val.Struct(userCredential); err != nil {
		return nil, e.ValidatorErrorParser(err)
	}

	query := `
		SELECT id, username, email, created_at FROM users WHERE username = $1 AND email = $2;
	`
	args := []any{userCredential.Username, userCredential.Email}
	var newUser d.UserSchema

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&newUser.Id, &newUser.Username, &newUser.Email, &newUser.Created_at)

	if err != nil {
		return nil, e.CheckForGenericErrors(err)
	}

	return &d.UserSchema{Id: newUser.Id, Username: newUser.Username, Email: newUser.Email}, nil
}
