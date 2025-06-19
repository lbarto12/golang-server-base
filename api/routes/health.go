package routes

import (
	"golang-server-base/api"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func AddHealthHandlers(server *api.Server) {
	server.AddHandlers(map[string]http.Handler{
		"GET /api/health": http.HandlerFunc(Health),
	})
}
