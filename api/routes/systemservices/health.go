package routes

import (
	"net/http"
)

func (handler SystemServicesHandlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
