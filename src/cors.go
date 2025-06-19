package src

import (
	"net/http"

	"github.com/rs/cors"
)

func ConfigureCors() *cors.Options { // Required, if you do not want cors, return nil from this method
	return &cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowCredentials: true,
	}
}
