package utils

import (
	"os"
	"strconv"
)

const (
	envDockerized = "DOCKERIZED"
)

func IsDockerized() bool {
	val := os.Getenv(envDockerized)

	ok, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return ok
}
