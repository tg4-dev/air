package discovery

import "time"

type Peer struct {
	Endpoint
	lastSeen time.Time
}
