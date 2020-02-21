package command

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
)

var checkSumCache map[string][]byte = make(map[string][]byte)

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
		if cachedSum, found := checkSumCache[filename]; found {
			hasher.Write(cachedSum)
		} else {
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				return "", err
			}
			cachedSum = fileHasher.Sum(data)
			checkSumCache[filename] = cachedSum
			hasher.Write(cachedSum)
			fileHasher.Reset()
		}

	}

	checksum := hasher.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}

func checkSumFromList(cmdName string, patterns []interface{}) (string, error) {
	var files []string
	for _, value := range patterns {
		if value, ok := value.(string); ok {
			files = append(files, value)
		} else {
			return "", newCommandError(
				"value of checksum list must be a string",
				cmdName,
				CHECKSUM,
				"",
			)
		}
	}
	calcChecksum, err := calculateChecksum(files)
	if err == nil {
		return calcChecksum, nil
	} else {
		return "", err
	}
}

// TODO make checksum calculation lazy
func parseAndValidateChecksum(checksum interface{}, newCmd *Command) error {
	patternsList, okList := checksum.([]interface{})
	patternsMap, okMap := checksum.(map[interface{}]interface{})
	if okList{
		calcChecksum, err := checkSumFromList(newCmd.Name, patternsList)
		if err != nil {
			return fmt.Errorf("failed to calculate checksum: %s", err)
		} else {
			newCmd.Checksum = calcChecksum
		}
	} else if okMap {
		hasher := sha1.New()
		newCmd.ChecksumMap = make(map[string]string)
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
			calcChecksum, err := checkSumFromList(newCmd.Name, patterns)
			if err != nil {
				return fmt.Errorf("failed to calculate checksum: %s", err)
			} else {
				newCmd.ChecksumMap[key] = calcChecksum
			}
				hasher.Write([]byte(calcChecksum))
			}
		newCmd.Checksum = fmt.Sprintf("%x", hasher.Sum(nil))
	} else {
		return newCommandError(
			"must be a list of string (files of glob patterns) or a map of lists of string",
			newCmd.Name,
			CHECKSUM,
			"",
		)
	}
	return nil
}
