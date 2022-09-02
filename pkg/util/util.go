package util

import (
	"fmt"
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
