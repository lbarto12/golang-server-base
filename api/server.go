package api

import (
	"errors"
	"fmt"
	"golang-server-base/api/apiservices"
	"log"
	"net/http"
	"slices"

	"github.com/rs/cors"
)

type Middleware func(next http.Handler) http.Handler
type Service int

type Server struct {
	Host       string
	Port       string
	Cors       *cors.Options
	handlers   map[string]http.Handler
	Mux        *http.ServeMux
	middleware []Middleware
	services   []Service
	Services   apiservices.ServicesAccess
}

func (server *Server) AddServices(services []Service) {
	server.services = append(server.services, services...)
}

func (server *Server) AddMiddleWare(middleware Middleware) {
	server.middleware = append(server.middleware, middleware)
}

func (server *Server) AddMiddleWares(middleware []Middleware) {
	server.middleware = append(server.middleware, middleware...)
}

func (server *Server) AddHandler(path string, handler http.Handler) {
	if server.handlers == nil {
		server.handlers = map[string]http.Handler{}
	}
	if _, exists := server.handlers[path]; exists {
		panic(fmt.Errorf("server handler `%s` already exists, and may not be changed", path))
	}
	server.handlers[path] = handler
}

func (server *Server) AddHandlers(handlers map[string]http.Handler) {
	for path, handler := range handlers {
		server.AddHandler(path, handler)

	}
}

func (server *Server) AddCors(options *cors.Options) {
	server.Cors = options
}

func (server *Server) Routes() []string {
	var result []string
	for path := range server.handlers {
		result = append(result, path)
	}
	return result
}

func NewServer(initial Server) Server {
	initial.Mux = &http.ServeMux{}
	return initial
}

func (server *Server) Launch() error {

	// Assert required options are not null
	if server.Host == "" {
		return errors.New("LaunchServer: options.Path undefined")
	}
	if server.Port == "" {
		return errors.New("LaunchServer: options.Port undefined")
	}

	// Create final API URL
	apiUrl := fmt.Sprintf("%s:%s", server.Host, server.Port)

	// Add Handlers
	for path, handler := range server.handlers {
		server.Mux.Handle(path, handler)
	}

	// Add CORS, if applied
	handler := http.Handler(server.Mux)
	if server.Cors != nil {
		handler = cors.New(*server.Cors).Handler(server.Mux)
	}

	// Add all applied middleware, in the order specified. Requires reversing the array
	slices.Reverse(server.middleware)
	for _, middleware := range server.middleware {
		handler = middleware(handler)
	}

	// Launch Server
	log.Printf("Started server on %s, listening for requests...\n", apiUrl)
	return http.ListenAndServe(apiUrl, handler)
}
