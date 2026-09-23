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

	for _, t := range types {
		fmt.Println(t)
	}
}
