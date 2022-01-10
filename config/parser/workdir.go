package parser

import (
	"fmt"

	"github.com/lets-cli/lets/config/config"
)

func parseWorkDir(rawWorkdir interface{}, newCmd *config.Command) error {
	workdir, ok := rawWorkdir.(string)
	if !ok {
		return fmt.Errorf("work_dir must be a string")
	}

	newCmd.WorkDir = workdir

	return nil
}
