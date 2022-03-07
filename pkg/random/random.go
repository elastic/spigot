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

// Port returns a random integer from 0 to 65535.
func Port() int {
	return rand.Intn(65536)
}
