package main

import (
	"os"

	"github.com/lets-cli/lets/cmd"
)

func main() {
	rootCmd := cmd.CreateRootCommand(os.Stdout, GetVersion())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
