package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/tg4-dev/air/src/internal/discovery"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Cannot get hostname: %s", err)
	}

	mdnsService, err := mdns.NewMDNSService(hostname, "_airnode._udp", "", "", 12345, nil, []string{"info=test"})
	if err != nil {
		log.Fatalf("Cannot create mdnsService: %s", err)
	}

	isRunning := true

	mdnsServer, err := mdns.NewServer(&mdns.Config{Zone: mdnsService})
	if err != nil {
		log.Fatalf("Cannot create mdnsServer: %s", err)
	}
	defer mdnsServer.Shutdown()
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
