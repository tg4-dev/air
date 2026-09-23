package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tg4-dev/air/src/internal/discovery"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	types := discovery.DiscoverServiceTypes(ctx)
	nodes := discovery.ScanAllServices(ctx, types)

	if len(nodes) == 0 {
		fmt.Println("Empty nodes list")
	}
	for _, node := range nodes {
		fmt.Printf("%+v\n", node)
	}
}
