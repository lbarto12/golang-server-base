package main

import (
	"golang-server-base/api"
	"golang-server-base/api/minioapi"
	"golang-server-base/api/postgresapi"
	"golang-server-base/api/routes"
	"golang-server-base/api/webtokens"
	"golang-server-base/src"
	"log"
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
	err = webtokens.Init()
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
	}

	// Add Built-In Handlers ============

	// Add Health
	routes.AddHealthHandlers(&server)

	// Add Session
	routes.AddSessionHandlers(&server)

	// =========================

	// Set user defined cors
	server.Cors = src.ConfigureCors()

	// Run user code
	src.Main(&server)

	// Run Server
	log.Fatal(server.Launch())
}
