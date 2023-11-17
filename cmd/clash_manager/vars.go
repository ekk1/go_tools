package main

import (
	"go_utils/utils/minikv"
	"sync"
)

const (
	clashManagerAddr = "http://127.0.0.1:9090"
	clashRulesKey    = "CC_Rules"
	clashProxyKey    = "CC_Proxy"
)

var (
	kv               *minikv.KV
	clashYamlOutPath string
	clashBinary      string
	clashConfigDir   string

	pageMsg      string
	CurrentProxy string
	CurrentNode  string
	AllNodes     []string

	ClashProxies map[string]string = map[string]string{}

	LastRuleType    = "DOMAIN-SUFFIX"
	LastRuleKeyword = "xxx.com"
	LastRuleTarget  = "DIRECT"
	LastRuleIndex   = "0"

	procChan chan int = make(chan int)
	procSync          = &struct {
		state bool
		lock  sync.Mutex
	}{
		state: false,
	}
)

type Proxy struct {
	All  []string `json:"all"`
	Name string   `json:"name"`
	Type string   `json:"type"`
	Now  string   `json:"now"`
}

type RuleGroup struct {
	Name   string   `json:"name"`
	Rules  []string `json:"rules"`
	Target string   `json:"target"`
}
