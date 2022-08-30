package util

import (
	"errors"
	"fmt"
	"os"
)

func EnsureEnvironmentSet(envKey string) error {
	env, ok := os.LookupEnv(envKey)
	if !ok {
		message := fmt.Sprintf("Environment %s is not set", envKey)
		return errors.New(message)
	}
	return nil
}
