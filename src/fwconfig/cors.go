package fwconfig

import (
	"net/http"

	"github.com/rs/cors"
)

/*
Required

This method allows you to configure cors from here.
The *cors.Options object that is returned is used as
the settings when the server is created.

If you do not want cors in your application, then return nil from this function.
*/
func ConfigureCors() *cors.Options {
	return &cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowCredentials: true,
	}
}
