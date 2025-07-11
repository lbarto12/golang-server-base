package main

import (
	"golang-server-base/api"
	"golang-server-base/api/apiservices"
	"golang-server-base/api/emailapi"
	"golang-server-base/api/meilisearchapi"
	"golang-server-base/api/minioapi"
	"golang-server-base/api/postgresapi"
	"golang-server-base/api/routes/systemservices"
	"golang-server-base/api/webtokensapi"
	"golang-server-base/src"
	"golang-server-base/src/fwconfig"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/joho/godotenv"
)

func main() {
	// Load Env
	// godotenv.Load()

	// Load Env -> Handle Errors
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .envx file:", envErr)
	}

	enabledServices := fwconfig.ConfigureServices()

	// Init connections with services

	serviceAccess := apiservices.ServicesAccess{}
	var err error

	// Init Postgres
	if slices.Contains(enabledServices, apiservices.Postgres) {
		err = postgresapi.Init(postgresapi.EnvGetOptions())
		if err != nil {
			panic(err)
		}
		db, err := postgresapi.Database()
		if err != nil {
			panic(err)
		}
		serviceAccess.Postgres = db
	}

	// Init Minio
	if slices.Contains(enabledServices, apiservices.Minio) {
		err := minioapi.Init(minioapi.EnvGetOptions())
		if err != nil {
			panic(err)
		}
		client := minioapi.Client()
		serviceAccess.Minio = client
	}

	// Init JWT library
	if slices.Contains(enabledServices, apiservices.Webtokens) {
		err = webtokensapi.Init()
		if err != nil {
			panic(err)
		}
	}

	// Init SMTP
	if slices.Contains(enabledServices, apiservices.Email) {
		err = emailapi.Init(emailapi.EnvGetOptions())
		if err != nil {
			panic(err)
		}
		serviceAccess.Email = emailapi.Dialer()
	}

	// Init Meilisearch
	if slices.Contains(enabledServices, apiservices.Meilisearch) {
		err = meilisearchapi.Init(meilisearchapi.EnvGetOptions())
		if err != nil {
			panic(err)
		}
		client, err := meilisearchapi.Client()
		if err != nil {
			panic(err)
		}
		serviceAccess.Meilisearch = &client
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

	server := api.NewServer(api.Server{
		Host: apiHost,
		Port: apiPort,
	})

	// Add services
	server.Services = serviceAccess

	// Add Handlers ============

	// Builtin

	systemhandlers := systemservices.SystemServicesHandlers{
		Postgres: server.Services.Postgres,
		Minio:    server.Services.Minio,
	}
	server.AddHandler("GET /public/api/health", http.HandlerFunc(systemhandlers.Health))

	if slices.Contains(enabledServices, apiservices.Sessions) {
		server.AddHandlers(map[string]http.Handler{
			"POST /public/api/sign-up": http.HandlerFunc(systemhandlers.SignUp),
			"POST /public/api/sign-in": http.HandlerFunc(systemhandlers.SignIn),
		})
	}

	// User defined

	// Add from `routes.go` in `src`
	server.AddHandlers(fwconfig.ConfigureRoutes(&server.Services))

	// =========================

	// Add Middleware
	server.AddMiddleWares(fwconfig.ConfigureMiddleware())

	// Set user defined cors
	server.Cors = fwconfig.ConfigureCors()

	// Run user code
	src.Main(&server)

	// Run Server
	log.Fatal(server.Launch())
}
