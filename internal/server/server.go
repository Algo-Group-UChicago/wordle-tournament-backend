package server

import (
	"fmt"
	"net/http"
	"wordle-tournament-backend/internal/config"
	"wordle-tournament-backend/internal/handlers"
)

type Server struct {
	config *config.Config
	mux    *http.ServeMux
}

func New(cfg *config.Config) *Server {
	s := &Server{
		config: cfg,
		mux:    http.NewServeMux(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/health", handlers.HealthHandler(s.config.Environment))
	s.mux.HandleFunc("/api/guesses", handlers.GuessesHandler())

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"Wordle Tournament API","version":"1.0.0"}`)
	})
}

func (s *Server) Start() error {
	addr := ":" + s.config.Port
	return http.ListenAndServe(addr, s.mux)
}
