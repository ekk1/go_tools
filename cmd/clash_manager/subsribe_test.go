package main

import (
	"encoding/base64"
	"go_utils/utils"
	"go_utils/utils/minikv"
	"strings"
	"testing"
)

func TestModels(t *testing.T) {
	utils.SetLogLevelDebug()
	kv = minikv.MustNewKV("ss", 0)
	kv.MustLoad()
	s := LoadSubscribe()
	for _, v := range s {
		t.Log(v.Name)
		t.Log(v.LastUpdated)
		switch v.Name {
		case "flower":
			for _, line := range utils.SplitByLine(v.Content) {
				if strings.Contains(line, "trojan") {
					tt, err := NewFlowerTrojan(line)
					utils.ErrExit(err)
					utils.LogPrintInfo(tt.Render())
				}
			}
		case "xianyu":
			data, err := base64.StdEncoding.DecodeString(v.Content)
			utils.ErrExit(err)
			for _, line := range utils.SplitByLine(string(data)) {
				if strings.Contains(line, "trojan") {
					tt, err := NewXianyuTrojan(line)
					utils.ErrExit(err)
					utils.LogPrintInfo(tt.Render())
				}
			}

		default:

		}
	}
}
