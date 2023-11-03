package utils

import (
	"bytes"
	"cmp"
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

func SortedMapKeys[K cmp.Ordered, V any](m map[K]V) []K {
	ret := []K{}
	for k := range m {
		ret = append(ret, k)
	}
	slices.Sort(ret)
	return ret
}
