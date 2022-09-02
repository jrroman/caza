package util

import (
	"fmt"
	"net"
	"os"
)

// EnsureEnvironmentSet is a quick check to tell us whether or not the envKey is
// set as an environment variable if the value is unset we return an error to the user
func EnsureEnvironmentSet(envKey string) error {
	_, ok := os.LookupEnv(envKey)
	if !ok {
		return fmt.Errorf("environment %s is not set", envKey)
	}
	return nil
}

func MergeNetworkMaps(networks []map[string]*net.IPNet) map[string]*net.IPNet {
	// if there is only one network return it
	if len(networks) == 1 {
		return networks[0]
	}
	merged := make(map[string]*net.IPNet)
	for _, nm := range networks {
		for name, network := range nm {
			merged[name] = network
		}
	}
	return merged
}
