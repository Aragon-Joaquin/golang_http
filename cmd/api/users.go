package main

import (
	d "golang-http/internal/dtos"
	"golang-http/internal/models"
	"net/http"
)

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {

	var userInfo d.UserSchema
	if err := s.ReadJSON(w, r, &userInfo); err != nil {
		s.WriteJSONError(w, http.StatusBadRequest, models.ErrMsg_JSONReading)
		return
	}

	user, err := s.storage.User.Create(r.Context(), &userInfo)

	if err != nil {
		switch err.Message {
		case models.ErrMsg_QueryTimeout:
			s.WriteJSONError(w, http.StatusRequestTimeout, err)
		case models.ErrMsg_NotFound:
			s.WriteJSONError(w, http.StatusNotFound, err)
		case models.ErrMsg_DBConflict:
			s.WriteJSONError(w, http.StatusConflict, err)
		case models.ErrMsg_OnValidations:
			s.WriteJSONError(w, http.StatusBadRequest, err)
		case models.ErrMsg_UndefinedCol:
			s.WriteJSONError(w, http.StatusBadRequest, err)
		default:
			s.WriteJSONError(w, http.StatusNotImplemented, err)
		}
		return
	}

	if err := s.WriteJSON(w, http.StatusCreated, user); err != nil {
		s.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
