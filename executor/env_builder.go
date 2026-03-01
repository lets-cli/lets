package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lets-cli/lets/checksum"
	"github.com/lets-cli/lets/config/config"
)

// EnvLayer represents a named layer of environment variables.
// Layers are applied in order, with later layers overriding earlier ones.
type EnvLayer struct {
	Name   string
	Values map[string]string
}

// EnvBuilder provides explicit, traceable environment variable assembly.
// Each layer is named for debugging and has clear precedence.
type EnvBuilder struct {
	layers []EnvLayer
}

// NewEnvBuilder creates a new environment builder.
func NewEnvBuilder() *EnvBuilder {
	return &EnvBuilder{
		layers: make([]EnvLayer, 0, 8),
	}
}

// Add appends a named layer of environment variables.
// Later layers override earlier ones.
func (b *EnvBuilder) Add(name string, values map[string]string) *EnvBuilder {
	if values == nil {
		return b
	}
	b.layers = append(b.layers, EnvLayer{Name: name, Values: values})
	return b
}

// Build constructs the final environment variable list.
// Starts with base (typically os.Environ()) and applies all layers in order.
func (b *EnvBuilder) Build(base []string) []string {
	result := make([]string, len(base), len(base)+b.estimateSize())
	copy(result, base)

	for _, layer := range b.layers {
		result = append(result, convertEnvMapToList(layer.Values)...)
	}

	return result
}

// Layers returns a copy of the layers for debugging/logging.
func (b *EnvBuilder) Layers() []EnvLayer {
	result := make([]EnvLayer, len(b.layers))
	copy(result, b.layers)
	return result
}

func (b *EnvBuilder) estimateSize() int {
	size := 0
	for _, layer := range b.layers {
		size += len(layer.Values)
	}
	return size
}

// BuildDefaultEnv creates the default environment variables for a command.
func BuildDefaultEnv(command *config.Command, workDir string, cfg *config.Config) map[string]string {
	shell := cfg.Shell
	if command.Shell != "" {
		shell = command.Shell
	}

	cmdWorkDir := cfg.WorkDir
	if command.WorkDir != "" {
		cmdWorkDir = command.WorkDir
	}

	return map[string]string{
		"LETS_COMMAND_NAME":     command.Name,
		"LETS_COMMAND_ARGS":     strings.Join(command.Args, " "),
		"LETS_COMMAND_WORK_DIR": cmdWorkDir,
		"LETS_CONFIG":           filepath.Base(cfg.FilePath),
		"LETS_CONFIG_DIR":       filepath.Dir(cfg.FilePath),
		"LETS_SHELL":            shell,
	}
}

// BuildChecksumEnv creates environment variables for checksums.
func BuildChecksumEnv(checksumMap map[string]string) map[string]string {
	envMap := make(map[string]string, len(checksumMap))

	for name, value := range checksumMap {
		envKey := "LETS_CHECKSUM_" + normalizeEnvKey(name)
		if name == checksum.DefaultChecksumKey {
			envKey = "LETS_CHECKSUM"
		}
		envMap[envKey] = value
	}

	return envMap
}

// BuildChecksumChangedEnv creates environment variables indicating checksum changes.
func BuildChecksumChangedEnv(checksumMap, persistedChecksumMap map[string]string) map[string]string {
	envMap := make(map[string]string, len(checksumMap))

	for checksumName, checksumValue := range checksumMap {
		normalizedKey := normalizeEnvKey(checksumName)

		envKey := fmt.Sprintf("LETS_CHECKSUM_%s_CHANGED", normalizedKey)
		if checksumName == checksum.DefaultChecksumKey {
			envKey = "LETS_CHECKSUM_CHANGED"
		}

		persistedChecksum, exists := persistedChecksumMap[checksumName]
		changed := isChecksumChanged(persistedChecksum, exists, checksumValue)

		envMap[envKey] = fmt.Sprintf("%t", changed)
	}

	return envMap
}

// BuildCommandEnv assembles the complete environment for a command execution.
func BuildCommandEnv(
	command *config.Command,
	cfg *config.Config,
	cmdEnv map[string]string,
) *EnvBuilder {
	builder := NewEnvBuilder()

	// Layer 1: Default lets environment variables
	builder.Add("defaults", BuildDefaultEnv(command, cfg.WorkDir, cfg))

	// Layer 2: Global config environment
	builder.Add("config", cfg.GetEnv())

	// Layer 3: Command-specific environment
	builder.Add("command", cmdEnv)

	// Layer 4: Docopt-parsed options
	builder.Add("options", command.Options)

	// Layer 5: CLI options
	builder.Add("cli_options", command.CliOptions)

	// Layer 6: Checksum values
	builder.Add("checksum", BuildChecksumEnv(command.ChecksumMap))

	// Layer 7: Checksum changed flags (only if persist_checksum is enabled)
	if command.PersistChecksum {
		builder.Add("checksum_changed", BuildChecksumChangedEnv(
			command.ChecksumMap,
			command.GetPersistedChecksums(),
		))
	}

	return builder
}

// FormatEnvForDebug formats environment variables for debug output.
func FormatEnvForDebug(env []string) string {
	if len(env) == 0 {
		return "(empty)"
	}

	var buf strings.Builder
	for _, entry := range env {
		buf.WriteString("\n  ")
		buf.WriteString(entry)
	}
	return buf.String()
}

// FormatLayersForDebug formats environment layers for debug output.
func FormatLayersForDebug(builder *EnvBuilder) string {
	var buf strings.Builder
	for _, layer := range builder.Layers() {
		if len(layer.Values) == 0 {
			continue
		}
		buf.WriteString(fmt.Sprintf("\n  [%s]:", layer.Name))
		for k, v := range layer.Values {
			buf.WriteString(fmt.Sprintf("\n    %s=%s", k, v))
		}
	}
	return buf.String()
}

// GetBaseEnv returns the base environment (OS environment).
func GetBaseEnv() []string {
	return os.Environ()
}
