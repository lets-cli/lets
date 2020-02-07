package main

import (
	"fmt"
	"os"

	"github.com/kindritskyiMax/lets/cmd"
	"github.com/kindritskyiMax/lets/config"
)

func main() {
	conf, err := config.Load("lets.yaml", "")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	rootCmd := cmd.CreateRootCommand(conf, os.Stdout, GetVersion())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
