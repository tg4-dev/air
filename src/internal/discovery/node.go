package discovery

import (
	"os"

	"github.com/tg4-dev/air/src/internal/utils"
)

type Node struct {
	self  Endpoint `json:"self"`
	peers []Peer
}

func NewNode() (*Node, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	// TODO: fill Meta
	iface, err := utils.GetActiveInterface()
	if err != nil {
		return nil, err
	}
	ips, err := utils.GetIpsFromInterface(iface)
	if err != nil {
		return nil, err
	}

	endpoint, err := NewEndpoint(hostname, ips, 12345, nil)
	if err != nil {
		return nil, err
	}

	return &Node{self: *endpoint, peers: nil}, nil
}
