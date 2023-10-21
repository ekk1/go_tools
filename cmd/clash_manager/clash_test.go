package main

import (
	"go_utils/utils"
	"go_utils/utils/minikv"
	"testing"
)

func TestClash(t *testing.T) {
	kv = minikv.MustNewKV("ss", 0)
	kv.MustLoad()

	LoadClashRules()

	subs := LoadSubscribe()
	utils.LogPrintInfo(RenderClashYaml(subs))

}
