package main

import (
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
)

func handleIndex(w http.ResponseWriter, req *http.Request) {
	ppName := req.URL.Query().Get("name")

	switch req.URL.Query().Get("action") {
	case "start":
		if !RunClash() {
			myhttp.ServerError("maybe already started", w, req)
			return
		}
	case "stop":
		if !StopClash() {
			myhttp.ServerError("maybe not running", w, req)
			return
		}
	case "generate":
		ret := RenderClashYaml()
		err := os.WriteFile(clashYamlOutPath, []byte(ret), 0600)
		if err != nil {
			myhttp.ServerError("Failed to write yaml", w, req)
			return
		}
		myhttp.ServerReply("Written", w)
		return
	case "proxy":
		if _, ok := ClashProxies[ppName]; ok {
			CurrentProxy = ppName
		} else {
			myhttp.ServerError("No proxy named "+ppName, w, req)
			return
		}
	case "node":
		if slices.Contains(AllNodes, ppName) {
			utils.LogPrintInfo("Selecting", ppName, "for", CurrentProxy)
			err := ChangeNodeForProxy(CurrentProxy, ppName)
			if err != nil {
				utils.LogPrintError(err)
				myhttp.ServerError("Failed to select node", w, req)
				return
			}
		} else {
			myhttp.ServerError("No proxy named "+ppName, w, req)
			return
		}
	}
	if CurrentProxy != "" {
		allNodes, nowNode, err := GetNodesByProxy(CurrentProxy)
		if err != nil {
			utils.LogPrintError(err)
			CurrentProxy = ""
		}
		CurrentNode = nowNode
		AllNodes = allNodes
	}
	renderIndex(w)
}

func handleRules(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if strings.Contains(req.URL.Path, "delete") {
			name := req.URL.Query().Get("name")
			if slices.Contains(ClashRules, name) {
				DeleteClashRules(name)
			} else {
				myhttp.ServerError("No rule named: "+name, w, req)
				return
			}
		}
	case http.MethodPost:
		ssRuleType := req.FormValue("rule_type")
		ssRuleKeyword := req.FormValue("rule_keyword")
		ssRuleTarget := req.FormValue("rule_target")
		ssRuleIndex := req.FormValue("index")
		if !myhttp.ServerCheckParam(ssRuleType, ssRuleIndex,
			ssRuleTarget, ssRuleIndex,
		) {
			myhttp.ServerError("field can not be empty", w, req)
			return
		}
		LastRuleType = ssRuleType
		LastRuleKeyword = ssRuleKeyword
		LastRuleTarget = ssRuleTarget
		LastRuleIndex = ssRuleIndex
		ssRuleInt, err := strconv.Atoi(ssRuleIndex)
		if err != nil {
			myhttp.ServerError("Index should be int", w, req)
			return
		}
		trueIndex := 0
		switch {
		case ssRuleInt == -1:
			trueIndex = len(ClashRules)
		case ssRuleInt < len(ClashRules):
			trueIndex = ssRuleInt
		case ssRuleInt >= len(ClashRules):
			trueIndex = len(ClashRules)
		}
		AddClashRules(trueIndex, strings.Join(
			[]string{ssRuleType, ssRuleKeyword, ssRuleTarget},
			",",
		))
	}
	renderRulesManager(w)
}

func handleProxies(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if strings.Contains(req.URL.Path, "delete") {
			name := req.URL.Query().Get("name")
			if _, ok := ClashProxies[name]; ok {
				DeleteClashProxies(name)
			} else {
				myhttp.ServerError("No proxy named: "+name, w, req)
				return
			}
		}
	case http.MethodPost:
		ppName := req.FormValue("name")
		ppFilter := req.FormValue("filter")
		AddClashProxies(ppName, ppFilter)
	}
	renderProxyManager(w)
}

func handleSubs(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	SelectedExtraNode = ""
	switch req.Method {
	case http.MethodGet:
		switch {
		case strings.Contains(req.URL.Path, "delete"):
			utils.LogPrintInfo("Deleting: ", name)
			kv.Delete(ssPrefix + name)
			if err := kv.Save(); err != nil {
				utils.LogPrintError(err)
				myhttp.ServerError("failed to save db", w, req)
				return
			}
		case strings.Contains(req.URL.Path, "update"):
			utils.LogPrintInfo("Updating: ", name)
			ss := LoadSingleSubscribe(ssPrefix + name)
			proxy := req.URL.Query().Get("proxy")
			if ss == nil {
				myhttp.ServerError("failed to get ss", w, req)
				return
			}
			useProxy := false
			if proxy != "" {
				useProxy = true
			}
			if err := ss.Update(useProxy); err != nil {
				utils.LogPrintError(err)
				myhttp.ServerError("failed to update", w, req)
				return
			}
			if err := ss.Save(); err != nil {
				utils.LogPrintError(err)
				myhttp.ServerError("failed to save", w, req)
				return
			}
		}
	case http.MethodPost:
		ssName := req.FormValue("name")
		ssURL := req.FormValue("url")
		utils.LogPrintInfo("Got sub:", ssName)
		if !myhttp.ServerCheckParam(ssName, ssURL) {
			myhttp.ServerError("Field can not be empty", w, req)
			return
		}
		ss := &Subscribe{
			Name: ssName,
			URL:  ssURL,
		}
		if err := ss.Save(); err != nil {
			utils.LogPrintError(err)
			myhttp.ServerError("failed to save", w, req)
			return
		}
	}
	renderSubManager(w)
}

func handleExtraNodes(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	SelectedExtraNode = ""
	switch req.Method {
	case http.MethodGet:
		switch {
		case strings.Contains(req.URL.Path, "delete"):
			DeleteCustomNode(name)
		case strings.Contains(req.URL.Path, "edit"):
			SelectedExtraNode = name
		}
	case http.MethodPost:
		trName := req.FormValue("name")
		trAddr := req.FormValue("address")
		trPort := req.FormValue("port")
		trPass := req.FormValue("password")
		trSni := req.FormValue("sni")
		utils.LogPrintInfo("Got extra:", trName)
		if !myhttp.ServerCheckParam(trName, trAddr, trPort, trPass, trSni) {
			myhttp.ServerError("Field can not be empty", w, req)
			return
		}
		tr := &Trojan{
			Name:     trName,
			Address:  trAddr,
			Port:     trPort,
			Password: trPass,
			Sni:      trSni,
		}
		if err := AddCustomNode(trName, tr); err != nil {
			utils.LogPrintError(err)
			myhttp.ServerError("failed to add custom node", w, req)
			return
		}
	}
	renderSubManager(w)
}
