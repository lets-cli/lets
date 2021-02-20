package test

import (
	"log"
	"os"
)

func CreateTempFile(dir string, name string) *os.File {
	file, err := os.CreateTemp(dir, name)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
