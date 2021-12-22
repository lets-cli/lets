package config

import (
	"fmt"
	"os"
)

// workDir is where lets.yaml found or rootDir points to.
func getWorkDir(filename string, rootDir string) (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get workdir for config %s: %w", filename, err)
	}

	if rootDir != "" {
		workDir = rootDir
	}

	return workDir, nil
}
