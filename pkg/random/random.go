// Package random provides functions for generating random objects using math/rand
package random

import (
	"math/rand"
	"net"
)

// IPv4 returns a random net.IP from the IPv4 address space.  No
// effort is made to prevent non-routable addresses.
func IPv4() net.IP {
	u32 := rand.Uint32()
	return net.IPv4(byte(u32&0xff), byte((u32>>8)&0xff), byte((u32>>16)&0xff), byte((u32>>24)&0xff))
}

// IPIn returns a random net.IP from the address spaces specified
// in the provided CIDR addresses ranges. The ranges may be IPv4 or IPv6.
func IPIn(ranges ...string) (net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(ranges[rand.Intn(len(ranges))])
	if err != nil {
		return nil, err
	}
	for i, m := range ipnet.Mask {
		ip[i] = byte(rand.Intn(255))&^m | ipnet.IP.Mask(ipnet.Mask)[i]
	}
	return ip[:len(ipnet.Mask)], nil
}

// Port returns a random integer from 0 to 65535.
func Port() int {
	return rand.Intn(65536)
}
