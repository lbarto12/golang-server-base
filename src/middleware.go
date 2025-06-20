package src

import (
	"golang-server-base/api"
	"golang-server-base/api/webtokens"
	"net/http"
)

func CongiureMiddleware() []api.Middleware {
	return []api.Middleware{
		func(next http.Handler) http.Handler {
			return webtokens.NewWebTokenMiddleware(next, webtokens.WebTokenMiddleWareConfig{
				PathPrefixExclusions: []string{"/public", "public"},
			})
		},
	}
}
