package util

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func OpenEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return errors.New("EDITOR is not set")
	}

	cmd := exec.Command(editor, path) //nolint:gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %s: %w", editor, err)
	}

	return nil
}
