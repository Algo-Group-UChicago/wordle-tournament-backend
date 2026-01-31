package common

import (
	"context"
	"log/slog"
)

// Justification for log.go in common:
// - Though exclusively used in the handlers, doesn't make sense to scope it there.
// - Reduce module bloat with another `logging` module

// Difference between LogWarning and LogError is that LogWarnings result from malformed requests,
// while LogErrors result from server errors.

func LogInfo(endpoint string, msg string, attrs ...slog.Attr) {
	all := make([]slog.Attr, 0, 1+len(attrs))
	all = append(all, slog.String("endpoint", endpoint))
	all = append(all, attrs...)
	slog.Default().LogAttrs(context.Background(), slog.LevelInfo, msg, all...)
}

func LogWarning(endpoint string, msg string, status int, attrs ...slog.Attr) {
	all := make([]slog.Attr, 0, 2+len(attrs))
	all = append(all,
		slog.String("endpoint", endpoint),
		slog.Int("status", status),
	)
	all = append(all, attrs...)
	slog.Default().LogAttrs(context.Background(), slog.LevelWarn, msg, all...)
}

func LogError(endpoint string, msg string, status int, attrs ...slog.Attr) {
	all := make([]slog.Attr, 0, 2+len(attrs))
	all = append(all,
		slog.String("endpoint", endpoint),
		slog.Int("status", status),
	)
	all = append(all, attrs...)
	slog.Default().LogAttrs(context.Background(), slog.LevelError, msg, all...)
}
