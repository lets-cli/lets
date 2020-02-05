package command

import (
	"crypto/sha1"
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
	fmt.Printf("files before sort: %s\n", files)
	// sort files list
	sort.Strings(files)
	fmt.Printf("files after sort: %s\n", files)
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
