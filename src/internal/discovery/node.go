package discovery

type Node struct {
	Name    string
	Host    string
	AddrV4  string
	AddrV6  string
	Port    int
	Info    string
	Service string
}
