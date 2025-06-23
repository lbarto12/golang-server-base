package src

import (
	"golang-server-base/api"
	"golang-server-base/api/apiservices"
)

/*
required

Here you can add and remove services provided by this framework.
To view all available services, see `golang-server-base/apiservices`, imported above

Note:
Any services that rely on docker containers should also be omitted from the docker compose startup script,

Example, disable postgres:
`docker compose up minio mielisearch`
*/
func ConfigureServices() []api.Service {
	return []api.Service{
		apiservices.Postgres,
		apiservices.Minio,
		apiservices.Email,
		apiservices.Meilisearch,
		apiservices.Sessions,
		apiservices.Webtokens,
	}
}
