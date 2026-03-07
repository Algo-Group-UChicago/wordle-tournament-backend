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

	// Uncomment this if we switch to json logs
	// JSON logs above INFO level are sent to stderr
	// slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	// Text logs above INFO level are sent to stderr
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	srv := server.New()

	common.LogInfo("ServerStart", &common.LogData{Msg: "Listening on :" + cfg.Port})
	if err := srv.Start(cfg.Port); err != nil {
		common.LogError("ServerStartFailure", &common.LogData{Msg: err.Error()})
		os.Exit(1)
	}
}
