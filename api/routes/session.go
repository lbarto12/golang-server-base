package routes

import (
	"golang-server-base/api"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sign Up!"))
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sign In!"))
}

// Add handlers
func AddSessionHandlers(server *api.Server) {
	server.AddHandlers(map[string]http.Handler{
		"POST /api/sign-up": http.HandlerFunc(SignUp),
		"POST /api/sign-in": http.HandlerFunc(SignIn),
	})
}
