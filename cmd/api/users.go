package main

import (
	d "golang-http/internal/dtos"
	"golang-http/internal/models"
	"net/http"
)

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {

	//TODO: this wont return the json data
	var userInfo d.UserSchema
	s.ReadJSON(w, r, userInfo)

	user, err := s.storage.User.Create(r.Context(), &userInfo)

	if err != nil {
		switch err.Message {
		case models.ErrNotFound.Error():
			s.WriteJSONError(w, http.StatusNotFound, err)
		default:
			s.WriteJSONError(w, http.StatusInternalServerError, err)
		}
		return
	}

	if err := s.WriteJSONDataField(w, http.StatusCreated, user); err != nil {
		s.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
