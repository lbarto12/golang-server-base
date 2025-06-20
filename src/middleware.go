package src

import (
	"golang-server-base/api"
	"golang-server-base/api/webtokens"
	"net/http"
)

/*
Required

This is where you can add middleware to your app, a description on how to create middleware
will be in the README.md

Middleware must be injected with a function like so:

	func(next http.Handler) http.Handler {
		...
	}

this allows you to customize your middleware before injecting it into the app.

By default, middleware affects all routes.
*/
func ConfigureMiddleware() []api.Middleware {
	return []api.Middleware{
		func(next http.Handler) http.Handler {
			return webtokens.NewWebTokenMiddleware(next, webtokens.WebTokenMiddleWareConfig{
				PathPrefixExclusions: []string{"/public", "public"},
			})
		},
	}
}
