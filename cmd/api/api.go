package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", s.HealthHandler)

		//  routes for /v1/user
		r.Route("/user", func(r chi.Router) {
			r.Get("/create", s.createUser)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		s.WriteJSONError(w, http.StatusNotFound, "Please specify the version of the api.")
	})

	return r
}
