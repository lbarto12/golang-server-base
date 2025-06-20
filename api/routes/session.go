package routes

import (
	"encoding/json"
	"golang-server-base/api"
	"golang-server-base/api/apiresponses"
	"golang-server-base/api/postgresapi"
	postgresmodels "golang-server-base/api/postgresapi/models"
	"golang-server-base/api/routes/models"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	var request models.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		apiresponses.Error(w, "invalid Request", http.StatusBadRequest)
		return
	}

	err = postgresapi.SignUp(postgresmodels.Account{
		UserName: request.UserName,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		apiresponses.Error(w, "error creating account", http.StatusInternalServerError)
		return
	}

	token, err := postgresapi.SignIn(postgresmodels.Account{
		Email:    request.Email,
		Password: request.Password,
	}, "")
	if err != nil {
		apiresponses.Error(w, "error signing in", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	apiresponses.Success(w, apiresponses.Response[models.SignUpResponse]{
		Body:    models.SignUpResponse{},
		Message: "successfully signed up",
	})
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	var request models.SignInRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		apiresponses.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	cookie, _ := r.Cookie("session")
	jwtSessionToken := ""
	if cookie != nil {
		jwtSessionToken = cookie.Value
	}

	token, err := postgresapi.SignIn(postgresmodels.Account{
		Email:    request.Email,
		Password: request.Password,
	}, jwtSessionToken)
	if err != nil {
		apiresponses.Error(w, "error signing in", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	apiresponses.Success(w, apiresponses.Response[models.SignInResponse]{
		Body:    models.SignInResponse{},
		Message: "successfully signed in",
	})

}

// Add handlers
func AddSessionHandlers(server *api.Server) {
	server.AddHandlers(map[string]http.Handler{
		"POST /public/api/sign-up": http.HandlerFunc(SignUp),
		"POST /public/api/sign-in": http.HandlerFunc(SignIn),
	})
}
