package utils

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	// TODO: log level from env
	level := slog.LevelDebug

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
