package main

import (
	"log"
	"wordle-tournament-backend/internal/config"
	"wordle-tournament-backend/internal/server"
)

func main() {
	log.Printf("Starting Wordle Tournament API...")
	log.Printf("Port: %s", config.Port)

	srv := server.New()

	log.Printf("Server listening on :%s", config.Port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
