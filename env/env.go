package env

import (
	"os"
	"strconv"
)

// IsDebug checks LETS_DEBUG env. If set to true or 1 - we in debug mode.
func IsDebug() bool {
	debug, err := strconv.ParseBool(os.Getenv("LETS_DEBUG"))
	if err != nil {
		return false
	}

	return debug
}

func IsNotColorOutput() bool {
	notColored, err := strconv.ParseBool(os.Getenv("NO_COLOR"))
	if err != nil {
		return false
	}

	return notColored
}
