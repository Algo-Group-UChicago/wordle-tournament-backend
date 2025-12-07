package server

import (
	"fmt"
	"net/http"
	"wordle-tournament-backend/internal/config"
	"wordle-tournament-backend/internal/handlers"
)

type Server struct {
	mux *http.ServeMux
}

func New() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/health", handlers.HealthHandler())
	s.mux.HandleFunc("/api/guesses", handlers.GuessesHandler())

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"Wordle Tournament API","version":"1.0.0"}`)
	})
}

func (s *Server) Start() error {
	addr := ":" + config.Port
	return http.ListenAndServe(addr, s.mux)
}
