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

func TestSome(t *testing.T) {
	aa := map[string]int64{}
	t.Log(aa)
	aa["test"]++
	aa["test"]++
	aa["test"]++
	aa["test"]++
	t.Log(aa)
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
