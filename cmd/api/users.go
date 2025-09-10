package main

import (
	"golang-http/internal/models"
	"net/http"
)

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {

	var userInfo models.UserSchema
	s.ReadJSON(w, r, userInfo)

	user, err := s.storage.User.Create(r.Context(), &userInfo)

	if err != nil {
		switch err {
		case models.ErrNotFound:
			s.WriteJSONError(w, http.StatusNotFound, err.Error())
		default:
			s.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		}
	}

	if err := s.WriteJSONDataField(w, http.StatusCreated, user); err != nil {
		s.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
