package runner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lets-cli/lets/commands/command"
)

func makeEnvEntry(k, v string) string {
	return fmt.Sprintf("%s=%s", k, v)
}

func normalizeEnvKey(origKey string) string {
	key := strings.ReplaceAll(origKey, "-", "_")
	key = strings.ToUpper(key)

	return key
}

func convertEnvMapToList(envMap map[string]string) []string {
	var envList []string
	for name, value := range envMap {
		envList = append(envList, makeEnvEntry(name, value))
	}

	return envList
}

func convertChecksumToEnvForCmd(checksum string) []string {
	return []string{makeEnvEntry("LETS_CHECKSUM", checksum)}
}

func convertChecksumMapToEnvForCmd(checksumMap map[string]string) []string {
	var envList []string

	for name, value := range checksumMap {
		if name != "" {
			envList = append(envList, makeEnvEntry(fmt.Sprintf("LETS_CHECKSUM_%s", normalizeEnvKey(name)), value))
		}
	}

	return envList
}

// persistedChecksumMap can be empty, and if so, we set env var LETS_CHECKSUM_[NAME]_CHANGED to false for all checksums
func convertChangedChecksumMapToEnvForCmd(
	defaultChecksum string,
	checksumMap map[string]string,
	persistedChecksumMap map[string]string,
) []string {
	var envList []string

	for name, value := range checksumMap {
		if name == "" { // TODO do we still have this empty key
			continue
		}

		normalizedKey := normalizeEnvKey(name)
		persistedValue, ok := persistedChecksumMap[name]
		checksumChanged := false

		if ok {
			checksumChanged = value != persistedValue
		}

		envList = append(
			envList,
			makeEnvEntry(fmt.Sprintf("LETS_CHECKSUM_%s_CHANGED", normalizedKey), strconv.FormatBool(checksumChanged)),
		)
	}

	persistedValue, ok := persistedChecksumMap[command.DefaultChecksumName]

	defaultChecksumChanged := false
	if ok {
		defaultChecksumChanged = defaultChecksum != persistedValue
	}

	envList = append(
		envList,
		makeEnvEntry("LETS_CHECKSUM_CHANGED", strconv.FormatBool(defaultChecksumChanged)),
	)

	return envList
}

func composeEnvs(envs ...[]string) []string {
	var composed []string
	for _, env := range envs {
		composed = append(composed, env...)
	}

	return composed
}
