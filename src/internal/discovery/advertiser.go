package discovery

import (
	"fmt"
	"net"
	"os"

	"github.com/grandcat/zeroconf"
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

	iface, err := getActiveInterface()
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

func getActiveInterface() (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		if iface.Flags&net.FlagMulticast == 0 {
			continue
		}

		addrs, _ := iface.Addrs()
		if len(addrs) == 0 {
			continue
		}
		return &iface, nil
	}
	return nil, fmt.Errorf("no suitable ifaces found")
}
