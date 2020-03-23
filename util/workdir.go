package util

import (
	"fmt"
	"os"
)

const dotLetsDir = ".lets"

// CreateDotLetsDir creates .lets dir where lets.yaml located.
// If directory already exists - skip creation.
func CreateDotLetsDir() error {
	if err := os.Mkdir(dotLetsDir, 0755); err != nil {
		if os.IsExist(err) {
			// its ok if we already have a dir, just return
			return nil
		}
		return fmt.Errorf("failed to create %s workdir: %s", dotLetsDir, err)
	}

	return nil
}
