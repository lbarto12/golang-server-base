package webtokens

import (
	"net/http"
	"strings"
)

type WebTokenMiddleWare struct {
	next   http.Handler
	config WebTokenMiddleWareConfig
}

type WebTokenMiddleWareConfig struct {
	PathPrefixExclusions []string
}

func NewWebTokenMiddleware(next http.Handler, config WebTokenMiddleWareConfig) *WebTokenMiddleWare {
	return &WebTokenMiddleWare{
		next,
		config,
	}
}

func (mw WebTokenMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(strings.ToLower(r.URL.Path))

	// Exclude middeware from paths
	if mw.config.PathPrefixExclusions != nil {
		for _, exclusion := range mw.config.PathPrefixExclusions {
			if strings.HasPrefix(path, exclusion) {
				mw.next.ServeHTTP(w, r)
				return
			}
		}
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := VerifyToken(cookie.Value)
	if err != nil || !token.Valid {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}

	//TODO: Validate User? the above just validates that the JWT was made from our API

	mw.next.ServeHTTP(w, r)
}
