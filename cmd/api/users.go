package main

import (
	d "golang-http/internal/dtos"
	er "golang-http/internal/errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "userIdentifier"

// todo: finish
func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil || id <= 0 {
		s.WriteJSONError(w, http.StatusBadRequest, &er.ErrorsStruct{Message: er.OnValidations})
		return
	}

}

func getUserFromContext(r *http.Request) *d.UserSchema {
	user, _ := r.Context().Value(userCtx).(*d.UserSchema)
	return user
}
