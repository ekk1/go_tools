package main

import (
	"encoding/base64"
	"encoding/json"
	"go_utils/utils"
	"slices"
)

const (
	clashRulesKey = "CC_Rules"
)

var (
	ClashRules []string = []string{}

	LastRuleType    = "DOMAIN-SUFFIX"
	LastRuleKeyword = "xxx.com"
	LastRuleTarget  = "DIRECT"
	LastRuleIndex   = "0"
)

type RuleGroup struct {
	Name   string   `json:"name"`
	Rules  []string `json:"rules"`
	Target string   `json:"target"`
}

func AddClashRules(index int, rule string) {
	ClashRules = slices.Insert(ClashRules, index, rule)
	SaveClashRules()
}

func DeleteClashRules(rule string) {
	i := slices.Index(ClashRules, rule)
	ClashRules = slices.Delete(ClashRules, i, i+1)
	SaveClashRules()
}

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
