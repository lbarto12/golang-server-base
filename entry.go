package main

import (
	"golang-server-base/api"
	"golang-server-base/api/emailapi"
	"golang-server-base/api/meilisearchapi"
	"golang-server-base/api/minioapi"
	"golang-server-base/api/postgresapi"
	routes "golang-server-base/api/routes/systemservices"
	"golang-server-base/api/webtokensapi"
	"golang-server-base/src"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

	// Init JWT library
	err = webtokensapi.Init()
	if err != nil {
		panic(err)
	}

	// Init SMTP
	err = emailapi.Init(emailapi.EnvGetOptions())
	if err != nil {
		panic(err)
	}

	// Init Meilisearch
	meilisearchapi.Init(meilisearchapi.EnvGetOptions())

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

	// Add Handlers ============

	// Builtin
	systemhandlers := routes.SystemServicesHandlers{}

	server.AddHandlers(map[string]http.Handler{
		"GET /api/health":          http.HandlerFunc(systemhandlers.Health),
		"POST /public/api/sign-up": http.HandlerFunc(systemhandlers.SignUp),
		"POST /public/api/sign-in": http.HandlerFunc(systemhandlers.SignIn),
	})

	// User defined
	// Add from `routes.go` in `src`
	server.AddHandlers(src.ConfigureRoutes())

	// =========================

	// Add Built-in Middleware
	server.AddMiddleWares(src.ConfigureMiddleware())

	// Set user defined cors
	server.Cors = src.ConfigureCors()

	// Run user code
	src.Main(&server)

	// Run Server
	log.Fatal(server.Launch())
}
