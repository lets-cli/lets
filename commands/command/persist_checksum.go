package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lets-cli/lets/util"
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

func getCmdChecksumPath(dotLetsDir string, cmdName string) string {
	return filepath.Join(dotLetsDir, checksumsDir, cmdName)
}

// returns dir path and full file path to checksum
// (.lets/checksums/[command_name]/, .lets/checksums/[command_name]/[checksum_name])
func getChecksumPath(dotLetsDir string, cmdName string, checksumName string) (string, string) {
	dirPath := getCmdChecksumPath(dotLetsDir, cmdName)
	return dirPath, filepath.Join(dirPath, checksumName)
}

func PersistCommandsChecksumToDisk(dotLetsDir string, cmd Command) error {
	if err := util.SafeCreateDir(filepath.Join(dotLetsDir, checksumsDir)); err != nil {
		return err
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
func ChecksumForCmdPersisted(dotLetsDir string, cmdName string) bool {
	// check if checksums for cmd exists
	if _, err := os.Stat(getCmdChecksumPath(dotLetsDir, cmdName)); err != nil {
		return !os.IsNotExist(err)
	}

	return true
}

// ReadChecksumsFromDisk reads all checksums for cmd into map
func (cmd *Command) ReadChecksumsFromDisk(dotLetsDir string, cmdName string, checksumMap map[string]string) error {
	checksums := make(map[string]string, len(checksumMap)+1)

	for checksumName := range checksumMap {
		checksum, err := readOneChecksum(dotLetsDir, cmdName, checksumName)
		if err != nil {
			return err
		}

		checksums[checksumName] = checksum
	}

	checksum, err := readOneChecksum(dotLetsDir, cmdName, DefaultChecksumName)
	if err != nil {
		return err
	}

	checksums[DefaultChecksumName] = checksum

	cmd.persistedChecksums = checksums

	return nil
}

func readOneChecksum(dotLetsDir, cmdName, checksumName string) (string, error) {
	_, checksumFilePath := getChecksumPath(dotLetsDir, cmdName, checksumName)

	fileData, err := os.ReadFile(checksumFilePath)
	if err != nil {
		return "", fmt.Errorf("can not open file %s to read checksum: %s", checksumFilePath, err)
	}

	return string(fileData), nil
}
