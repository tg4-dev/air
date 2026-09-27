package discovery

type Node struct {
	self  Endpoint
	peers []Peer
}

func NewNode() *Node
