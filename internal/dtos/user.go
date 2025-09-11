package handlers

import "time"

type UserSchema struct {
	Id         int64     `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Created_at time.Time `json:"created_at"`
}
