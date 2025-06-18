package main

import (
	"golang-server-base/api"
	"golang-server-base/api/minioapi"
	"golang-server-base/api/postgresapi"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load Env
	godotenv.Load()

	// Init connections with services
	err := minioapi.Init(minioapi.EnvGetOptions())
	if err != nil {
		panic(err)
	}

	err = postgresapi.Init(postgresapi.EnvGetOptions())
	if err != nil {
		panic(err)
	}

	// Create Server
	apiHost, ok := os.LookupEnv("API_HOST")
	if !ok {
		panic("API_HOST environment variable not set")
	}

	apiPort, ok := os.LookupEnv("API_PORT")
	if !ok {
		panic("API_PORT environment variable not set")
	}

	server := api.Server{
		Host: apiHost,
		Port: apiPort,
		Cors: &cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost},
			AllowCredentials: true,
		},
	}

	// Add Handlers
	server.AddHandlers(map[string]http.Handler{
		"/ping": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		}),
		"/yurr": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("skrt"))
		}),
	})

	// Run Server
	log.Fatal(server.Launch())
}
