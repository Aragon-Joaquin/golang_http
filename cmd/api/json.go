package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Message any  `json:"message"`
	Error   bool `json:"error"`
}

// general func
func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		writeJSON(w, http.StatusInternalServerError, &response{Message: err.Error(), Error: true})
		log.Fatalf("error handling JSON marshal. Err: %v", err)
		return err
	}

	return nil
}

// methods for server.
func (s *Server) WriteJSONError(w http.ResponseWriter, status int, message any) error {
	return writeJSON(w, status, &response{Message: message, Error: true})
}

func (s *Server) WriteJSON(w http.ResponseWriter, status int, message string) error {
	return writeJSON(w, status, &response{Message: message, Error: false})
}

func (s *Server) WriteJSONDataField(w http.ResponseWriter, status int, data any) error {
	type responseData struct {
		Data  any  `json:"data"`
		Error bool `json:"error"`
	}

	return writeJSON(w, status, &responseData{Data: data, Error: false})
}

// decode/read
func (s *Server) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // ~1mb, underscores only improves readability and have no impact whatsoever
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	//! need to return the readed data as well
	return decoder.Decode(data)
}
