package discovery

import (
	"os"

	"github.com/google/uuid"
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
	// TODO: fill Addrs and Meta
	return &Node{self: Endpoint{
		ID:    uuid.New().String(),
		Name:  hostname,
		Addrs: nil,
		Port:  12345,
		Meta:  nil,
	}, peers: nil}, nil
}
