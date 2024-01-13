package main

import (
	"encoding/base64"
	"fmt"
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"net/http"
	"strings"
)

func RenderClashYaml() string {
	nodes := ""
	proxies := ""
	allNodeList := []string{}
	for k, v := range CustomNodes {
		nodes += v.Render()
		allNodeList = append(allNodeList, k)
	}
	for _, sub := range LoadSubscribe() {
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
					allNodeList = append(allNodeList, tt.Name)
				}
			}
		}
	}

	for _, ppName := range utils.SortedMapKeys(ClashProxies) {
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

const (
	clashYamlBase = `
---
port: 7890
socks-port: 10099
redir-port: 7892
allow-lan: true
mode: rule
log-level: info
ipv6: false
external-controller: 127.0.0.1:9090
dns:
  enable: true
  listen: 0.0.0.0:53
  default-nameserver:
    - 114.114.114.114
  enhanced-mode: fake-ip # or fake-ip
  fake-ip-range: 198.18.0.1/16 # Fake IP addresses pool CIDR
  nameserver:
    - 114.114.114.114 # default value

proxies:
  # socks5
  - name: "socks"
    type: socks5
    server: 127.0.0.1
    port: 10099
    udp: true
%s

proxy-groups:
%s

rules:
%s
`
)
