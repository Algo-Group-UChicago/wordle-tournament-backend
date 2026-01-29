package main

import (
	"log"
	"log/slog"
	"os"
	"wordle-tournament-backend/internal/config"
	"wordle-tournament-backend/internal/server"
)

func main() {
	cfg := config.Get()

	// JSON logs above INFO level are sent to stderr
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	log.Printf("Starting Wordle Tournament API...")
	log.Printf("Port: %s", cfg.Port)

	srv := server.New()

	log.Printf("Server listening on :%s", cfg.Port)
	if err := srv.Start(cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
