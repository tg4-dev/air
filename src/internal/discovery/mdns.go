package discovery

import (
	"context"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
)

const metaServiceQuery = "_services._dns-sd._udp"
const airServiceQuery = "_airnode._udp"
const scanTimeout = 3 * time.Second

func firstAddr(addrs []net.IP) string {
	if len(addrs) == 0 {
		return ""
	}
	return addrs[0].String()
}

func DiscoverServiceTypes(ctx context.Context) []string {
	entriesChannel := make(chan *zeroconf.ServiceEntry, 64)
	seen := map[string]bool{}
	var types []string
	var wg sync.WaitGroup

	wg.Add(1)
	wg.Go(func() {
		defer wg.Done()
		for entry := range entriesChannel {
			name := entry.Instance
			if !seen[name] {
				seen[name] = true
				types = append(types, name)
			}
		}
	})

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		close(entriesChannel)
		wg.Wait()
		return types
	}

	browseCtx, cancel := context.WithTimeout(context.Background(), scanTimeout)
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			cancel()
		case <-browseCtx.Done():
		}
	}()

	_ = resolver.Browse(browseCtx, metaServiceQuery, "local.", entriesChannel)
	<-browseCtx.Done()

	wg.Wait()

	return types
}

func ScanAirNodes(ctx context.Context) ([]Node, error) {
	entriesChannel := make(chan *zeroconf.ServiceEntry, 64)
	var nodes []Node
	var wg sync.WaitGroup

	wg.Go(func() {
		for entry := range entriesChannel {
			nodes = append(nodes, Node{
				Name:   entry.Instance,
				Host:   entry.HostName,
				AddrV4: firstAddr(entry.AddrIPv4),
				AddrV6: firstAddr(entry.AddrIPv6),
				Port:   entry.Port,
				Info:   strings.Join(entry.Text, " "),
				// Service: svcType,
			})
		}
	})

	var opts []zeroconf.ClientOption
	if runtime.GOOS == "darwin" {
		iface, err := net.InterfaceByName("en0")
		if err != nil {
			close(entriesChannel)
			wg.Wait()
			return nil, err
		}
		opts = append(opts, zeroconf.SelectIfaces([]net.Interface{*iface}))
	}

	resolver, err := zeroconf.NewResolver(opts...)
	if err != nil {
		close(entriesChannel)
		wg.Wait()
		return nil, err
	}

	browseCtx, cancel := context.WithTimeout(context.Background(), scanTimeout)
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			cancel()
		case <-browseCtx.Done():
		}
	}()

	err = resolver.Browse(browseCtx, airServiceQuery, "local.", entriesChannel)
	if err != nil {
		cancel()
		wg.Wait()
		return nil, err
	}

	<-browseCtx.Done()
	wg.Wait()

	return nodes, nil
}

func ScanAllServices(ctx context.Context, types []string) []Node {
	var nodes []Node
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, svcType := range types {
		wg.Add(1)
		go func(svcType string) {
			defer wg.Done()

			entriesChannel := make(chan *zeroconf.ServiceEntry, 32)

			done := make(chan struct{})
			go func() {
				defer close(done)
				for entry := range entriesChannel {
					mu.Lock()
					nodes = append(nodes, Node{
						Name:    entry.Instance,
						Host:    entry.HostName,
						AddrV4:  firstAddr(entry.AddrIPv4),
						AddrV6:  firstAddr(entry.AddrIPv6),
						Port:    entry.Port,
						Info:    strings.Join(entry.Text, " "),
						Service: svcType,
					})
					mu.Unlock()
				}
			}()

			resolver, err := zeroconf.NewResolver(nil)
			if err != nil {
				close(entriesChannel)
				<-done
				return
			}

			browseCtx, cancel := context.WithTimeout(context.Background(), scanTimeout)
			defer cancel()
			go func() {
				select {
				case <-ctx.Done():
					cancel()
				case <-browseCtx.Done():
				}
			}()

			_ = resolver.Browse(browseCtx, svcType, "local.", entriesChannel)
			<-browseCtx.Done()
			<-done
		}(svcType)
	}

	wg.Wait()
	return nodes
}
