package discovery

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/mdns"
)

const metaServiceQuery = "_services._dns-sd._udp"
const scanTimeout = 3 * time.Second

func DiscoverServiceTypes(ctx context.Context) []string {
	entriesChannel := make(chan *mdns.ServiceEntry, 64)
	seen := map[string]bool{}
	var types []string
	var wg sync.WaitGroup

	wg.Add(1)
	wg.Go(func() {
		defer wg.Done()
		for entry := range entriesChannel {
			name := entry.Name
			if !seen[name] {
				seen[name] = true
				types = append(types, name)
			}
		}
	})

	params := mdns.DefaultParams(metaServiceQuery)
	params.Entries = entriesChannel
	params.DisableIPv6 = false
	params.Timeout = scanTimeout

	_ = mdns.Query(params)
	close(entriesChannel)

	wg.Wait()

	return types
}

func ScanAllServices(ctx context.Context, types []string) []Node {
	var nodes []Node
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, svcType := range types {
		wg.Add(1)
		go func(svcType string) {
			defer wg.Done()

			entriesChannel := make(chan *mdns.ServiceEntry, 32)
			done := make(chan struct{})

			go func() {
				defer close(done)
				for entry := range entriesChannel {
					mu.Lock()
					nodes = append(nodes, Node{
						Name:    entry.Name,
						Host:    entry.Host,
						AddrV4:  string(entry.AddrV4),
						AddrV6:  string(entry.AddrV6),
						Port:    entry.Port,
						Info:    entry.Info,
						Service: svcType,
					})
					mu.Unlock()
				}
			}()
			params := mdns.DefaultParams(svcType)
			params.Entries = entriesChannel
			params.Timeout = scanTimeout

			_ = mdns.Query(params)
			close(entriesChannel)
			<-done
		}(svcType)
	}

	wg.Wait()
	return nodes
}
