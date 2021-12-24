package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/lets-cli/lets/test"
)

// get filename without last 4 chars - to make tests more predictable.
func getFilePrefix(filename string) string {
	return fmt.Sprintf("%s*", filename[:len(filename)-4])
}

const expectChecksum = "56a89c168888554d9cafa50c2f37c249dde6e37d"

func TestCalculateChecksumSimpleFilename(t *testing.T) {
	tempDir := os.TempDir()
	file1 := test.CreateTempFile(tempDir, "lets_checksum_test_1")
	file2 := test.CreateTempFile(tempDir, "lets_checksum_test_2")

	defer os.Remove(file1.Name())
	defer os.Remove(file2.Name())

	_, err := file1.Write([]byte("qwerty1"))
	if err != nil {
		t.Errorf("Can not write test file. Error: %s", err)
	}

	_, err = file2.Write([]byte("asdfg2"))
	if err != nil {
		t.Errorf("Can not write test file. Error: %s", err)
	}

	checksum, err := calculateChecksum(tempDir, []string{
		file1.Name(),
		file2.Name(),
	})
	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	if expectChecksum != checksum {
		t.Errorf("Checksum is not correct. Expect: %s, got: %s", expectChecksum, checksum)
	}
}

func TestCalculateChecksumGlobPattern(t *testing.T) {
	tempDir := os.TempDir()
	file1 := test.CreateTempFile(tempDir, "lets_checksum_test_1")
	file2 := test.CreateTempFile(tempDir, "lets_checksum_test_2")

	defer os.Remove(file1.Name())
	defer os.Remove(file2.Name())

	_, err := file1.Write([]byte("qwerty1"))
	if err != nil {
		t.Errorf("Can not write test file. Error: %s", err)
	}

	_, err = file2.Write([]byte("asdfg2"))
	if err != nil {
		t.Errorf("Can not write test file. Error: %s", err)
	}

	f1Prefix := getFilePrefix(file1.Name())
	f2Prefix := getFilePrefix(file2.Name())
	checksum, err := calculateChecksum(tempDir, []string{
		f1Prefix,
		f2Prefix,
	})
	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	if expectChecksum != checksum {
		t.Errorf("Checksum is not correct. Expect: %s, got: %s", expectChecksum, checksum)
	}
}

func TestCalculateChecksumFromListOrMap(t *testing.T) {
	tempDir := os.TempDir()
	file1 := test.CreateTempFile(tempDir, "lets_checksum_test_1")
	file2 := test.CreateTempFile(tempDir, "lets_checksum_test_2")

	defer os.Remove(file1.Name())
	defer os.Remove(file2.Name())

	_, err := file1.Write([]byte("qwerty1"))
	if err != nil {
		t.Errorf("Can not write test file. Error: %s", err)
	}

	_, err = file2.Write([]byte("asdfg2"))
	if err != nil {
		t.Errorf("Can not write test file. Error: %s", err)
	}

	// declare command with checksum as list
	cmdChAsList := NewCommand("checksum-as-list")
	f1Prefix := getFilePrefix(file1.Name())
	f2Prefix := getFilePrefix(file2.Name())
	cmdChAsList.ChecksumSource = map[string][]string{
		"": {f1Prefix, f2Prefix},
	}

	err = calculateChecksumFromSource(tempDir, &cmdChAsList)
	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	if cmdChAsList.Checksum != expectChecksum {
		t.Errorf(
			"Checksum is not correct for command with checksum as list. Expect: %s, got: %s",
			expectChecksum,
			cmdChAsList.Checksum,
		)
	}

	// declare command with checksum as map but with same files
	cmdChAsMap := NewCommand("checksum-as-map")
	cmdChAsMap.ChecksumSource = map[string][]string{
		"misc": {f1Prefix, f2Prefix},
	}

	err = calculateChecksumFromSource(tempDir, &cmdChAsMap)
	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	if cmdChAsMap.ChecksumMap["misc"] != expectChecksum {
		t.Errorf(
			"Checksum is not correct for command with checksum as map. Expect: %s, got: %s",
			expectChecksum,
			cmdChAsMap.ChecksumMap["misc"],
		)
	}
}
