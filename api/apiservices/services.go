package apiservices

import (
	"github.com/jmoiron/sqlx"
	"github.com/meilisearch/meilisearch-go"
	"github.com/minio/minio-go/v7"
	"gopkg.in/gomail.v2"
)

type ServicesAccess struct {
	Postgres    *sqlx.DB
	Minio       *minio.Client
	Email       *gomail.Dialer
	Meilisearch *meilisearch.ServiceManager
}

// Available Services
const (
	Postgres    = iota
	Minio       = iota
	Email       = iota
	Meilisearch = iota
	Sessions    = iota
	Webtokens   = iota
)
