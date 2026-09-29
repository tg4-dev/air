package discovery

import (
	"os"

	"github.com/google/uuid"
	"github.com/tg4-dev/air/src/internal/utils"
)

type Node struct {
	self  Endpoint
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

	return &Node{self: Endpoint{
			ID:    uuid.New().String(),
			Name:  hostname,
			Addrs: ips,
			Port:  12345,
			Meta:  nil,
		},
			peers: nil},
		nil
}
