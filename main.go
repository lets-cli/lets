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
)

var version = "0.0.0-dev"

func main() {
	ctx := getContext()

	logging.InitLogging(env.IsDebug())

	cfg, err := config.Read(version)
	if err != nil {
		logging.Log.Error(err.Error())
		os.Exit(1)
	}

	if err = workdir.CreateDotLetsDir(cfg.WorkDir); err != nil {
		logging.Log.Error(err.Error())
		os.Exit(1)
	}

	rootCmd := cmd.CreateRootCommand(os.Stdout, cfg)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logging.Log.Error(err.Error())

		exitCode := 1
		if e, ok := err.(*runner.RunErr); ok {
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
// so command's process can begin finishing earlier than cancel will say it to
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
