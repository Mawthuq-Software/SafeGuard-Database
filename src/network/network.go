package network

import "net"

func ParseIP(address string) (ipvX *net.IPNet, err error) { //parses string into IP address
	_, ipvX, err = net.ParseCIDR(address)
	return
}
