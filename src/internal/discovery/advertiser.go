package discovery

import (
	"context"

	"github.com/grandcat/zeroconf"
)

type Advertiser struct {
	node   Node
	server zeroconf.Server
}

func NewAdvertiser(node Node) *Advertiser {

	server, err := zeroconf.Register()
	return &Advertiser{node: node, server: )}
}
func (a *Advertiser) AdvertiseOnce(ctx context.Context) {

}
