package command

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
)

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
	for _, filename := range files {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", err
		}
		hasher.Write(data)
	}

	checksum := hasher.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}

func parseAndValidateChecksum(checksum interface{}, newCmd *Command) error {
	patterns, ok := checksum.([]interface{})
	if !ok {
		return newCommandError(
			"must be a list of string (files of glob patterns)",
			newCmd.Name,
			CHECKSUM,
			"",
		)
	}

	var files []string
	for _, value := range patterns {
		if value, ok := value.(string); ok {
			files = append(files, value)
		} else {
			return newCommandError(
				"value of checksum list must be a string",
				newCmd.Name,
				CHECKSUM,
				"",
			)
		}
	}
	calcChecksum, err := calculateChecksum(files)
	if err == nil {
		newCmd.Checksum = calcChecksum
	} else {
		return errors.New(fmt.Sprintf("failed to calculate checksum: %s", err))
	}
	return nil
}
