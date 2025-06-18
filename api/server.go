package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

type Server struct {
	Host     string
	Port     string
	Cors     *cors.Options
	handlers map[string]http.Handler
}

func (server *Server) AddHandler(path string, handler http.Handler) error {
	if server.handlers == nil {
		server.handlers = map[string]http.Handler{}
	}
	if _, exists := server.handlers[path]; exists {
		return fmt.Errorf("server handler `%s` already exists, and may not be overwritten", path)
	}
	server.handlers[path] = handler
	return nil
}

func (server *Server) AddHandlers(handlers map[string]http.Handler) error {
	for path, handler := range handlers {
		err := server.AddHandler(path, handler)
		if err != nil {
			return err
		}
	}
	return nil
}

func (server *Server) Launch() error {
	// Create MUX
	mux := http.NewServeMux()

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
		mux.Handle(path, handler)
	}

	// Add CORS, if applied
	handler := http.Handler(mux)
	if server.Cors != nil {
		handler = cors.New(*server.Cors).Handler(mux)
	}

	// Launch Server
	log.Printf("Started server on %s, listening for requests...\n", apiUrl)
	return http.ListenAndServe(apiUrl, handler)
}
