package config

import (
	// #nosec G505
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/lets-cli/lets/util"
)

const (
	DefaultChecksumName = "lets_default_checksum"
	checksumsDir        = "checksums"
)

var checksumCache = make(map[string][]byte)

// files can have either absolute or relative path:
// - in case of absolute path we just read that file
// - in case of relative file we trying to read file in work dir
//
// return sorted list of files read by glob patterns.
func readFilesFromPatterns(workDir string, patterns []string) ([]string, error) {
	var files []string

	for _, pattern := range patterns {
		absPatternPath := pattern
		if !filepath.IsAbs(pattern) {
			absPatternPath = filepath.Join(workDir, pattern)
		}

		matches, err := filepath.Glob(absPatternPath)
		if err != nil {
			return []string{}, fmt.Errorf("can not read file to calculate checksum: %w", err)
		}

		files = append(files, matches...)
	}
	// sort files list
	sort.Strings(files)

	return files, nil
}

// calculate sha1 hash from files content and return hex digest
// It calculates sha1 for each file, cache checksum for each file.
// Resulting checksum is sha1 from all files sha1's.
func calculateChecksum(workDir string, patterns []string) (string, error) {
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

func getChecksumsKeys(mapping map[string][]string) []string {
	keys := make([]string, len(mapping))
	idx := 0

	for key := range mapping {
		keys[idx] = key
		idx++
	}

	return keys
}

// calculate checksum from files listed in command.checksum.
func calculateChecksumFromSource(workDir string, newCmd *Command) error {
	newCmd.ChecksumMap = make(map[string]string)
	// if checksum is a list of patterns
	if patterns, ok := newCmd.ChecksumSource[""]; ok {
		calcChecksum, err := calculateChecksum(workDir, patterns)
		if err != nil {
			return fmt.Errorf("calculate checksum error: %w", err)
		}

		newCmd.Checksum = calcChecksum

		return nil
	}

	// if checksum is a map of key: patterns
	hasher := sha1.New() // #nosec G401

	keys := getChecksumsKeys(newCmd.ChecksumSource)
	// sort keys to make checksum deterministic
	sort.Strings(keys)

	for _, key := range keys {
		patterns := newCmd.ChecksumSource[key]

		calcChecksum, err := calculateChecksum(workDir, patterns)
		if err != nil {
			return fmt.Errorf("failed to calculate checksum: %w", err)
		}

		newCmd.ChecksumMap[key] = calcChecksum

		_, err = hasher.Write([]byte(calcChecksum))
		if err != nil {
			return fmt.Errorf("failed to update hasher with checksum: %w", err)
		}
	}

	newCmd.Checksum = fmt.Sprintf("%x", hasher.Sum(nil))

	return nil
}

func readOneChecksum(dotLetsDir, cmdName, checksumName string) (string, error) {
	_, checksumFilePath := getChecksumPath(dotLetsDir, cmdName, checksumName)

	fileData, err := os.ReadFile(checksumFilePath)
	if err != nil {
		return "", fmt.Errorf("can not open file %s to read checksum: %w", checksumFilePath, err)
	}

	return string(fileData), nil
}

func getCmdChecksumPath(dotLetsDir string, cmdName string) string {
	return filepath.Join(dotLetsDir, checksumsDir, cmdName)
}

// returns dir path and full file path to checksum
// (.lets/checksums/[command_name]/, .lets/checksums/[command_name]/[checksum_name]).
func getChecksumPath(dotLetsDir string, cmdName string, checksumName string) (string, string) {
	dirPath := getCmdChecksumPath(dotLetsDir, cmdName)

	return dirPath, filepath.Join(dirPath, checksumName)
}

func PersistCommandsChecksumToDisk(dotLetsDir string, cmd Command) error {
	checksumPath := filepath.Join(dotLetsDir, checksumsDir)
	if err := util.SafeCreateDir(checksumPath); err != nil {
		return fmt.Errorf("can not create %s: %w", checksumPath, err)
	}

	// TODO if at least one write failed do we have to revert all writes ???
	for checksumName, checksum := range cmd.ChecksumMap {
		err := persistOneChecksum(dotLetsDir, cmd.Name, checksumName, checksum)
		if err != nil {
			return err
		}
	}

	err := persistOneChecksum(dotLetsDir, cmd.Name, DefaultChecksumName, cmd.Checksum)
	if err != nil {
		return err
	}

	return nil
}

func persistOneChecksum(dotLetsDir string, cmdName string, checksumName string, checksum string) error {
	checksumDirPath, checksumFilePath := getChecksumPath(dotLetsDir, cmdName, checksumName)
	if err := util.SafeCreateDir(checksumDirPath); err != nil {
		return fmt.Errorf("can not create checksum dir at %s: %w", checksumDirPath, err)
	}

	f, err := os.OpenFile(checksumFilePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("can not open file %s to persist checksum: %w", checksumFilePath, err)
	}

	_, err = f.Write([]byte(checksum))
	if err != nil {
		return fmt.Errorf("can not write checksum to file %s: %w", checksumFilePath, err)
	}

	return nil
}

// ChecksumForCmdPersisted checks if checksums for cmd exists and persisted.
func ChecksumForCmdPersisted(dotLetsDir string, cmdName string) bool {
	// check if checksums for cmd exists
	if _, err := os.Stat(getCmdChecksumPath(dotLetsDir, cmdName)); err != nil {
		return !os.IsNotExist(err)
	}

	return true
}
