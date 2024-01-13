package utils

import (
	"os"
	"strings"
	"time"
)

func ListDirFiles(path string) ([]string, error) {
	f, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for _, v := range f {
		if !v.Type().IsDir() {
			ret = append(ret, v.Name())
		}
	}
	return ret, nil
}

func GetLatestFileInDir(path string) (string, error) {
	f, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	ret := ""
	retTime := time.Unix(0, 0)
	for _, v := range f {
		if !v.Type().IsDir() {
			_i, err := v.Info()
			if err != nil {
				return "", err
			}
			if _i.ModTime().After(retTime) {
				ret = _i.Name()
				retTime = _i.ModTime()
			}
		}
	}
	return ret, nil
}

func ListDirDirs(path string) ([]string, error) {
	f, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for _, v := range f {
		if v.Type().IsDir() {
			ret = append(ret, v.Name())
		}
	}
	return ret, nil
}

func ListDirItems(path string) []string {
	ret := []string{}
	f, err := os.ReadDir(path)
	if err != nil {
		return ret
	}
	LogPrintDebug("ReadDir: ", f)
	for _, v := range f {
		ret = append(ret, v.Name())
	}
	return ret
}

func FindLineInContent(content string, keywords []string, once bool) []string {
	matchedLines := []string{}
	for _, line := range strings.Split(content, "\n") {
		matched := true
		for _, kw := range keywords {
			if !strings.Contains(line, kw) {
				matched = false
			}
		}
		if matched {
			matchedLines = append(matchedLines, line)
		}
	}
	return matchedLines
}
