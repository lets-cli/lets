package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const defaultChecksum = "da39a3ee5e6b4b0d3255bfef95601890afd80709"

func createTempFile(dir string, name string) *os.File {
	file, err := ioutil.TempFile(dir, name)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func TestCalculateChecksumSimpleFilename(t *testing.T) {
	tempDir := os.TempDir()
	file1 := createTempFile(tempDir, "lets_checksum_test_1")
	file2 := createTempFile(tempDir, "lets_checksum_test_2")

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

	checksum, err := calculateChecksum([]string{
		file1.Name(),
		file2.Name(),
	})

	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	if defaultChecksum != checksum {
		t.Errorf("Checksum is not correct. Expect: %s, got: %s", defaultChecksum, checksum)
	}
}

func TestCalculateChecksumGlobPattern(t *testing.T) {
	tempDir := os.TempDir()
	file1 := createTempFile(tempDir, "lets_checksum_test_1")
	file2 := createTempFile(tempDir, "lets_checksum_test_2")

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

	checksum, err := calculateChecksum([]string{
		fmt.Sprintf("%s/lets_checksum_test_*", tempDir),
	})

	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	if defaultChecksum != checksum {
		t.Errorf("Checksum is not correct. Expect: %s, got: %s", defaultChecksum, checksum)
	}
}

func TestCalculateChecksumFromListOrMap(t *testing.T) {
	tempDir := os.TempDir()
	file1 := createTempFile(tempDir, "lets_checksum_test_1")
	file2 := createTempFile(tempDir, "lets_checksum_test_2")

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
	cmdChAsList.checksumSource = map[string][]string{
		"": {"lets_checksum_test_1", "lets_checksum_test_2"},
	}

	err = calculateChecksumFromSource(&cmdChAsList)
	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	if cmdChAsList.Checksum != defaultChecksum {
		t.Errorf(
			"Checksum is not correct for command with checksum as list. Expect: %s, got: %s",
			defaultChecksum,
			cmdChAsList.Checksum,
		)
	}

	// declare command with checksum as map but with same files
	cmdChAsMap := NewCommand("checksum-as-map")
	cmdChAsMap.checksumSource = map[string][]string{
		"misc": {"lets_checksum_test_1", "lets_checksum_test_2"},
	}

	err = calculateChecksumFromSource(&cmdChAsMap)
	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	if cmdChAsMap.ChecksumMap["misc"] != defaultChecksum {
		t.Errorf(
			"Checksum is not correct for command with checksum as map. Expect: %s, got: %s",
			defaultChecksum,
			cmdChAsMap.ChecksumMap["misc"],
		)
	}
}
