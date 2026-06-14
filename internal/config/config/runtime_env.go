package config

import (
	"path/filepath"
	"runtime"
	"strings"
)

func (c *Config) BuiltinEnv(shell string) map[string]string {
	letsConfig := filepath.Base(c.FilePath)
	letsConfigDir := filepath.Dir(c.FilePath)

	if c.RemoteSource != "" {
		letsConfig = c.RemoteSource
		letsConfigDir = c.WorkDir
	}

	return map[string]string{
		"LETS_CONFIG":     letsConfig,
		"LETS_CONFIG_DIR": letsConfigDir,
		"LETS_OS":         runtime.GOOS,
		"LETS_ARCH":       runtime.GOARCH,
		"LETS_SHELL":      shell,
	}
}

func (c *Config) CommandBuiltinEnv(command *Command, shell string, workDir string) map[string]string {
	envMap := c.BuiltinEnv(shell)
	envMap["LETS_COMMAND_NAME"] = command.Name
	envMap["LETS_COMMAND_ARGS"] = strings.Join(command.Args, " ")
	envMap["LETS_COMMAND_WORK_DIR"] = workDir

	return envMap
}
