package webtokensapi

import (
	"net/http"
	"strings"
)


// JWT -> Authentication 
// JWT -> {Roles, Id, Claims} -> Authorization 


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

	mw.next.ServeHTTP(w, r)
}
