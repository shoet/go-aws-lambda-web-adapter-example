package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ServerConfig struct {
	Port int
}

type Server struct {
	server   *http.Server
	listener net.Listener
}

func (s *Server) Start(l net.Listener) error {
	return s.server.Serve(l)
}

func NewServer(config *ServerConfig, handler http.Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	server := &http.Server{
		Handler: handler,
	}

	return &Server{server, listener}, nil
}

func NewHandlers() (*chi.Mux, error) {
	hch := NewHealthCheckHandler()

	mux := chi.NewMux()
	mux.Get("/health", hch.ServeHTTP)

	return mux, nil
}

type HealthCheckHandler struct{}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}
func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: "OK",
		Status:  http.StatusOK,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	config := &ServerConfig{Port: 8080}
	handler, err := NewHandlers()
	if err != nil {
		panic(err)
	}
	server, err := NewServer(config, handler)
	if err != nil {
		panic(err)
	}
	if err := server.Start(server.listener); err != nil {
		panic(err)
	}
}
