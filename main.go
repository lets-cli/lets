package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lets-cli/lets/cmd"
	"github.com/lets-cli/lets/config"
	"github.com/lets-cli/lets/env"
	"github.com/lets-cli/lets/executor"
	"github.com/lets-cli/lets/logging"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version = "0.0.0-dev"

func main() {
	ctx := getContext()

	configFile := os.Getenv("LETS_CONFIG")
	configDir := os.Getenv("LETS_CONFIG_DIR")

	logging.InitLogging(os.Stdout, os.Stderr)

	rootCmd := cmd.CreateRootCommand(version)
	rootCmd.InitDefaultHelpFlag()
	rootCmd.InitDefaultVersionFlag()
	reinitCompletionCmd := cmd.InitCompletionCmd(rootCmd, nil)
	rootCmd.InitDefaultHelpCmd()

	command, args, err := rootCmd.Traverse(os.Args[1:])
	if err != nil {
		log.Errorf("lets: traverse flags error: %s", err)
		os.Exit(1)
	}

	if err = rootCmd.ParseFlags(args); err != nil {
		log.Errorf("lets: parse flags error: %s", err)
		os.Exit(1)
	}

	_, err = env.ParseDebugLevel(rootCmd)
	if err != nil {
		log.Errorf("lets: parse debug level error: %s", err)
		os.Exit(1)
	}

	if env.IsDebug() {
		log.SetLevel(log.DebugLevel)
	}

	cfg, err := config.Load(configFile, configDir, version)
	if err != nil {
		if failOnConfigError(rootCmd, command) {
			log.Errorf("lets: config error: %s", err)
			os.Exit(1)
		}
	}

	if cfg != nil {
		reinitCompletionCmd(cfg)
		cmd.InitSubCommands(rootCmd, cfg, os.Stdout)
	}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Error(err.Error())

		exitCode := 1
		if e, ok := err.(*executor.ExecuteError); ok { //nolint:errorlint
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
		log.Printf("lets: signal received: %s", sig)
		cancel()
	}()

	return ctx
}

func failOnConfigError(root *cobra.Command, current *cobra.Command) bool {
	return root.Flags().NFlag() == 0 &&
		current.Name() != "completion" &&
		current.Name() != "help"
}
