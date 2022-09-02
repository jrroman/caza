package testutils

import "net"

func CreateIPNetHelper(cidrBlock string) *net.IPNet {
	_, ipNet, _ := net.ParseCIDR(cidrBlock)
	return ipNet
}
