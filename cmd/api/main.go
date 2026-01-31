package main

import (
	"log/slog"
	"os"
	"wordle-tournament-backend/internal/common"
	"wordle-tournament-backend/internal/config"
	"wordle-tournament-backend/internal/server"
)

func main() {
	cfg := config.Get()

	// JSON logs above INFO level are sent to stderr
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	srv := server.New()

	common.LogInfo("main", "Wordle Tournament API listening on :"+cfg.Port)
	if err := srv.Start(cfg.Port); err != nil {
		common.LogError("main", "Failed to start server", 500, slog.String("error", err.Error()))
		os.Exit(1)
	}
}
