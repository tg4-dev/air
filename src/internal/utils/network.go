package utils

import (
	"fmt"
	"net"
)

func GetActiveInterface() (*net.Interface, error) {
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

func GetIpsFromInterface(iface *net.Interface) ([]net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	var ips []net.IP
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok {
			var ip net.IP = ipNet.IP
			ips = append(ips, ip)
		}
	}

	return ips, nil
}
