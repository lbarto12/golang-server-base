package src

import "net/http"

func ConfigureRoutes() map[string]http.Handler {
	return map[string]http.Handler{
		"GET /public/ping": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		}),
	}
}
