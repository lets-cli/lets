package util

import (
	"fmt"
	"os/exec"
	"runtime"
)

func browserCommand(goos string, url string) (string, []string, error) {
	switch goos {
	case "darwin":
		return "open", []string{url}, nil
	case "linux":
		return "xdg-open", []string{url}, nil
	default:
		return "", nil, fmt.Errorf("unsupported platform %q", goos)
	}
}

func OpenURL(url string) error {
	name, args, err := browserCommand(runtime.GOOS, url)
	if err != nil {
		return err
	}

	cmd := exec.Command(name, args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start %s: %w", name, err)
	}

	if cmd.Process != nil {
		_ = cmd.Process.Release()
	}

	return nil
}
