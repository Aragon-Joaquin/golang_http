package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserSchema struct {
	Id         int64     `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Created_at time.Time `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *UserSchema) (*UserSchema, error) {
	query := `
		SELECT * FROM CreateUser($1, $2);
	`
	args := []any{user.Username, user.Email}

	var newUser UserSchema
	// 	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	// defer cancel()

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&newUser.Id, &newUser.Username, &newUser.Email)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	return &UserSchema{Id: newUser.Id, Username: newUser.Username, Email: newUser.Email}, nil
}
