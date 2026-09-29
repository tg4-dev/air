package utils

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger() *slog.Logger {
	level := getLogLevelFromEnv()

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}

func getLogLevelFromEnv() slog.Level {
	env := strings.ToLower(os.Getenv("AIR_LOG"))
	level := slog.LevelInfo

	switch env {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	return level
}
