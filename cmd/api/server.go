package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"golang-http/internal/auth"
	"golang-http/internal/database"
	"golang-http/internal/models"
)

type Server struct {
	port int

	storage       models.ModelsStorageStruct
	db            *database.Service
	authenticator auth.Authenticator
}

func newServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	databaseServ := database.New()
	NewServer := &Server{
		port:          port,
		storage:       *models.ModelsStorage(databaseServ.Db),
		db:            databaseServ,
		authenticator: auth.NewJWTAuthenticator(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
