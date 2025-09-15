package main

import (
	"fmt"
	d "golang-http/internal/dtos"
	er "golang-http/internal/errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "userIdentifier"

// ! get user
func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil || id <= 0 {
		s.WriteJSONError(w, http.StatusBadRequest, &er.ErrorsStruct{Message: er.OnValidations})
		return
	}

	user, err2 := s.storage.User.Get(r.Context(), id)

	if err2 != nil {
		fmt.Println(err2)
		s.WriteJSONError(w, er.MatchErrorCodes(err2.Message), err2)
		return
	}

	if err := s.WriteJSON(w, http.StatusOK, user); err != nil {
		s.WriteJSONError(w, http.StatusInternalServerError, &er.ErrorsStruct{Message: err.Error()})
	}

}

// ! get own user
func (s *Server) getOwnUser(w http.ResponseWriter, r *http.Request) {
	u := s.getUserFromContext(r)

	user, err := s.storage.User.Get(r.Context(), u.Id)

	if err != nil {
		s.WriteJSONError(w, er.MatchErrorCodes(err.Message), err)
		return
	}

	if err := s.WriteJSON(w, http.StatusOK, user); err != nil {
		s.WriteJSONError(w, http.StatusInternalServerError, &er.ErrorsStruct{Message: err.Error()})
	}
}

func (s *Server) getUserFromContext(r *http.Request) *d.UserSchema {
	user, _ := r.Context().Value(userCtx).(*d.UserSchema)
	return user
}
