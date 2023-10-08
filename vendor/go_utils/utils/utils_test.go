package utils

import (
	"testing"
)

func TestRunCmd(t *testing.T) {
	for _, cmd := range []string{"sleep 1", "date"} {
		if ret, err := RunCmd(cmd, nil); err == nil {
			LogPrintInfo(ret)
		} else {
			LogPrintError(err)
		}
	}
}

func TestListDirFiles(t *testing.T) {
	if f, err := ListDirFiles("."); err == nil {
		LogPrintInfo(f)
	} else {
		LogPrintError(err)
	}
}
