package executor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lets-cli/lets/checksum"
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

func getChecksumEnvMap(checksumMap map[string]string) map[string]string {
	envMap := make(map[string]string)

	for name, value := range checksumMap {
		envKey := fmt.Sprintf("LETS_CHECKSUM_%s", normalizeEnvKey(name))
		if name == checksum.DefaultChecksumKey {
			envKey = "LETS_CHECKSUM"
		}
		envMap[envKey] = value
	}

	return envMap
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
func getChangedChecksumEnvMap(
	checksumMap map[string]string,
	persistedChecksumMap map[string]string,
) map[string]string {
	envMap := make(map[string]string)

	for checksumName, checksumValue := range checksumMap {
		normalizedKey := normalizeEnvKey(checksumName)

		envKey := fmt.Sprintf("LETS_CHECKSUM_%s_CHANGED", normalizedKey)
		if checksumName == checksum.DefaultChecksumKey {
			envKey = "LETS_CHECKSUM_CHANGED"
		}

		persistedChecksum, persistedChecksumExists := persistedChecksumMap[checksumName]

		checksumChanged := isChecksumChanged(persistedChecksum, persistedChecksumExists, checksumValue)

		envMap[envKey] = strconv.FormatBool(checksumChanged)
	}

	return envMap
}

func fmtEnv(env []string) string {
	buf := ""

	for _, entry := range env {
		buf = fmt.Sprintf("%s\n  %s", buf, entry)
	}

	return buf
}
