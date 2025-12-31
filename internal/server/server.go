package server

import (
	"fmt"
	"net/http"
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

// Handler returns the HTTP handler for the server (for testing)
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/health", handlers.HealthHandler())
	s.mux.HandleFunc("/start", handlers.StartHandler())
	s.mux.HandleFunc("/api/guesses", handlers.GuessesHandler())

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"Wordle Tournament API","version":"1.0.0"}`)
	})
}

func (s *Server) Start(port string) error {
	addr := ":" + port
	return http.ListenAndServe(addr, s.mux)
}
