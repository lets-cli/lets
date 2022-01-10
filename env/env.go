package env

import (
	"os"
	"strconv"
)

// GetConfigPathFromEnv return config file name and config dir
// LETS_CONFIG_DIR convenient to use in tests or when you want to run lets in another dir.
func GetConfigPathFromEnv() (string, string) {
	return os.Getenv("LETS_CONFIG"), os.Getenv("LETS_CONFIG_DIR")
}

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
