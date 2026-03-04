package test

import (
	"os"
)

// MockArgs mocks os.Args with values passed to thi func.
func MockArgs(args []string) {
	os.Args = append([]string{"lets"}, args...)
}
