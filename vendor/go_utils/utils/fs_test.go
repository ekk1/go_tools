package utils

import (
	"testing"
)

func TestGetLatestFileInDir(t *testing.T) {
	LogPrintInfo(GetLatestFileInDir("."))
}
