package checksum

import (
	// #nosec G505
	"crypto/sha1"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lets-cli/lets/set"
	"github.com/lets-cli/lets/util"
)

const (
	DefaultChecksumKey      = "__default_checksum__"
	DefaultChecksumFileName = "lets_default_checksum"
)

var checksumCache = make(map[string][]byte)

// files can have either absolute or relative path:
// - in case of absolute path we just read that file
// - in case of relative file we trying to read file in work dir
//
// return sorted list of files read by glob patterns.
func readFilesFromPatterns(workDir string, patterns []string) ([]string, error) {
	filesSet := set.NewSet[string]()

	for _, pattern := range patterns {
		absPatternPath := pattern
		if !filepath.IsAbs(pattern) {
			absPatternPath = filepath.Join(workDir, pattern)
		}

		matches, err := filepath.Glob(absPatternPath)
		if err != nil {
			return []string{}, fmt.Errorf("can not read file to calculate checksum: %w", err)
		}

		filesSet.Add(matches...)
	}
	// sort files list
	files := filesSet.ToList()
	sort.Strings(files)

	return files, nil
}

// CalculateChecksum calculates sha1 hash from files content and return hex digest
// It calculates sha1 for each file, cache checksum for each file.
// Resulting checksum is sha1 from all files sha1's.
func CalculateChecksum(workDir string, patterns []string) (string, error) {
	// read filenames from patterns
	files, err := readFilesFromPatterns(workDir, patterns)
	if err != nil {
		return "", err
	}

	hasher := sha1.New()     // #nosec G401
	fileHasher := sha1.New() // #nosec G401

	for _, filename := range files {
		if cachedSum, found := checksumCache[filename]; found {
			_, err = hasher.Write(cachedSum)
			if err != nil {
				return "", fmt.Errorf("can not write cached checksum to hasher: %w", err)
			}
		} else {
			data, err := os.ReadFile(filename)
			if err != nil {
				return "", fmt.Errorf("can not read file to calculate checksum: %w", err)
			}
			cachedSum = fileHasher.Sum(data)
			checksumCache[filename] = cachedSum
			_, err = hasher.Write(cachedSum)
			if err != nil {
				return "", fmt.Errorf("can not write checksum to hasher: %w", err)
			}
			fileHasher.Reset()
		}
	}

	checksum := hasher.Sum(nil)

	return fmt.Sprintf("%x", checksum), nil
}

// getMapKeys get keys as array.
func getChecksumsKeys(mapping map[string][]string) []string {
	keys := make([]string, len(mapping))
	idx := 0

	for key := range mapping {
		keys[idx] = key
		idx++
	}

	return keys
}

func CalculateChecksumFromCmd(shell string, workDir string, script string) (string, error) {
	cmd := exec.Command(shell, "-c", script)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("can not calculate checksum from cmd: %s: %w", script, err)
	}

	res := string(out)
	return strings.TrimSpace(res), nil
}

// CalculateChecksumFromSources calculates checksum from checksumSources.
func CalculateChecksumFromSources(workDir string, checksumSources map[string][]string) (map[string]string, error) {
	checksumMap := make(map[string]string)

	// if checksum is a list of patterns
	if patterns, ok := checksumSources[DefaultChecksumKey]; ok {
		calcChecksum, err := CalculateChecksum(workDir, patterns)
		if err != nil {
			return nil, fmt.Errorf("calculate checksum error: %w", err)
		}

		checksumMap[DefaultChecksumKey] = calcChecksum

		return checksumMap, nil
	}

	// if checksum is a map of key: patterns
	hasher := sha1.New() // #nosec G401

	keys := getChecksumsKeys(checksumSources)
	// sort keys to make checksum deterministic
	sort.Strings(keys)

	for _, key := range keys {
		patterns := checksumSources[key]

		calcChecksum, err := CalculateChecksum(workDir, patterns)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate checksum: %w", err)
		}

		checksumMap[key] = calcChecksum

		_, err = hasher.Write([]byte(calcChecksum))
		if err != nil {
			return nil, fmt.Errorf("failed to update hasher with checksum: %w", err)
		}
	}

	checksumMap[DefaultChecksumKey] = fmt.Sprintf("%x", hasher.Sum(nil))

	return checksumMap, nil
}

func ReadChecksumFromDisk(checksumsDir, cmdName, checksumName string) (string, error) {
	_, checksumFilePath := getChecksumPath(checksumsDir, cmdName, checksumName)

	fileData, err := os.ReadFile(checksumFilePath)
	if err != nil {
		return "", fmt.Errorf("can not open file %s to read checksum: %w", checksumFilePath, err)
	}

	return string(fileData), nil
}

func getCmdChecksumPath(checksumsDir string, cmdName string) string {
	return filepath.Join(checksumsDir, cmdName)
}

// returns dir path and full file path to checksum
// (.lets/checksums/[command_name]/, .lets/checksums/[command_name]/[checksum_name]).
func getChecksumPath(checksumsDir string, cmdName string, checksumName string) (string, string) {
	dirPath := getCmdChecksumPath(checksumsDir, cmdName)

	return dirPath, filepath.Join(dirPath, checksumName)
}

// TODO maybe checksumMap has to be separate struct ?
func PersistCommandsChecksumToDisk(checksumsDir string, checksumMap map[string]string, cmdName string) error {
	// TODO if at least one write failed do we have to revert all writes ???
	for checksumName, checksum := range checksumMap {
		filename := checksumName
		if checksumName == DefaultChecksumKey {
			filename = DefaultChecksumFileName
		}
		err := persistOneChecksum(checksumsDir, cmdName, filename, checksum)
		if err != nil {
			return err
		}
	}

	return nil
}

func persistOneChecksum(checksumsDir string, cmdName string, checksumName string, checksum string) error {
	checksumDirPath, checksumFilePath := getChecksumPath(checksumsDir, cmdName, checksumName)
	if err := util.SafeCreateDir(checksumDirPath); err != nil {
		return fmt.Errorf("can not create checksum dir at %s: %w", checksumDirPath, err)
	}

	f, err := os.OpenFile(checksumFilePath, os.O_CREATE|os.O_WRONLY, 0o755) //nolint:nosnakecase
	if err != nil {
		return fmt.Errorf("can not open file %s to persist checksum: %w", checksumFilePath, err)
	}

	_, err = f.Write([]byte(checksum))
	if err != nil {
		return fmt.Errorf("can not write checksum to file %s: %w", checksumFilePath, err)
	}

	return nil
}

// IsChecksumForCmdPersisted checks if checksums for cmd exists and persisted.
func IsChecksumForCmdPersisted(checksumsDir string, cmdName string) bool {
	// check if checksums for cmd exists
	if _, err := os.Stat(getCmdChecksumPath(checksumsDir, cmdName)); err != nil {
		return !os.IsNotExist(err)
	}

	return true
}
