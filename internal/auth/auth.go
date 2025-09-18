package auth

import "github.com/golang-jwt/jwt/v5"

type Authenticator interface {
	GenerateToken(identifier any) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
