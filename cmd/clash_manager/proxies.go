package main

import (
	"encoding/base64"
	"encoding/json"
	"go_utils/utils"
)

const (
	clashProxyKey = "CC_Proxy"
)

var (
	ClashProxies map[string]string = map[string]string{}
)

type Proxy struct {
	All  []string `json:"all"`
	Name string   `json:"name"`
	Type string   `json:"type"`
	Now  string   `json:"now"`
}

func LoadClashProxies() {
	if !kv.Exists(clashProxyKey) {
		return
	}
	pps := kv.Get(clashProxyKey)
	ppBytes, err := base64.URLEncoding.DecodeString(pps)
	utils.ErrExit(err)
	utils.ErrExit(json.Unmarshal(ppBytes, &ClashProxies))
}

func SaveClashProxies() {
	ppBytes, err := json.Marshal(ClashProxies)
	utils.ErrExit(err)
	ppStr := base64.URLEncoding.EncodeToString(ppBytes)
	utils.ErrExit(kv.Set(clashProxyKey, ppStr))
	utils.ErrExit(kv.Save())
}

func AddClashProxies(name, filter string) {
	ClashProxies[name] = filter
	SaveClashProxies()
}

func DeleteClashProxies(name string) {
	if _, ok := ClashProxies[name]; ok {
		delete(ClashProxies, name)
	}
	SaveClashProxies()
}
