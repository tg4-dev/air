package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tg4-dev/air/src/internal/discovery"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	isRunning := true

	for isRunning {
		select {
		case <-sigs:
			fmt.Println("Shutting down...")
			isRunning = false
		default:
			nodes, err := discovery.ScanAirNodes(ctx)
			if err != nil {
				log.Fatal("Cannot scan nodes ", err)
			}

			if len(nodes) == 0 {
				fmt.Println("Empty nodes list")
			}
			for _, node := range nodes {
				fmt.Printf("%+v\n", node)
			}
		}
	}
}
