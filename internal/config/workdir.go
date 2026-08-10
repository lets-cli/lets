package config

import (
	"fmt"
	"os"
)

// getSearchDir is where lets starts looking for the config file: the process
// cwd, or rootDir when the user pinned one via --config-dir / LETS_CONFIG_DIR.
func getSearchDir(filename string, rootDir string) (string, error) {
	searchDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get workdir for config %s: %w", filename, err)
	}

	if rootDir != "" {
		searchDir = rootDir
	}

	return searchDir, nil
}
