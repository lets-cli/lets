package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lets-cli/lets/cmd"
	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/env"
	"github.com/lets-cli/lets/logging"
	"github.com/lets-cli/lets/runner"
	"github.com/lets-cli/lets/workdir"
	"github.com/spf13/cobra"
)

var version = "0.0.0-dev"

func main() {
	ctx := getContext()

	logging.InitLogging(env.IsDebug(), os.Stdout, os.Stderr)

	cfg, readConfigErr := config.Load(version)

	var rootCmd *cobra.Command
	if cfg != nil {
		rootCmd = cmd.CreateRootCommandWithConfig(os.Stdout, cfg, version)

		if err := workdir.CreateDotLetsDir(cfg.WorkDir); err != nil {
			logging.Log.Error(err)
			os.Exit(1)
		}
	} else {
		rootCmd = cmd.CreateRootCommand(version)
	}

	if readConfigErr != nil {
		cmd.ConfigErrorCheck(rootCmd, readConfigErr)
	}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logging.Log.Error(err.Error())

		exitCode := 1
		if e, ok := err.(*runner.RunErr); ok { //nolint:errorlint
			exitCode = e.ExitCode()
		}

		os.Exit(exitCode)
	}
}

// getContext returns context and kicks of a goroutine
// which waits for SIGINT, SIGTERM and cancels global context.
//
// Note that since we setting stdin to command we run, that command
// will receive SIGINT, SIGTERM at the same time as we here,
// so command's process can begin finishing earlier than cancel will say it to.
func getContext() context.Context {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-ch
		logging.Log.Printf("lets: signal received: %s", sig)
		cancel()
	}()

	return ctx
}
