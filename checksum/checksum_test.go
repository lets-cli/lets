package checksum

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

	checksum, err := CalculateChecksum(tempDir, []string{
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
	checksum, err := CalculateChecksum(tempDir, []string{
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
	f1Prefix := getFilePrefix(file1.Name())
	f2Prefix := getFilePrefix(file2.Name())
	checksumSources := map[string][]string{
		DefaultChecksumKey: {f1Prefix, f2Prefix},
	}

	checksumMap, err := CalculateChecksumFromSources(tempDir, checksumSources)
	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	if checksumMap[DefaultChecksumKey] != expectChecksum {
		t.Errorf(
			"Checksum is not correct for command with checksum as list. Expect: %s, got: %s",
			expectChecksum,
			checksumMap[DefaultChecksumKey],
		)
	}

	// declare command with checksum as map but with same files
	checksumSources1 := map[string][]string{
		"misc": {f1Prefix, f2Prefix},
	}

	checksumMap1, err := CalculateChecksumFromSources(tempDir, checksumSources1)

	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	if checksumMap1["misc"] != expectChecksum {
		t.Errorf(
			"Checksum is not correct for command with checksum as map. Expect: %s, got: %s",
			expectChecksum,
			checksumMap1["misc"],
		)
	}
}
