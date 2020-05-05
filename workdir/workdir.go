package workdir

import (
	"path/filepath"

	"github.com/lets-cli/lets/util"
)

const dotLetsDir = ".lets"

// CreateDotLetsDir creates .lets dir where lets.yaml located.
// If directory already exists - skip creation.
func CreateDotLetsDir(workDir string) error {
	fullPath, err := GetDotLetsDir(workDir)
	if err != nil {
		return err
	}

	return util.SafeCreateDir(fullPath)
}

func GetDotLetsDir(workDir string) (string, error) {
	return filepath.Abs(filepath.Join(workDir, dotLetsDir))
}
