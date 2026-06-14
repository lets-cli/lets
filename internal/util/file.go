package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// WriteFileAtomic writes data to dst via a sibling temp file and an os.Rename,
// so dst is never left in a partially-written state if the process is interrupted.
// The resulting file has 0o644 permissions.
func WriteFileAtomic(dst string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(dst), "*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	tmpPath := tmp.Name()

	_, writeErr := tmp.Write(data)
	closeErr := tmp.Close()

	if writeErr != nil || closeErr != nil {
		os.Remove(tmpPath)

		if writeErr != nil {
			return fmt.Errorf("failed to write temp file: %w", writeErr)
		}

		return fmt.Errorf("failed to close temp file: %w", closeErr)
	}

	//#nosec G306
	if err := os.Chmod(tmpPath, 0o644); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to set permissions on temp file: %w", err)
	}

	if err := os.Rename(tmpPath, dst); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp file to %s: %w", dst, err)
	}

	return nil
}
