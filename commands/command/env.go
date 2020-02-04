package command

import "os/exec"

func evalEnvVariable(rawCmd string) (string, error) {
	cmd := exec.Command("sh", "-c", rawCmd)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
