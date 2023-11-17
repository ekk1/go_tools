package main

import (
	"encoding/base64"
	"encoding/json"
	"go_utils/utils"
)

var (
	ClashRules []string = []string{}
)

func LoadClashRules() {
	if !kv.Exists(clashRulesKey) {
		return
	}
	rules := kv.Get(clashRulesKey)
	rulesBytes, err := base64.URLEncoding.DecodeString(rules)
	utils.ErrExit(err)
	utils.ErrExit(json.Unmarshal(rulesBytes, &ClashRules))
}

func SaveClashRules() {
	rulesBytes, err := json.Marshal(ClashRules)
	utils.ErrExit(err)
	rulesStr := base64.URLEncoding.EncodeToString(rulesBytes)
	utils.ErrExit(kv.Set(clashRulesKey, rulesStr))
	utils.ErrExit(kv.Save())
}
