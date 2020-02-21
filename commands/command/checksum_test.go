package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

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

	file1.Write([]byte("qwerty1"))
	file2.Write([]byte("asdfg2"))

	checksum, err := calculateChecksum([]string{
		file1.Name(),
		file2.Name(),
	})

	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	expected := "56a89c168888554d9cafa50c2f37c249dde6e37d"
	if expected != checksum {
		t.Errorf("Checksum is not correct. Expect: %s, got: %s", expected, checksum)
	}
}

func TestCalculateChecksumGlobPattern(t *testing.T) {
	tempDir := os.TempDir()
	file1 := createTempFile(tempDir, "lets_checksum_test_1")
	file2 := createTempFile(tempDir, "lets_checksum_test_2")

	defer os.Remove(file1.Name())
	defer os.Remove(file2.Name())

	file1.Write([]byte("qwerty1"))
	file2.Write([]byte("asdfg2"))

	checksum, err := calculateChecksum([]string{
		fmt.Sprintf("%s/lets_checksum_test_*", tempDir),
	})

	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	expected := "56a89c168888554d9cafa50c2f37c249dde6e37d"
	if expected != checksum {
		t.Errorf("Checksum is not correct. Expect: %s, got: %s", expected, checksum)
	}
}
