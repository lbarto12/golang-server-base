package apiresponses

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Body    T      `json:"body"`
	Message string `json:"message,omitempty"`
}

func Error(w http.ResponseWriter, message string, status int) {
	response := Response[any]{
		Body:    nil,
		Message: message,
	}
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusInternalServerError)
	}
	http.Error(w, string(res), status)
}

func Success[T any](w http.ResponseWriter, response Response[T]) {
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
