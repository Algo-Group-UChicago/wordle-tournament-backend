package main

import (
	"log"
	"wordle-tournament-backend/internal/config"
	"wordle-tournament-backend/internal/registry"
	"wordle-tournament-backend/internal/server"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting Wordle Tournament API...")
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("Port: %s", cfg.Port)

	registry.Initialize("corpus.txt", "possible_answers.txt")

	srv := server.New(cfg)

	log.Printf("Server listening on :%s", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
