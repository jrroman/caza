package cloud

import (
	"net"

	"github.com/jrroman/caza/pkg/config"
)

type Cloud interface {
	GetNetworks(*config.Config) (map[string]*net.IPNet, error)
}
