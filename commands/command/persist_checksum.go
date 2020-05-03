package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lets-cli/lets/util"
	"github.com/lets-cli/lets/workdir"
)

const DefaultChecksumName = "lets_default_checksum"
const checksumsDir = "checksums"

func parseAndValidatePersistChecksum(persistChecksum interface{}, newCmd *Command) error {
	shouldPersist, ok := persistChecksum.(bool)

	if !ok {
		return newParseCommandError(
			"must be a bool",
			newCmd.Name,
			PersistChecksum,
			"",
		)
	}

	if !newCmd.hasChecksum {
		return newParseCommandError(
			"you must declare 'checksum' for command to use 'persist_checksum'",
			newCmd.Name,
			PersistChecksum,
			"",
		)
	}

	newCmd.PersistChecksum = shouldPersist

	return nil
}

func getCmdChecksumPath(cmdName string) string {
	return filepath.Join(workdir.DotLetsDir, checksumsDir, cmdName)
}

// returns dir path and full file path to checksum
// (.lets/checksums/[command_name]/, .lets/checksums/[command_name]/[checksum_name])
func getChecksumPath(cmdName string, checksumName string) (string, string) {
	dirPath := getCmdChecksumPath(cmdName)
	return dirPath, filepath.Join(dirPath, checksumName)
}

func PersistCommandsChecksumToDisk(cmd Command) error {
	if err := util.SafeCreateDir(filepath.Join(workdir.DotLetsDir, checksumsDir)); err != nil {
		return err
	}

	// TODO if at least one write failed do we have to revert all writes ???
	for checksumName, checksum := range cmd.ChecksumMap {
		err := persistOneChecksum(cmd.Name, checksumName, checksum)
		if err != nil {
			return err
		}
	}

	err := persistOneChecksum(cmd.Name, DefaultChecksumName, cmd.Checksum)
	if err != nil {
		return err
	}

	return nil
}

func persistOneChecksum(cmdName string, checksumName string, checksum string) error {
	checksumDirPath, checksumFilePath := getChecksumPath(cmdName, checksumName)
	if err := util.SafeCreateDir(checksumDirPath); err != nil {
		return err
	}

	f, err := os.OpenFile(checksumFilePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("can not open file %s to persist checksum: %s", checksumFilePath, err)
	}

	_, err = f.Write([]byte(checksum))
	if err != nil {
		return fmt.Errorf("can not write checksum to file %s: %s", checksumFilePath, err)
	}

	return nil
}

// ChecksumForCmdPersisted checks if checksums for cmd exists and persisted
func ChecksumForCmdPersisted(cmdName string) bool {
	// check if checksums for cmd exists
	if _, err := os.Stat(getCmdChecksumPath(cmdName)); err != nil {
		return !os.IsNotExist(err)
	}

	return true
}

// ReadChecksumsFromDisk reads all checksums for cmd into map
func ReadChecksumsFromDisk(cmdName string, checksumMap map[string]string) (map[string]string, error) {
	checksums := make(map[string]string, len(checksumMap)+1)

	for checksumName := range checksumMap {
		checksum, err := readOneChecksum(cmdName, checksumName)
		if err != nil {
			return nil, err
		}

		checksums[checksumName] = checksum
	}

	checksum, err := readOneChecksum(cmdName, DefaultChecksumName)
	if err != nil {
		return nil, err
	}

	checksums[DefaultChecksumName] = checksum

	return checksums, nil
}

func readOneChecksum(cmdName, checksumName string) (string, error) {
	_, checksumFilePath := getChecksumPath(cmdName, checksumName)

	fileData, err := ioutil.ReadFile(checksumFilePath)
	if err != nil {
		return "", fmt.Errorf("can not open file %s to read checksum: %s", checksumFilePath, err)
	}

	return string(fileData), nil
}
