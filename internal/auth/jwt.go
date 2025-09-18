package auth

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtKey = os.Getenv("JWT_KEY")
	exp    = time.Hour * 24 * 2 // 2 days
	iss    = "golang_http"      // identifies the source
)

type JWTAuthenticator struct {
	secret string //! key
	aud    string //! audience
	// (identifies the recipients that the JWT is intended for.)

	iss string //! issuer
	// identifies the principal that issued the JWT (can be a human user, an organization, or a service.)

	//? both iss and aud prevents token abuse, and provides some basic information about the context in which was issued.
}

// ! for creating the new instance
func NewJWTAuthenticator() *JWTAuthenticator {
	return &JWTAuthenticator{secret: jwtKey, aud: iss, iss: iss}
}

// ! generates tokens
type JWTCore struct {
	Token        *jwt.Token
	SignedString string
}

func (j *JWTAuthenticator) GenerateToken(identifier any) (string, error) {

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		//*JWT specification (RFC 7519) FIELDS
		"sub": identifier,                 //subject: identifies the users uniquely (this needs to be an userID or similar)
		"exp": time.Now().Add(exp).Unix(), //expiration: when the token expires
		"iat": time.Now().Unix(),          //issued at: when the token was issued/created
		"nbf": time.Now().Unix(),          //not before: the time before which the token is not accepted for processing
		"iss": iss,                        //issuer: identifies the source (who made the jwt)
		"aud": iss,                        //audience: identifies uniquely the consumer/receiver (only those services specified can use this jwt)
	})

	signedToken, err := t.SignedString([]byte(j.secret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ! validates tokens
func (j *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	tokenString, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(j.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(j.aud),
		jwt.WithIssuer(j.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	if err != nil {
		switch strings.Split(err.Error(), ":")[0] {
		case jwt.ErrTokenMalformed.Error():
			return nil, errors.New("jwt unable to parse, invalid jwt interface")
		case jwt.ErrTokenExpired.Error():
			return nil, errors.New("jwt is expired. please refresh your token again")
		default:
			return nil, err
		}
	}

	return tokenString, nil
}
