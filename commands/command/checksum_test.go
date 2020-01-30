package command

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func createTempFile(name string) *os.File {
	file, err := ioutil.TempFile("", name)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func TestCalculateChecksumSimpleFilename(t *testing.T) {
	file1 := createTempFile("lets_checksum_test_1")
	file2 := createTempFile("lets_checksum_test_2")

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

	expected := "f064d5c8f0433a574f6f5cdef5b18b850ac4029e"
	if expected != checksum {
		t.Errorf("Checksum is not correct. Expect: %s, got: %s", expected, checksum)
	}
}

func TestCalculateChecksumGlobPattern(t *testing.T) {
	file1 := createTempFile("lets_checksum_test_1")
	file2 := createTempFile("lets_checksum_test_2")

	defer os.Remove(file1.Name())
	defer os.Remove(file2.Name())

	file1.Write([]byte("qwerty1"))
	file2.Write([]byte("asdfg2"))

	checksum, err := calculateChecksum([]string{
		"/tmp/lets_checksum_test_*",
	})

	if err != nil {
		t.Errorf("Checksum is not correct. Error: %s", err)
	}

	expected := "f064d5c8f0433a574f6f5cdef5b18b850ac4029e"
	if expected != checksum {
		t.Errorf("Checksum is not correct. Expect: %s, got: %s", expected, checksum)
	}
}
