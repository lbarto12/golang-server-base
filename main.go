package main

import (
	"golang-server-base/api"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {

	server := api.Server{
		Host: "localhost",
		Port: "8080",
		Cors: &cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost},
			AllowCredentials: true,
		},
	}

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

	log.Fatal(server.Launch())
}
