package utils

import (
	"math/rand"
	"os"
)

func ErrExit(err error, rc int) {
	if err != nil {
		LogPrintError(err)
		os.Exit(rc)
	}
}

func RandomString(length int64) string {
	ret := make([]rune, length)
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	for i := range ret {
		ret[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(ret)
}
