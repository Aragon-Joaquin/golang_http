package main

import (
	"encoding/json"
	err "golang-http/internal/errors"
	"log"
	"net/http"
)

type DataResponse[T any] struct {
	Data  T    `json:"data"`
	Error bool `json:"error"`
}

// general func
func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		writeJSON(w, http.StatusInternalServerError, &DataResponse[any]{Data: err.Error(), Error: true})
		log.Fatalf("error handling JSON marshal. Err: %v", err)
		return err
	}

	return nil
}

// methods for server.

func (s *Server) WriteJSONError(w http.ResponseWriter, status int, message *err.ErrorsStruct) error {
	return writeJSON(w, status, &DataResponse[*err.ErrorsStruct]{Data: message, Error: true})
}

func (s *Server) WriteJSON(w http.ResponseWriter, status int, message any) error {
	return writeJSON(w, status, &DataResponse[any]{Data: message, Error: false})
}

// decode/read
func (s *Server) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // ~1mb, underscores only improves readability and have no impact whatsoever
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}
