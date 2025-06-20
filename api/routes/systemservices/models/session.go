package models

type SignUpRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
}

type SignInRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type SignInResponse struct {
}
