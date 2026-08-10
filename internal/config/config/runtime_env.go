package config

import (
	"path/filepath"
	"runtime"
	"strings"
)

func (c *Config) BuiltinEnv(shell string) map[string]string {
	letsConfig := filepath.Base(c.FilePath)
	// ConfigDir, not RootDir: for a remote config this is the cache dir holding the
	// downloaded yaml. The project root is simply the cwd now, so $PWD covers it.
	letsConfigDir := c.ConfigDir

	if c.RemoteSource != "" {
		letsConfig = c.RemoteSource
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
