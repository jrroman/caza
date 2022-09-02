package cloud

import (
	"net"
)

// TODO leave this here for now, will be useful once we move past a single cloud provider
type Cloud interface {
	GetNetworks(string) (map[string]*net.IPNet, error)
}
