package systemservices

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

// Enforce Routes
var _ SystemServicesInterface = new(SystemServicesHandlers)

type SystemServicesInterface interface {
	Health(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
}

type SystemServicesHandlers struct {
	Postgres *sqlx.DB
	Minio    *minio.Client
}
