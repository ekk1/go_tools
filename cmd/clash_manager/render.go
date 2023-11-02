package main

import (
	"go_utils/utils/webui"
	"net/http"
	"slices"
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
	ppList := []string{}
	for p := range ClashProxies {
		ppList = append(ppList, p)
	}
	slices.Sort(ppList)
	for _, p := range ppList {
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

	infoTable := webui.NewTable(webui.NewTableRow(true, "Name", "Last", "", ""))
	for _, s := range LoadSubscribe() {
		infoTable.AddChild(webui.NewTableRow(false,
			s.Name, s.LastUpdated,
			webui.NewLink("Update", "/subs/update?name="+s.Name).Render(),
			webui.NewLink("Delete", "/subs/delete?name="+s.Name).Render(),
		))
	}
	bb.AddChild(webui.NewColumnDiv(
		webui.NewDiv3C(webui.NewForm(
			"/subs", "Add sub",
			webui.NewTextInputWithValue("name", "ss"),
			webui.NewTextInputWithValue("url", "https://"),
			webui.NewSubmitBtn("Add", "submit"),
		)),
		webui.NewDiv6C(infoTable),
	))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(bb.Render()))
}

func renderProxyManager(w http.ResponseWriter) {
	bb := generatePageStruct()
	proxyTable := webui.NewTable(webui.NewTableRow(true, "Name", "Filter", ""))
	ppList := []string{}
	for p := range ClashProxies {
		ppList = append(ppList, p)
	}
	slices.Sort(ppList)
	for _, p := range ppList {
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
