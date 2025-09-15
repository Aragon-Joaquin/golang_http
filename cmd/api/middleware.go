package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.WriteJSONError(w, http.StatusUnauthorized, "missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			s.WriteJSONError(w, http.StatusUnauthorized, "invalid Authorization header format, should be 'Bearer <token>")
			return
		}

		token := parts[1]
		jwtToken, err := s.authenticator.ValidateToken(token)
		if err != nil {
			fmt.Println(err)
			s.WriteJSONError(w, http.StatusUnauthorized, err.Error())
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		//transforms it into from a string > to a float > to a base10 64 int

		if err != nil {
			s.WriteJSONError(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := r.Context()
		user, err2 := s.storage.User.Get(ctx, userID)

		if err2 != nil {
			s.WriteJSONError(w, http.StatusUnauthorized, err2.Message)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
