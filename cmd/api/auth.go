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

// Register User godoc
//
//	@Router			/auth/user [post]
//	@Summary		Creates an user and returns JWT
//	@Version		1.0
//	@Description	Creates an user with the fields "email" and "username" and returns a jwt if its valid
//	@Tags			users, auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.CreateUser	true	"User data"
//	@Failure		400		{object}	DataResponse[er.ReturnedError]
//	@Failure		401		{object}	DataResponse[er.ReturnedError]
//	@Failure		404		{object}	DataResponse[er.ReturnedError]
//	@Failure		500		{object}	DataResponse[er.ReturnedError]
//	@Success		201		{object}	DataResponse[UserWithToken]
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

// LogInUser godoc
//
//	@Router			/auth/token [post]
//	@Summary		Generates a JWT based on user credentials
//	@Version		1.0
//	@Description	Creates a JWT with the fields "email" and "username" and returns its if its valid
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.CreateUser	true	"User data"
//	@Failure		400		{object}	DataResponse[er.ReturnedError]
//	@Failure		401		{object}	DataResponse[er.ReturnedError]
//	@Failure		404		{object}	DataResponse[er.ReturnedError]
//	@Failure		500		{object}	DataResponse[er.ReturnedError]
//	@Success		200		{object}	DataResponse[UserWithToken]
func (s *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	var userInfo d.UserSchema
	if err := s.ReadJSON(w, r, &userInfo); err != nil {
		s.WriteJSONError(w, http.StatusBadRequest, &er.ErrorsStruct{Message: er.JSONReading})
		return
	}

	user, err := s.storage.User.GetByCredentials(r.Context(), &userInfo)

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
