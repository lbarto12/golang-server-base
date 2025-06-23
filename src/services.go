package src

import (
	"golang-server-base/api"
	"golang-server-base/api/apiservices"
)

func ConfigureServices() []api.Service {
	return []api.Service{
		apiservices.Postgres,
		apiservices.Minio,
		apiservices.Meilisearch,
		apiservices.Email,
		apiservices.Sessions,
		apiservices.Webtokens,
	}
}
