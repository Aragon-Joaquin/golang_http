package main

import (
	d "golang-http/internal/dtos"
	er "golang-http/internal/errors"
	"net/http"
)

type UserWithToken struct {
	*d.UserSchema
	Token string `json:"token"`
}

func (s *Server) registerUser(w http.ResponseWriter, r *http.Request) {

	var userInfo d.UserSchema
	if err := s.ReadJSON(w, r, &userInfo); err != nil {
		s.WriteJSONError(w, http.StatusBadRequest, &er.ErrorsStruct{Message: er.JSONReading})
		return
	}

	user, err := s.storage.User.Create(r.Context(), &userInfo)

	if err != nil {
		s.WriteJSONError(w, er.MatchErrorCodes(err.Message), err)
		return
	}

	token, err2 := s.authenticator.GenerateToken(user.Id)

	if err2 != nil {
		s.WriteJSONError(w, http.StatusInternalServerError, &er.ErrorsStruct{Message: err2.Error()})
	}

	userWToken := &UserWithToken{UserSchema: user, Token: token}

	if err := s.WriteJSON(w, http.StatusCreated, userWToken); err != nil {
		s.WriteJSONError(w, http.StatusInternalServerError, &er.ErrorsStruct{Message: err.Error()})
	}

}
