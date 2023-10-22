package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"net/http"
	"slices"
	"strings"
)

const (
	clashManagerAddr = "http://127.0.0.1:9090"
	clashRulesKey    = "CC_Rules"
)

type Proxy struct {
	All  []string `json:"all"`
	Name string   `json:"name"`
	Type string   `json:"type"`
	Now  string   `json:"now"`
}

var ClashRules []string = []string{}

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

func AddClashRules(rule string) {
	ClashRules = append(ClashRules, rule)
	SaveClashRules()
}

func DeleteClashRules(rule string) {
	i := slices.Index(ClashRules, rule)
	ClashRules = slices.Delete(ClashRules, i, i+1)
	SaveClashRules()
}

func RenderClashYaml(subs []*Subscribe) string {
	nodes := ""
	nodesNames := ""
	for _, sub := range subs {
		switch sub.Name {
		case "xianyu":
			data, err := base64.StdEncoding.DecodeString(sub.Content)
			utils.ErrExit(err)
			for _, line := range utils.SplitByLine(string(data)) {
				if strings.Contains(line, "trojan") {
					tt, err := NewXianyuTrojan(line)
					utils.ErrExit(err)
					nodes += tt.Render()
					nodesNames += "      - \"" + tt.Name + "\"\n"
				}
			}
		case "flower":
			for _, line := range utils.SplitByLine(sub.Content) {
				if strings.Contains(line, "trojan") {
					tt, err := NewFlowerTrojan(line)
					utils.ErrExit(err)
					nodes += tt.Render()
					nodesNames += "      - \"" + tt.Name + "\"\n"
				}
			}
		}
	}
	rules := ""
	for _, r := range ClashRules {
		rules += "  - " + r + "\n"
	}
	rules += "  - MATCH,p2"
	return fmt.Sprintf(clashYamlBase, nodes, nodesNames, nodesNames, rules)
}

func GetNodesByProxy(p string) ([]string, string, error) {
	c := myhttp.NewHTTPClient()
	ret, err := c.SendReq(http.MethodGet, clashManagerAddr+"/proxies/"+p, nil)
	if err != nil {
		return nil, "", err
	}
	utils.LogPrintDebug(ret.Text())
	utils.LogPrintDebug(ret)
	var pp *Proxy = &Proxy{}
	if err := ret.JSON(pp); err != nil {
		return nil, "", err
	}

	return pp.All, pp.Now, nil
}

func ChangeNodeForProxy(p, n string) error {
	c := myhttp.NewHTTPClient()
	c.SetSendJSON(true)
	sendDict := map[string]string{"name": n}
	ret, err := c.SendReq(http.MethodPut, clashManagerAddr+"/proxies/"+p, sendDict)
	if err != nil {
		return err
	}
	utils.LogPrintInfo(ret.Text())
	return nil
}
