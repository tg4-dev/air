package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/tg4-dev/air/src/internal/discovery"
	"github.com/tg4-dev/air/src/internal/nodeinfo"
	"github.com/tg4-dev/air/src/internal/utils"
)

func main() {
	logger := utils.NewLogger()
	startedAt := time.Now()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	node, err := discovery.NewNode()
	if err != nil {
		logger.Error("failed to create node", "err", err)
		os.Exit(1)
	}

	provider := nodeinfo.NewDefaultProvider(node.ID(), node.Hostname(), startedAt)
	infoServer := nodeinfo.NewServer(provider)
	addr := ":" + strconv.Itoa(node.Port())
	if err := infoServer.Start(addr); err != nil {
		logger.Error("failed to start nodeinfo HTTP server", "addr", addr, "err", err)
		os.Exit(1)
	}
	logger.Info("nodeinfo HTTP server listening", "addr", addr)

	advertiser, err := discovery.NewAdvertiser(*node)
	if err != nil {
		logger.Error("failed to start mDNS advertiser", "err", err)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = infoServer.Shutdown(shutdownCtx)
		shutdownCancel()
		os.Exit(1)
	}
	defer advertiser.Shutdown()

	browser := discovery.Browser{}
	browser.Update(ctx)
	peers := browser.GetPeers()

	fmt.Println("===== PEERS =====")
	for _, peer := range peers {
		fmt.Printf("%+v\n", peer)
	}

	<-sigs
	logger.Info("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := infoServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("nodeinfo HTTP shutdown", "err", err)
	}
}
