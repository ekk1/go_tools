package main

import (
	"go_utils/utils"
	"go_utils/utils/webui"
	"net/http"
	"strconv"
)

func generatePageStruct() *webui.Base {
	stateStr := "Clash core is not running"
	procSync.lock.Lock()
	defer procSync.lock.Unlock()
	if procSync.state {
		stateStr = "Clash core is up and running"
	}

	base := webui.NewBase("cc")
	base.AddChild(
		webui.NewDiv(
			webui.NewHeader("Clash Manager", "h1"),
			webui.NewText("A simple manager"),
			webui.NewLinkBtn("Refresh", "/"),
			webui.NewLinkBtn("Generate", "/?action=generate"),
			webui.NewLinkBtn("Start", "/?action=start"),
			webui.NewLinkBtn("Stop", "/?action=stop"),
			webui.NewText(stateStr),
		),
		webui.NewDiv(
			webui.NewLinkBtn("Index", "/"),
			webui.NewLinkBtn("Subs", "/subs"),
			webui.NewLinkBtn("Proxies", "/proxies"),
			webui.NewLinkBtn("Rules", "/rules"),
		),
	)
	return base
}

func renderIndex(w http.ResponseWriter) {
	base := generatePageStruct()

	proxyTable := webui.NewTable(
		webui.NewTableRow(true, "Proxy", "Now: "+CurrentProxy),
	)
	for _, p := range utils.SortedMapKeys(ClashProxies) {
		proxyTable.AddChild(webui.NewTableRow(false,
			p, webui.NewLink("Select", "/?action=proxy&name="+p).Render(),
		))
	}

	nodeTable := webui.NewTable(
		webui.NewTableRow(true,
			"Node", "Now: "+CurrentNode,
		),
	)
	for _, p := range AllNodes {
		nodeTable.AddChild(webui.NewTableRow(false,
			p, webui.NewLink("Select", "/?action=node&name="+p).Render(),
		))
	}

	base.AddChild(webui.NewColumnDiv(
		webui.NewDiv4C(proxyTable),
		webui.NewDiv4C(nodeTable),
	))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base.Render()))
}

func renderSubManager(w http.ResponseWriter) {
	bb := generatePageStruct()

	infoTable := webui.NewTable(webui.NewTableRow(true, "Name", "Last", "", "", ""))
	extraTable := webui.NewTable(webui.NewTableRow(true, "Name", "", ""))
	for _, s := range LoadSubscribe() {
		infoTable.AddChild(webui.NewTableRow(false, s.Name, s.LastUpdated,
			webui.NewLink("Update", "/subs/update?name="+s.Name).Render(),
			webui.NewLink("UpdateWithProxy", "/subs/update?proxy=1&name="+s.Name).Render(),
			webui.NewLink("Delete", "/subs/delete?name="+s.Name).Render(),
		))
	}
	for k := range CustomNodes {
		extraTable.AddChild(webui.NewTableRow(false, k,
			webui.NewLink("Edit", "/extra/edit?name="+k).Render(),
			webui.NewLink("Delete", "/extra/delete?name="+k).Render(),
		))
	}
	var (
		extraNodeName     = "tr1"
		extraNodeAddress  = ""
		extraNodePort     = ""
		extraNodePassword = ""
		extraNodeSni      = ""
	)
	if SelectedExtraNode != "" {
		if t, ok := CustomNodes[SelectedExtraNode]; ok {
			extraNodeName = t.Name
			extraNodeAddress = t.Address
			extraNodePort = t.Port
			extraNodePassword = t.Password
			extraNodeSni = t.Sni
		}
	}
	bb.AddChild(webui.NewColumnDiv(
		webui.NewDiv3C(
			webui.NewDiv(webui.NewForm(
				"/subs", "Add sub",
				webui.NewTextInputWithValue("name", "ss"),
				webui.NewTextInputWithValue("url", "https://"),
				webui.NewSubmitBtn("Add", "submit"),
			)),
			webui.NewDiv(webui.NewForm(
				"/extra", "Add Extra Trojan",
				webui.NewTextInputWithValue("name", extraNodeName),
				webui.NewTextInputWithValue("address", extraNodeAddress),
				webui.NewTextInputWithValue("port", extraNodePort),
				webui.NewTextInputWithValue("password", extraNodePassword),
				webui.NewTextInputWithValue("sni", extraNodeSni),
				webui.NewSubmitBtn("Add", "submit"),
			)),
		),
		webui.NewDiv6C(webui.NewDiv(infoTable), webui.NewDiv(extraTable)),
	))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(bb.Render()))
}

func renderProxyManager(w http.ResponseWriter) {
	bb := generatePageStruct()
	proxyTable := webui.NewTable(webui.NewTableRow(true, "Name", "Filter", ""))
	for _, p := range utils.SortedMapKeys(ClashProxies) {
		proxyTable.AddChild(webui.NewTableRow(false, p, ClashProxies[p],
			webui.NewLink("Delete", "/proxies/delete?name="+p).Render(),
		))
	}
	bb.AddChild(webui.NewColumnDiv(
		webui.NewDiv3C(webui.NewForm(
			"/proxies", "Add proxy",
			webui.NewTextInputWithValue("name", "pp"),
			webui.NewTextInputWithValue("filter", "JP"),
			webui.NewSubmitBtn("Add", "submit"),
		)),
		webui.NewDiv6C(proxyTable),
	))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(bb.Render()))
}

func renderRulesManager(w http.ResponseWriter) {
	bb := generatePageStruct()
	rulesTable := webui.NewTable(
		webui.NewTableRow(true,
			"Index", "Rules", "",
		),
	)
	for no, p := range ClashRules {
		rulesTable.AddChild(webui.NewTableRow(false, strconv.Itoa(no),
			p, webui.NewLink("Delete", "/rules/delete?name="+p).Render(),
		))
	}
	bb.AddChild(webui.NewColumnDiv(
		webui.NewDiv3C(webui.NewForm(
			"/rules", "Add rules",
			webui.NewTextInputWithValue("rule_type", LastRuleType),
			webui.NewTextInputWithValue("rule_keyword", LastRuleKeyword),
			webui.NewTextInputWithValue("rule_target", LastRuleTarget),
			webui.NewTextInputWithValue("index", LastRuleIndex),
			webui.NewSubmitBtn("Add", "submit2"),
		)),
		webui.NewDiv6C(rulesTable),
	))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(bb.Render()))
}
