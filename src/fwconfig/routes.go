package fwconfig

import (
	"golang-server-base/api/apiservices"
	"net/http"
)

/*
required

Here is where you can add your API's routes.

Do not use anonymous functions as handler; instead use idomatic ways of defining your routes and include them here

These routes will be automatically wrapped in any applied middleware from `middleware.go`
*/
func ConfigureRoutes(services *apiservices.ServicesAccess) map[string]http.Handler {
	return map[string]http.Handler{
		"GET /public/ping": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		}),
	}
}
