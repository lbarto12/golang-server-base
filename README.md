# What is this project?
This project is a framework written in go that has a variety prebuilt services, libraries, and driver code common 
to many apps being built.

- Current services and libraries include:
    - Postgres
    - Minio
    - Meilisearch       <---- NOT IMPLEMENTED
    - Email             <---- NOT IMPLEMENTED
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
```

# How to run

For first time installation, run the bash script to generate and initialize the services:
```bash
./rebuild-docker.sh
```

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
Within this directory are two files: `cors.go` and `main.go`

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
