package utils

import (
	"os/exec"
)

// RunCmd Runs a cmd, env can be nil or a map
func RunCmd(command string, env map[string]string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	if env != nil {
		for k, v := range env {
			cmd.Env = append(cmd.Environ(), k+"="+v)
		}
	}
	ret, err := cmd.CombinedOutput()
	if err != nil {
		return string(ret), err
	}
	return string(ret), nil
}
