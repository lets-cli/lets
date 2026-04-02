package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func LetsUserDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %w", err)
	}

	return filepath.Join(homeDir, ".config", "lets"), nil
}

func LetsUserFile(name string) (string, error) {
	dir, err := LetsUserDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, name), nil
}
