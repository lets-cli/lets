package workdir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lets-cli/lets/util"
	log "github.com/sirupsen/logrus"
)

const dotLetsDir = ".lets"

const defaultLetsYaml = `version: %s
shell: bash

commands:
  hello:
	description: Say hello
	options: |
		Usage: lets hello [<name>]
		Examples:
			lets hello
			lets hello Friend
	cmd: echo Hello, "${LETSOPT_NAME:-world}"!
`

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

// InitLetsFile creates lets.yaml int the current dir.
func InitLetsFile(workDir string, version string) error {
	f := filepath.Join(workDir, "lets.yaml")

	if _, err := os.Stat(f); err == nil {
		return fmt.Errorf("lets.yaml already exists in %s", workDir)
	}

	output := fmt.Sprintf(defaultLetsYaml, version)
	if err := os.WriteFile(f, []byte(output), 0644); err != nil {
		return err
	}

	log.Println("lets.yaml created in the current directory")

	return nil
}
