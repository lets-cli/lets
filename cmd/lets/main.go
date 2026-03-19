package main

import (
	"os"

	"github.com/lets-cli/lets/internal/cli"
)

var Version = "0.0.0-dev"
var BuildDate = ""

func main() {
	os.Exit(cli.Main(Version, BuildDate))
}
