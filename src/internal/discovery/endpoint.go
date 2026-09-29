package discovery

import "net"

type Endpoint struct {
	ID    string
	Name  string
	Addrs []net.IP
	Port  int
	Meta  map[string]string
}
