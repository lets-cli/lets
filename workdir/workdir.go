package workdir

import "github.com/lets-cli/lets/util"

const DotLetsDir = ".lets"

// CreateDotLetsDir creates .lets dir where lets.yaml located.
// If directory already exists - skip creation.
func CreateDotLetsDir() error {
	return util.SafeCreateDir(DotLetsDir)
}
