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

func AddClashRules(index int, rule string) {
	ClashRules = slices.Insert(ClashRules, index, rule)
	SaveClashRules()
}

func DeleteClashRules(rule string) {
	i := slices.Index(ClashRules, rule)
	ClashRules = slices.Delete(ClashRules, i, i+1)
	SaveClashRules()
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

func RenderClashYaml(subs []*Subscribe) string {
	nodes := ""
	nodesNames := ""
	proxies := ""
	allNodeList := []string{}
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
					allNodeList = append(allNodeList, tt.Name)
				}
			}
		case "flower":
			for _, line := range utils.SplitByLine(sub.Content) {
				if strings.Contains(line, "trojan") {
					tt, err := NewFlowerTrojan(line)
					utils.ErrExit(err)
					nodes += tt.Render()
					nodesNames += "      - \"" + tt.Name + "\"\n"
					allNodeList = append(allNodeList, tt.Name)
				}
			}
		}
	}
	ppList := []string{}
	for p := range ClashProxies {
		ppList = append(ppList, p)
	}
	slices.Sort(ppList)
	for _, ppName := range ppList {
		proxies += "  - name: " + ppName + "\n"
		proxies += "    type: select\n"
		proxies += "    proxies:\n"
		filterKeys := strings.Split(ClashProxies[ppName], ",")
		for _, nn := range allNodeList {
			for _, kk := range filterKeys {
				if strings.Contains(nn, kk) {
					proxies += "      - \"" + nn + "\"\n"
				}
			}
		}
		proxies += "\n"
	}
	rules := ""
	for _, r := range ClashRules {
		rules += "  - " + r + "\n"
	}
	rules += "  - MATCH,all"
	return fmt.Sprintf(clashYamlBase, nodes, proxies, rules)
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
