package cloud

import (
	"net"

	"github.com/jrroman/caza/pkg/config"
)

// TODO leave this here for now, will be useful once we move past a single cloud provider
type Cloud interface {
	GetNetworks(*config.Config) (map[string]*net.IPNet, error)
}
