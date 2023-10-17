package utils

import (
	"math/rand"
	"os"
)

func ErrExit(err error) {
	if err != nil {
		LogPrintError(err)
		os.Exit(1)
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
