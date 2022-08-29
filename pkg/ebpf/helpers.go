package ebpf

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
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

// TODO pull this network data in from aws or whatever cloud provider we are utilizing
func createNetworkMap() (map[string]*net.IPNet, error) {
	cidrs := []string{"127.0.0.1/16", "192.168.0.0/16", "172.17.0.0/16"}
	netMap := make(map[string]*net.IPNet)
	for idx, cidr := range cidrs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		netMap[strconv.Itoa(idx)] = ipNet
	}
	log.Println(netMap)
	return netMap, nil
}
