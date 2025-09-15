package main

import (
	er "golang-http/internal/errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	errorSpecifyAPIVer = "Please specify the version of the api."
	incorrectPath      = "The path provided is not valid. Maybe it could be a wrong HTTP Method?"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) //adds a unique counter for each request
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recovers from panics, logs the panic, and returns a 500 status if possible

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
			r.Use(s.AuthMiddleware)

			r.Get("/{id}", s.getUser)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", s.registerUser)
			// r.Post("/token", func(w http.ResponseWriter, r *http.Request) {})
		})

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			s.WriteJSONError(w, http.StatusNotFound, &er.ErrorsStruct{Message: incorrectPath})
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		s.WriteJSONError(w, http.StatusNotFound, &er.ErrorsStruct{Message: errorSpecifyAPIVer})
	})

	return r
}
