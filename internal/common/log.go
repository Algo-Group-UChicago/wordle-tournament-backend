package common

import (
	"context"
	"log/slog"
)

// Justification for log.go in common:
// - Though exclusively used in the handlers, doesn't make sense to scope it there.
// - Reduce module bloat with another `logging` module

type LogData struct {
	TeamID string
	RunID  string
	Msg    string
}

func LogInfo(name string, data *LogData) {
	logWithLevel(slog.LevelInfo, name, data)
}

func LogWarning(name string, data *LogData) {
	logWithLevel(slog.LevelWarn, name, data)
}

func LogError(name string, data *LogData) {
	logWithLevel(slog.LevelError, name, data)
}

func logWithLevel(level slog.Level, name string, data *LogData) {
	attrs := []slog.Attr{slog.String("name", name)}
	if data != nil {
		if data.TeamID != "" {
			attrs = append(attrs, slog.String("team_id", data.TeamID))
		}
		if data.RunID != "" {
			attrs = append(attrs, slog.String("run_id", data.RunID))
		}
		if data.Msg != "" {
			attrs = append(attrs, slog.String("msg", data.Msg))
		}
	}

	slog.Default().LogAttrs(context.Background(), level, name, attrs...)
}
