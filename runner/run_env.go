package runner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lets-cli/lets/config/config"
)

func makeEnvEntry(k, v string) string {
	return fmt.Sprintf("%s=%s", k, v)
}

// helper function for convenient use with composeEnvs.
func makeEnvEntryList(k, v string) []string {
	return []string{makeEnvEntry(k, v)}
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

func isChecksumChanged(persistedChecksum string, persistedChecksumExists bool, newChecksum string) bool {
	if !persistedChecksumExists {
		// We set true here because if there was no persisted checksum that means that its a brand new checksum.
		// Hence it was changed from none to some value.
		return true
	}

	// But if we have persisted checksum - we check for checksum change below.
	return persistedChecksum != newChecksum
}

// persistedChecksumMap can be empty, and if so, we set env var LETS_CHECKSUM_[NAME]_CHANGED to false for all checksums.
func convertChangedChecksumMapToEnvForCmd(
	defaultChecksum string,
	checksumMap map[string]string,
	persistedChecksumMap map[string]string,
) []string {
	var envList []string

	for checksumName, checksumValue := range checksumMap {
		if checksumName == "" { // TODO do we still have this empty key
			continue
		}

		normalizedKey := normalizeEnvKey(checksumName)
		persistedChecksum, persistedChecksumExists := persistedChecksumMap[checksumName]

		checksumChanged := isChecksumChanged(persistedChecksum, persistedChecksumExists, checksumValue)

		envList = append(
			envList,
			makeEnvEntry(fmt.Sprintf("LETS_CHECKSUM_%s_CHANGED", normalizedKey), strconv.FormatBool(checksumChanged)),
		)
	}

	persistedDefaultChecksum, persistedDefaultChecksumExists := persistedChecksumMap[config.DefaultChecksumName]

	defaultChecksumChanged := isChecksumChanged(persistedDefaultChecksum, persistedDefaultChecksumExists, defaultChecksum)

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
