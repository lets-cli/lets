package util

import (
	"fmt"
	"os"
)

func SafeCreateDir(dirPath string) error {
	if err := os.Mkdir(dirPath, 0755); err != nil {
		if os.IsExist(err) {
			// its ok if we already have a dir, just return
			return nil
		}

		return fmt.Errorf("failed to create %s dir: %w", dirPath, err)
	}

	return nil
}
