package utils

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"
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

// Wait <timeout> seconds for <f> to return true, return <warning> when failed
// <f> is called every 3 secs
func GenericWaiter(timeout int64, f func() bool, warning string) error {
	startTime := time.Now()
	for {
		nowTime := time.Now()
		if nowTime.Sub(startTime) > (time.Duration(timeout) * time.Second) {
			msg := fmt.Sprintf(
				"Failed to %s within %d seconds", warning, timeout,
			)
			return errors.New(msg)
		}
		if f() {
			break
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}
