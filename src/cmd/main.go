package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tg4-dev/air/src/internal/discovery"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	browser := discovery.Browser{}

	browser.Update(ctx)
	peers := browser.GetPeers()

	fmt.Println(peers)
}
