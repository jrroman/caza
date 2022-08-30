package awscloud

import (
	"net"

	"github.com/jrroman/caza/pkg/config"
)

type AwsCloud struct{}

func GetNetworks(cfg *config.Config) (map[string]*net.IPNet, error) {
	_, ipNet, _ := net.ParseCIDR("10.0.0.0/16")
	mockNetworks := map[string]*net.IPNet{
		"a": ipNet,
	}
	return mockNetworks, nil
}
