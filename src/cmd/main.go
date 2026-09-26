package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/grandcat/zeroconf"
	"github.com/tg4-dev/air/src/internal/discovery"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Cannot get hostname: %s", err)
	}

	server, err := zeroconf.Register(hostname, "_airnode._udp", "local.", 12345, []string{"info=test"}, nil)
	if err != nil {
		log.Fatalf("Cannot create mdnsService: %s", err)
	}
	defer server.Shutdown()

	for {
		select {
		case <-sigs:
			fmt.Println("Shutting down...")
			cancel()
			return
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
