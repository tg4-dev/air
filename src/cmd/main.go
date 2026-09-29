package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tg4-dev/air/src/internal/discovery"
	"github.com/tg4-dev/air/src/internal/utils"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger := utils.NewLogger()

	logger.Info("starting air")
	browser := discovery.Browser{}
	logger.Info("Browser successfully created")
	node, err := discovery.NewNode()
	fmt.Printf("%+v\n", node)
	if err != nil {
		panic(err)
	}
	advertiser, err := discovery.NewAdvertiser(*node)
	if err != nil {
		panic(err)
	}
	logger.Info("Advertiser successfully created")
	defer advertiser.Shutdown()

	browser.Update(ctx)
	peers := browser.GetPeers()

	fmt.Println("===== PEERS =====")
	for i, peer := range peers {
		logger.Debug("peer info", "index", i+1, "peer", peer)
	}

	<-sigs
	logger.Info("shutting down...")
}
