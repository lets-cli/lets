package main

import (
	"os"

	"github.com/lets-cli/lets/cmd"
	"github.com/lets-cli/lets/env"
	"github.com/lets-cli/lets/logging"
)

func main() {
	rootCmd := cmd.CreateRootCommand(os.Stdout, GetVersion())

	logging.InitLogging(env.IsDebug())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
