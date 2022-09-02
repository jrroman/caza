package config

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/jrroman/caza/internal"
)

// Config defines configuration for the application runtime
type Config struct {
	// Specifies whether or not we should be interacting with a cloud provider.
	// If this value is false the program will run only with the networks specified
	// within our Networks slice above
	CloudEnabled bool
	// Networks is an optional list of net.IPNet. If you are also pulling subnet
	// network data from a cloud provider the nework CIDRs specified here will
	// added to the list pulled from the cloud
	Networks map[string]*net.IPNet
	// Region signifies the cloud provider region you would like to interact with.
	// If VpcID is specified this value is required
	Region string
	// VpcID is the AWS private network ID to pull network data from
	VpcID string
}

// We need to ensure the CIDRs passed into opts.Networks are valid CIDR blocks
func validateNetworks(networkString string) (map[string]*net.IPNet, error) {
	// After splitting the slice of networks by comma we end up with a slice with
	// a format that looks like the following e.g. [local:127.0.0.1/32, router:192.168.0.0/16]
	networks := strings.Split(networkString, ",")
	networkMap := make(map[string]*net.IPNet)
	for _, network := range networks {
		networkParts := strings.Split(network, ":")
		if len(networkParts) != 2 {
			return nil, fmt.Errorf("invalid formatting of network string: %s", network)
		}
		name, cidrString := networkParts[0], networkParts[1]
		_, ipNet, err := net.ParseCIDR(cidrString)
		if err != nil {
			return nil, fmt.Errorf("invalid cidr block: %s", cidrString)
		}
		networkMap[name] = ipNet
	}
	return networkMap, nil
}

// Validate all of the input options to ensure we have a valid runtime configuration
func validateConfig(opts internal.Options) (*Config, error) {
	if opts.CloudEnabled && (opts.VpcID == "" || opts.Region == "") {
		return nil, errors.New("cloud enabled, region and vpc id values must be set")
	}
	cfg := &Config{
		CloudEnabled: opts.CloudEnabled,
		Region:       opts.Region,
		VpcID:        opts.VpcID,
	}
	var err error
	if opts.Networks != "" {
		cfg.Networks, err = validateNetworks(opts.Networks)
		if err != nil {
			return nil, err
		}
	}
	if len(cfg.Networks) == 0 && !cfg.CloudEnabled {
		return nil, errors.New("no networks specified; enable cloud or pass in networks")
	}
	return cfg, nil
}

func New(opts internal.Options) (*Config, error) {
	return validateConfig(opts)
}
