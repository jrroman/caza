package ebpf

import (
	"encoding/binary"
	"net"
	"unsafe"
)

var (
	NativeEndian binary.ByteOrder
)

func init() {
	// Determine the endianness of the host machine to translate network addrs
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		NativeEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		NativeEndian = binary.BigEndian
	default:
		panic("could not determine endianness of host.")
	}
}

// intToIP converts IPv4 number to net.IP
func intToIP(ipNum uint32) net.IP {
	ip := make(net.IP, 4)
	NativeEndian.PutUint32(ip, ipNum)
	return ip
}
