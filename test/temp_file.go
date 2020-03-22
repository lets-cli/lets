package test

import (
	"io/ioutil"
	"log"
	"os"
)

func CreateTempFile(dir string, name string) *os.File {
	file, err := ioutil.TempFile(dir, name)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
