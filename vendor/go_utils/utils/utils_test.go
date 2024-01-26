package utils

import (
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	for _, cmd := range []string{"sleep 1", "date", "aa"} {
		if ret, err := RunCmd(cmd, nil); err == nil {
			LogPrintInfo(ret)
		} else {
			LogPrintError(ret)
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

func TestSome(t *testing.T) {
	aa := "aaa   0/1     Completed   0          1d"
	LogPrintInfo(strings.Fields(aa))
}

func TestSortMapKeys(t *testing.T) {
	a := map[string]string{
		"11":         "33",
		"111":        "33",
		"34":         "33",
		"122":        "33",
		"99":         "33",
		"91912xxadw": "33",
		"xadawda":    "33",
		"1axzq2":     "33",
		"a9a8987":    "33",
	}
	for _, v := range SortedMapKeys(a) {
		t.Log(v)
	}
	b := map[int]string{
		01:       "33",
		22:       "33",
		99:       "33",
		19012893: "33",
		91:       "33",
		99319:    "33",
	}
	for _, v := range SortedMapKeys(b) {
		t.Log(v)
	}
}

func TestGenericWaiter(t *testing.T) {
	tester := "test"

	ff := func() bool {
		if tester == "test" {
			return true
		} else {
			t.Log("tester is:", tester)
		}
		return false
	}

	if err := GenericWaiter(10, ff, "aa"); err != nil {
		t.Log("Failed wait")
		t.Log(err)
	} else {
		t.Log("Finished")
	}
}
