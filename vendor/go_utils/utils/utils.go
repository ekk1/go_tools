package utils

import (
	"os"
)

func ErrExit(err error) {
	if err != nil {
		LogPrintError(err)
		os.Exit(1)
	}
}
