package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New(environment string) *slog.Logger {
	level := slog.LevelInfo

	switch strings.ToLower(environment) {
	case "development":
		level = slog.LevelDebug
	case "production":
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
