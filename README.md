# What is this project?
This project is a framework written in go that has a variety prebuilt services, libraries, and driver code common 
to many apps being built.

- Current services and libraries include:
    - Postgres
    - Minio
    - Meilisearch
    - Email
    - PDF Generation    <---- NOT IMPLEMENTED
    - Cors

- Current features
    - Session Handling
        - Managed by a database table called 'accounts' in postgres
        - Accessed by the endpoints `"POST /api/sign-up"` and `"POST /api/sign-in"`
    - API Health checks
        - Accessed by endpoint(s): `"GET /api/health"`
    - Managed Service Clients
        - Services are accessed via the package `<servicename>api`
            - e.g. `db, err := postgresapi.Database()`
            - e.g. `client, err := minioapi.Client()`
			- e.g. `err := emailapi.Send(...)`


# Installation
To install, simply clone the repo and run `go mod tidy` or `go mod download` to install all required packages

To customize the schema of your postgres instance, modify the `init.sql` file located in `<project-root>/api/init/postgres/`

Run the bash script to generate and initialize the services:
```bash
./rebuild-docker.sh
```
This will also launch the services with `docker compose up`



### Example .env

```.env
# API
API_HOST=localhost
API_PORT=8080

# session
JWT_SECRET_KEY=my_secret_jwt_key

# postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5050
POSTGRES_DATABASE=postgres-database-name
POSTGRES_USER=postgres_user
POSTGRES_PASSWORD=postgres_password
POSTGRES_MAX_OPEN_CONNECTIONS=75

# minio
MINIO_ENDPOINT=localhost:9000
MINIO_USER=minio_user
MINIO_PASSWORD=minio_password
MINIO_DEFAULT_BUCKET=default_bucket
MINIO_USE_SSL=false

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_SENDER=youremail
SMTP_PASSWORD=yourpassword

# meilisearch
MEILI_HOST=http://localhost
MEILI_PORT=7700
MEILI_MASTER_KEY=my_secret_meili_master_key
```

# How to run

In the root directory, follow these steps:

Launch Postrgres and Minio
```bash
docker compose up
```

Launch the server
```bash
go run .
```

# Writing your code

To write programs using this framework, keep all of your code within the `src` directory.
Within this directory are a few files:

### main.go
This is where you will put your code, a server object is injected for you to modify its various properties before it is launched.

```go
package src

import (
	"golang-server-base/api"
)


func Main(server *api.Server) {
	// Run your code here
}
```


### cors.go
This is where you can configure cors easily, the return value is a `*cors.Options` object. Its configuration will be applied to the server.
If you do not wish to have cors in your application, then return `nil` from this function.

```go
package src

import (
	"net/http"

	"github.com/rs/cors"
)

func ConfigureCors() *cors.Options {
	return &cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowCredentials: true,
	}
}
```


### routes.go
This is where you can add your custom routes to the application.

Routes that you want to add should be specified as a key value pair in the returned `map`, with the key being the path of the endpoint, and the value
being the function itself.

It is `not reccomended` that you pass in anonymous functions as the value for your routes, and instead define them in their own packages. The following
does so for the sake of demonstration.

Here is how you might add a `ping` endpoint to your application:

```go
package src

import "net/http"

func ConfigureRoutes() map[string]http.Handler {
	return map[string]http.Handler{
		"GET /public/ping": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		}),
	}
}
```


### middleware.go
This is where you can add middleware to your application's endpoints.

To create new middleware, it must be designed as a struct that implements the `ServeHTTP` method from the `http.Handler` interface. 

e.g.
```go
type Middleware struct {
	next   http.Handler
	config MiddlewareConfig
}

type MiddlewareConfig struct {
	PathPrefixExclusions []string
}

func NewMiddleware(next http.Handler, config MiddlewareConfig) *Middleware {
	return &Middleware{
		next,
		config,
	}
}

func (mw Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(strings.ToLower(r.URL.Path))

	// Example: Exclude middeware from paths
	if mw.config.PathPrefixExclusions != nil {
		for _, exclusion := range mw.config.PathPrefixExclusions {
			if strings.HasPrefix(path, exclusion) {
				mw.next.ServeHTTP(w, r)
				return
			}
		}
	}

    // Your middleware's functionality

	mw.next.ServeHTTP(w, r)
}
```

This is because, under the hood, the app will repeatedly wrap the `*http.ServeMux` in the next included middleware until it has applied all of it.

In the `middleware.go` file, you may include any preconfigured or user made middleware so long as it matches the specificiations. To include middleware,
it must be written as a function that accepts, and returns, an `http.Handler` object. This allows for configurations prior to adding the middleware, and
lets each middleware propogate to the next.

Example of including the above middleware in your project, with settings:

```go
package src

import (
	"golang-server-base/api"
	"golang-server-base/api/webtokens"
	"net/http"
)

func ConfigureMiddleware() []api.Middleware {
	return []api.Middleware{
		func(next http.Handler) http.Handler {
			return webtokens.NewMiddleware(next, webtokens.MiddlewareConfig{
				PathPrefixExclusions: []string{"/public", "public"},
			})
		},
	}
}
```

## !!! IMPORTANT !!!

Middleware listed in this function will be applied `IN ORDER`, from top to bottom. Internally, the framework will reverse the array so that the
middleware that comes first in the list will be run first. 

This is what the code to apply middleware looks like:
```go
// Add all applied middleware, in the order specified. Requires reversing the array
slices.Reverse(server.middleware)
for _, middleware := range server.middleware {
    handler = middleware(handler)
}
```

### services.go

Here you can specify which services you want to enable in your app. To enable a service, import it from `/api/apiservices` and include it in the returned array:

```go
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
```
Any skipped services will not be initialized and will throw errors when their methods are accessed.

In the event you want to remove services that rely on docker containers, you can launch the app with only certain containers active like so:

e.g. launch with *only* the postgres service active

```bash
docker compose up database
```
And omit unused services from the configuration
```go
func ConfigureServices() []api.Service {
	return []api.Service{
		apiservices.Postgres,
		apiservices.Email,
		apiservices.Sessions,
		apiservices.Webtokens,
	}
}
```

# Services

## Postgres
The internal postgres api is implemented using [sqlx](https://github.com/jmoiron/sqlx). You can access the database by doing the following:

```go
import "golang-server-base/api/postgresapi"

db, err := postgresapi.Database()
if err != nil {
	...
}
```
This will return an `*sqlx.DB` object. The usage of this object remains the same as the package it is imported from.


## Minio
The internal minio client uses [minio-go](https://github.com/minio/minio-go) and can be accessed like so:

```go
import "golang-server-base/api/minioapi"

client := minioapi.Client()
```
Since the minio client is not instantiated with an active connection to the store, no error is returned. Errors regarding this connection will
occur when the client is used. (TODO: There must be a clean way to test the connection...)


## Email
The internal email api is implemented using [gomail](https://github.com/go-gomail/gomail), and can be used as follows:

```go
import "golang-server-base/api/emailapi"
import emailmodels "golang-server-base/api/emailapi/models"

err := emailapi.Send(emailmodels.EmailOptions{
	...
})
if err != nil {
	...
}

```


## Meilisearch
The internal meilisearch client uses [meilisearch-go](https://github.com/meilisearch/meilisearch-go) and can be accessed like so:

```go
import "golang-server-base/api/meilisearchapi"

client, err := meilisearchapi.Client()
if err != nil {

}
```

`meilisearch.Client()` will return `nil, error` if the service is deemed unhealthy.