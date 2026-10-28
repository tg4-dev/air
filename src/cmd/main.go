package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tg4-dev/air/src/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, err := app.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to init app:", err)
		os.Exit(1)
	}
	if err := a.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "app error:", err)
		os.Exit(1)
	}
}
