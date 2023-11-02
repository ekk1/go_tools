package utils

import (
	"bytes"
	"math/rand"
	"slices"
	"strings"
)

func ErrExit(err error) {
	if err != nil {
		LogPrintError(err)
		panic(err)
	}
}

func SplitByLine(content string) []string {
	if bytes.Contains([]byte(content), []byte("\r\n")) {
		LogPrintDebug("Spliting content with CRLF")
		return strings.Split(content, "\r\n")
	}
	return strings.Split(content, "\n")
}

func RandomString(length int64) string {
	ret := make([]rune, length)
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	for i := range ret {
		ret[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(ret)
}

func SortedKeysStr(m map[string]any) []string {
	ret := make([]string, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	slices.Sort(ret)
	return ret
}

func SortedKeysInt(m map[int]any) []int {
	ret := make([]int, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	slices.Sort(ret)
	return ret
}
