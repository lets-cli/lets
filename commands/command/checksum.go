package command

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
)

var checksumCache map[string][]byte = make(map[string][]byte)

// calculate sha1 hash from files content and return hex digest
func calculateChecksum(patterns []string) (string, error) {
	// read filenames from patterns
	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return "", err
		}
		files = append(files, matches...)
	}
	// sort files list
	sort.Strings(files)
	hasher := sha1.New()
	fileHasher := sha1.New()
	for _, filename := range files {
		if cachedSum, found := checksumCache[filename]; found {
			hasher.Write(cachedSum)
		} else {
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				return "", err
			}
			cachedSum = fileHasher.Sum(data)
			checksumCache[filename] = cachedSum
			hasher.Write(cachedSum)
			fileHasher.Reset()
		}

	}

	checksum := hasher.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}

func parseAndValidateChecksum(checksum interface{}, newCmd *Command) error {
	patternsList, okList := checksum.([]interface{})
	patternsMap, okMap := checksum.(map[interface{}]interface{})
	checksumSource := make(map[string][]string)
	if okList {
		for _, value := range patternsList {
			if value, ok := value.(string); ok {
				checksumSource[""] = append(checksumSource[""], value)
			} else {
				return newCommandError(
					"value of checksum list must be a string",
					newCmd.Name,
					CHECKSUM,
					"",
				)
			}
		}
	} else if okMap {
		for key, patterns := range patternsMap {
			key, ok := key.(string)
			if !ok {
				return newCommandError(
					"key of checksum list must be a string",
					newCmd.Name,
					CHECKSUM,
					"",
				)
			}
			patterns, ok := patterns.([]interface{})
			if !ok {
				return newCommandError(
					"value of checksum map must be a list",
					newCmd.Name,
					CHECKSUM,
					"",
				)
			}
			for _, value := range patterns {
				if value, ok := value.(string); ok {
					checksumSource[key] = append(checksumSource[key], value)
				} else {
					return newCommandError(
						"value of checksum list must be a string",
						newCmd.Name,
						CHECKSUM,
						"",
					)
				}
			}
		}
	} else {
		return newCommandError(
			"must be a list of string (files of glob patterns) or a map of lists of string",
			newCmd.Name,
			CHECKSUM,
			"",
		)
	}
	newCmd.ChecksumCalculator = func() error {
		return calculateChecksumFromSource(newCmd, checksumSource)
	}
	return nil
}

func calculateChecksumFromSource(newCmd *Command, checksumSource map[string][]string) error {
	newCmd.ChecksumMap = make(map[string]string)
	hasher := sha1.New()
	for key, patterns := range checksumSource {
		calcChecksum, err := calculateChecksum(patterns)
		if err != nil {
			return fmt.Errorf("failed to calculate checksum: %s", err)
		} else {
			newCmd.ChecksumMap[key] = calcChecksum
		}
		hasher.Write([]byte(calcChecksum))
	}
	newCmd.Checksum = fmt.Sprintf("%x", hasher.Sum(nil))

	return nil
}
