package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tg4-dev/air/src/internal/discovery"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	browser := discovery.Browser{}
	node, err := discovery.NewNode()
	fmt.Printf("%+v\n", node)
	if err != nil {
		panic(err)
	}
	advertiser, err := discovery.NewAdvertiser(*node)
	if err != nil {
		panic(err)
	}
	defer advertiser.Shutdown()

	browser.Update(ctx)
	peers := browser.GetPeers()

	fmt.Println("===== PEERS =====")
	for _, peer := range peers {
		fmt.Printf("%+v\n", peer)
	}

	<-sigs
	fmt.Println("Shutting down")
}
