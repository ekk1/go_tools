package utils

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func ReadUserInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	ret, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	ret = strings.TrimSuffix(ret, "\n")
	return ret, nil
}

func AskUserConfirm(prompt string) bool {
	LogPrintInfo(prompt)
	ret, err := ReadUserInput("请确认 [Yy]: ")
	if err != nil {
		return false
	}
	trueAnswers := []string{"Y", "y"}
	if slices.Contains(trueAnswers, ret) {
		return true
	}
	return false
}

func ReadEnvVar(envVar string) (string, bool) {
	return os.LookupEnv(envVar)
}
