package models

import (
	"time"
)

type User struct {
	Username   string
	Email      string
	Created_at time.Time
}
