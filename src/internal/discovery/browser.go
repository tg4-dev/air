package discovery

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
)

const airServiceQuery = "_airnode._udp"
const scanTimeout = 3 * time.Second

type Browser struct {
	peers []Peer
}

func (b *Browser) Update(ctx context.Context) bool {
	peers, err := b.Browse(ctx)
	if err != nil {
		return false
	}
	b.peers = peers
	return true
}
func (b *Browser) Browse(ctx context.Context) ([]Peer, error) {
	entriesChannel := make(chan *zeroconf.ServiceEntry, 64)
	var peers []Peer
	var wg sync.WaitGroup

	wg.Go(func() {
		for entry := range entriesChannel {
			peers = append(peers, Peer{
				Endpoint: Endpoint{
					Name:  entry.Instance,
					Addrs: getIpsFromEntry(entry),
					Port:  entry.Port,
					Meta:  parseMeta(entry.Text),
				},
				lastSeen: time.Now(),
			})
		}
	})

	var opts []zeroconf.ClientOption

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
	return peers, nil
}

func (b *Browser) GetPeers() []Peer {
	return b.peers
}

func getIpsFromEntry(entry *zeroconf.ServiceEntry) []net.IP {
	var addrs []net.IP

	for _, addr := range entry.AddrIPv4 {
		addrs = append(addrs, addr)
	}
	for _, addr := range entry.AddrIPv6 {
		addrs = append(addrs, addr)
	}

	return addrs
}

func parseMeta(txt []string) map[string]string {
	meta := make(map[string]string)
	for _, t := range txt {
		entries := strings.Fields(t)
		for _, entry := range entries {
			splitted := strings.Split(entry, "=")
			k, v := splitted[0], splitted[1]
			meta[k] = v
		}
	}
	return meta
}
