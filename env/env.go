package env

import "os"

// GetConfigPathFromEnv return config file name and config dir
// LETS_CONFIG_DIR convenient to use in tests or when you want to run lets in another dir
func GetConfigPathFromEnv() (string, string) {
	return os.Getenv("LETS_CONFIG"), os.Getenv("LETS_CONFIG_DIR")
}
