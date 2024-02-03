package main

import (
	"go_utils/utils"
	"testing"
)

func TestXianyu(t *testing.T) {
	testStr := ""
	tr, err := NewXianyuTrojan(testStr)
	utils.LogPrintInfo(tr, err)

}
