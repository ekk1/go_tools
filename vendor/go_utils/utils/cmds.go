package utils

import (
	"os/exec"
)

func RunCmd(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	ret, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(ret), nil
}
