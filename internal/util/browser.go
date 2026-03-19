package util

import (
	"fmt"
	"os/exec"
	"runtime"
)

func browserCommand(goos string, url string) (*exec.Cmd, error) {
	switch goos {
	case "darwin":
		return exec.Command("open", url), nil
	case "linux":
		return exec.Command("xdg-open", url), nil
	default:
		return nil, fmt.Errorf("unsupported platform %q", goos)
	}
}

func OpenURL(url string) error {
	cmd, err := browserCommand(runtime.GOOS, url)
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start %s: %w", cmd.Path, err)
	}

	if cmd.Process != nil {
		_ = cmd.Process.Release()
	}

	return nil
}
