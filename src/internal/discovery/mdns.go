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
