package discovery

import (
	"fmt"
	"net"
	"os"

	"github.com/grandcat/zeroconf"
	"github.com/tg4-dev/air/src/internal/utils"
)

type Advertiser struct {
	node   Node
	server *zeroconf.Server
}

func NewAdvertiser(node Node) (*Advertiser, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	iface, err := utils.GetActiveInterface()
	if err != nil {
		return nil, err
	}
	if iface == nil {
		return nil, fmt.Errorf("no active iface")
	}
	ifaces := []net.Interface{*iface}
	server, err := zeroconf.Register(hostname, "_airnode._udp", "local.", node.self.Port, nil, ifaces)
	if err != nil {
		return nil, err
	}
	return &Advertiser{node: node, server: server}, nil
}

func (a *Advertiser) Shutdown() {
	a.server.Shutdown()
}
