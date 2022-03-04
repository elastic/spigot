package random

import (
	"math/rand"
	"net"
)

func IPv4() net.IP {
	u32 := rand.Uint32()
	return net.IPv4(byte(u32&0xff), byte((u32>>8)&0xff), byte((u32>>16)&0xff), byte((u32>>24)&0xff))
}
func Port() int {
	return rand.Intn(65536)
}
