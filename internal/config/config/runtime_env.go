package config

import (
	"path/filepath"
	"runtime"
	"strings"
)

func (c *Config) BuiltinEnv(shell string) map[string]string {
	return map[string]string{
		"LETS_CONFIG":     filepath.Base(c.FilePath),
		"LETS_CONFIG_DIR": filepath.Dir(c.FilePath),
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
