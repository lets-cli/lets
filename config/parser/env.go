package parser

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/lets-cli/lets/checksum"
	"github.com/lets-cli/lets/config/config"
)

type EnvEntry interface {
	Value() (string, error)
	Name() string
}

type EnvString struct {
	name  string
	value string
}

func (entry EnvString) Value() (string, error) {
	return entry.value, nil
}

func (entry EnvString) Name() string {
	return entry.name
}

type EnvSh struct {
	name   string
	script string
}

func (entry EnvSh) Value() (string, error) {
	computedVal, err := entry.executeScript(entry.script)
	if err != nil {
		return "", parseDirectiveError(
			"env",
			fmt.Sprintf("failed to eval '%s' env variable: %s", entry.name, err),
		)
	}

	return computedVal, nil
}

// eval env value and trim result string.
func (entry EnvSh) executeScript(script string) (string, error) {
	// TODO maybe use cfg.Shell instead of sh.
	// TODO pass env from cfg.env - it will allow to use static env in eval_env
	cmd := exec.Command("sh", "-c", script)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("can not get output from eval_env script: %s: %w", script, err)
	}

	res := string(out)
	// TODO get rid of TrimSpace
	return strings.TrimSpace(res), nil
}

func (entry EnvSh) Name() string {
	return entry.name
}

type EnvChecksum struct {
	name     string
	patterns []string
	workDir  string
}

func (entry EnvChecksum) Value() (string, error) {
	checksumResult, err := checksum.CalculateChecksum(entry.workDir, entry.patterns)
	if err != nil {
		return "", err
	}

	return checksumResult, nil
}

func (entry EnvChecksum) Name() string {
	return entry.name
}

func parseEnvEntries(rawEnv map[string]interface{}, cfg *config.Config) ([]EnvEntry, error) {
	var envEntries []EnvEntry

	for name, rawValue := range rawEnv {
		switch value := rawValue.(type) {
		case string:
			envEntries = append(envEntries, EnvString{name, value})
		case int:
			envEntries = append(envEntries, EnvString{name: name, value: fmt.Sprintf("%d", value)})
		case map[string]interface{}:
			for mode, modeValue := range value {
				switch mode {
				case "sh":
					envEntries = append(envEntries, EnvSh{name: name, script: fmt.Sprintf("%s", modeValue)})
				case "checksum":
					patternsList, ok := modeValue.([]interface{})
					errMsg := fmt.Sprintf(
						"failed to parse checksum patterns list for '%s' env variable: must be list of strings",
						name,
					)
					if !ok {
						return []EnvEntry{}, parseDirectiveError("env", errMsg)
					}
					patterns := make([]string, 0, len(patternsList))

					for _, value := range patternsList {
						if value, ok := value.(string); ok {
							patterns = append(patterns, value)
						} else {
							return []EnvEntry{}, parseDirectiveError("env", errMsg)
						}
					}
					envEntries = append(envEntries, EnvChecksum{name: name, patterns: patterns, workDir: cfg.WorkDir})
				default:
					return []EnvEntry{}, parseDirectiveError(
						"env",
						fmt.Sprintf("invalid execution mode '%s' for '%s' env variable", mode, name),
					)
				}
			}
		}
	}

	return envEntries, nil
}
