package path

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/lets-cli/lets/util"
)

var ErrFileNotExists  = errors.New("file not exists")
var ErrConfigNotFound = errors.New("can not find config")


// find config file non-recursively
// filename is a file to find and work dir is where to start.
func GetFullConfigPath(filename string, workDir string) (string, error) {
	fileAbsPath, err := filepath.Abs(filepath.Join(workDir, filename))
	if err != nil {
		return "", fmt.Errorf("can not get absolute workdir path: %w", err)
	}

	if !util.FileExists(fileAbsPath) {
		return "", fmt.Errorf("%w: %s", ErrFileNotExists, fileAbsPath)
	}

	return fileAbsPath, nil
}



// find config file recursively
// filename is a file to find and work dir is where to start.
func GetFullConfigPathRecursive(filename string, workDir string) (string, error) {
	fileAbsPath, err := filepath.Abs(filepath.Join(workDir, filename))
	if err != nil {
		return "", fmt.Errorf("can not get absolute workdir path: %w", err)
	}

	if util.FileExists(fileAbsPath) {
		return fileAbsPath, nil
	}

	// else we get parent and try again up until we reach roof of fs
	parentDir := filepath.Dir(workDir)
	if parentDir == "/" {
		return "", ErrConfigNotFound
	}

	return GetFullConfigPathRecursive(filename, parentDir)
}
