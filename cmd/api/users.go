package main

import (
	er "golang-http/internal/errors"
	"net/http"
	"strconv"

	d "golang-http/internal/dtos"

	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "userIdentifier"

// GetUser godoc
//
//	@Router			/user/{id} [get]
//	@Summary		Fetches a user profile
//	@Version		1.0
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Failure		400	{object}	DataResponse[er.ReturnedError]
//	@Failure		401	{object}	DataResponse[er.ReturnedError]
//	@Failure		404	{object}	DataResponse[er.ReturnedError]
//	@Failure		500	{object}	DataResponse[er.ReturnedError]
//	@Success		200	{object}	DataResponse[d.UserSchema]
//	@Security		ApiKeyAuth
func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil || id <= 0 {
		s.WriteJSONError(w, http.StatusBadRequest, &er.ErrorsStruct{Message: er.OnValidations})
		return
	}

	user, err2 := s.storage.User.Get(r.Context(), id)

	if err2 != nil {
		s.WriteJSONError(w, er.MatchErrorCodes(err2.Message), err2)
		return
	}

	if err := s.WriteJSON(w, http.StatusOK, user); err != nil {
		s.WriteJSONError(w, http.StatusInternalServerError, &er.ErrorsStruct{Message: err.Error()})
	}

}

// getOwnUser godoc
//
//	@Router			/user [get]
//	@Summary		Fetches your own profile
//	@Version		1.0
//	@Description	Returns information about the same user if its authenticated
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Failure		400	{object}	DataResponse[er.ReturnedError]
//	@Failure		401	{object}	DataResponse[er.ReturnedError]
//	@Failure		404	{object}	DataResponse[er.ReturnedError]
//	@Failure		500	{object}	DataResponse[er.ReturnedError]
//	@Success		200	{object}	DataResponse[d.UserSchema]
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
