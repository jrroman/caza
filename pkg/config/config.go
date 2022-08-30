package config

import (
	"errors"
	"net"
	"strings"

	"github.com/jrroman/caza/internal"
)

// Config defines configuration for the application runtime
type Config struct {
	// AdditionalNetworks is an optional list of net.IPNet. If you are also pulling subnet
	// network data from a cloud provider the nework CIDRs specified here will
	// added to the list pulled from the cloud
	AdditionalNetworks []*net.IPNet
	// Specifies whether or not we should be interacting with a cloud provider.
	// If this value is false the program will run only with the networks specified
	// within our AdditionalNetworks slice above
	CloudEnabled bool
	// Region signifies the cloud provider region you would like to interact with.
	// If VpcID is specified this value is required
	Region string
	// VpcID is the AWS private network ID to pull network data from
	VpcID string
}

// We need to ensure the CIDRs passed into opts.Networks are valid CIDR blocks
func validateNetworksOpt(netString string) ([]*net.IPNet, error) {
	networks := strings.Split(netString, ",")
	var cidrList []*net.IPNet
	for _, network := range networks {
		_, ipNet, err := net.ParseCIDR(network)
		if err != nil {
			return nil, err
		}
		cidrList = append(cidrList, ipNet)
	}
	return cidrList, nil
}

// Validate all of the input options to ensure we have a valid runtime configuration
func validateConfig(opts internal.Options) (*Config, error) {
	if opts.CloudEnabled && (opts.VpcID == "" || opts.Region == "") {
		return nil, errors.New("Cloud enabled, region and vpc id values must be set")
	}
	cfg := &Config{
		CloudEnabled: opts.CloudEnabled,
		Region:       opts.Region,
		VpcID:        opts.VpcID,
	}
	var err error
	if opts.Networks != "" {
		cfg.AdditionalNetworks, err = validateNetworksOpt(opts.Networks)
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func New(opts internal.Options) (*Config, error) {
	return validateConfig(opts)
}
