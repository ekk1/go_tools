package main

import (
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"go_utils/utils/webui"
	"net/http"
	"os"
	"slices"
)

var (
	pageMsg      string
	CurrentProxy string   = "p2"
	AllProxies   []string = []string{"p2", "p3"}
	CurrentNode  string
	AllNodes     []string
)

func onlineSaveSS(ss *Subscribe, req *http.Request, w http.ResponseWriter) bool {
	if err := ss.Save(); err != nil {
		utils.LogPrintError(err)
		myhttp.ServerError("failed to save db", w, req)
		return false
	}
	return true
}

func renderPage(w http.ResponseWriter, req *http.Request) {
	allNodes, nowNode, err := GetNodesByProxy(CurrentProxy)
	if err != nil {
		utils.LogPrintError(err)
		myhttp.ServerError("Failed get nodes", w, req)
		return
	}
	CurrentNode = nowNode
	AllNodes = allNodes

	base := webui.NewBase("cc")
	base.AddChild(webui.NewDiv(
		webui.NewHeader("Clash Manager", "h1"),
		webui.NewText("A simple manager"),
		webui.NewLinkBtn("Refresh", "/"),
		webui.NewLinkBtn("Generate", "/gen"),
	))
	infoTable := webui.NewTable(
		webui.NewTableRow(true, "Name", "Last", "", ""),
	)
	for _, s := range LoadSubscribe() {
		infoTable.AddChild(webui.NewTableRow(false,
			s.Name, s.LastUpdated,
			webui.NewLink("Update", "/update?name="+s.Name).Render(),
			webui.NewLink("Delete", "/delete?name="+s.Name).Render(),
		))
	}
	base.AddChild(webui.NewDiv(infoTable))

	base.AddChild(webui.NewDiv(webui.NewForm(
		"/", "Add sub",
		webui.NewTextInputWithValue("name", "ss"),
		webui.NewTextInputWithValue("url", "https://"),
		webui.NewSubmitBtn("Add", "submit"),
	)))
	base.AddChild(webui.NewDiv(webui.NewForm(
		"/rules", "Add rules",
		webui.NewTextInputWithValue("rule", "DOMAIN-SUFFIX,xxx.com,DIRECT"),
		webui.NewSubmitBtn("Add", "submit2"),
	)))

	rulesTable := webui.NewTable(
		webui.NewTableRow(true,
			"Rules", "",
		),
	)
	for _, p := range ClashRules {
		rulesTable.AddChild(webui.NewTableRow(false,
			p, webui.NewLink("Delete", "/deleterule?name="+p).Render(),
		))
	}

	proxyTable := webui.NewTable(
		webui.NewTableRow(true,
			"Proxy", "Now: "+CurrentProxy,
		),
	)
	for _, p := range AllProxies {
		proxyTable.AddChild(webui.NewTableRow(false,
			p, webui.NewLink("Select", "/selectproxy?name="+p).Render(),
		))
	}

	nodeTable := webui.NewTable(
		webui.NewTableRow(true,
			"Node", "Now: "+CurrentNode,
		),
	)
	for _, p := range AllNodes {
		nodeTable.AddChild(webui.NewTableRow(false,
			p, webui.NewLink("Select", "/selectnode?name="+p).Render(),
		))
	}

	base.AddChild(
		webui.NewDiv(rulesTable),
		webui.NewDiv(proxyTable),
		webui.NewDiv(nodeTable),
		webui.NewDiv(webui.NewText(pageMsg)),
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base.Render()))
}

func handleAddRules(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		myhttp.ServerError("method not post", w, req)
		return
	}
	ssRule := req.FormValue("rule")
	AddClashRules(ssRule)
	renderPage(w, req)
}

func handleDeleteRules(w http.ResponseWriter, req *http.Request) {
	ppName := req.URL.Query().Get("name")
	if slices.Contains(ClashRules, ppName) {
		DeleteClashRules(ppName)
	} else {
		myhttp.ServerError("No rule named: "+ppName, w, req)
		return
	}
	renderPage(w, req)
}

func handleSelectProxy(w http.ResponseWriter, req *http.Request) {
	ppName := req.URL.Query().Get("name")
	if slices.Contains(AllProxies, ppName) {
		CurrentProxy = ppName
	} else {
		myhttp.ServerError("No proxy named "+ppName, w, req)
		return
	}
	renderPage(w, req)
}

func handleSelectNode(w http.ResponseWriter, req *http.Request) {
	ppName := req.URL.Query().Get("name")
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
	renderPage(w, req)
}

func handleUpdete(w http.ResponseWriter, req *http.Request) {
	ssName := req.URL.Query().Get("name")
	utils.LogPrintInfo("Updating: ", ssName)

	ss := LoadSingleSubscribe(ssPrefix + ssName)
	if ss == nil {
		myhttp.ServerError("failed to get ss", w, req)
		return
	}
	if err := ss.Update(); err != nil {
		utils.LogPrintError(err)
		myhttp.ServerError("failed to update", w, req)
		return
	}
	if !onlineSaveSS(ss, req, w) {
		return
	}
	renderPage(w, req)
}

func handleDelete(w http.ResponseWriter, req *http.Request) {
	ssName := req.URL.Query().Get("name")
	utils.LogPrintInfo("Deleting: ", ssName)
	kv.Delete(ssPrefix + ssName)
	if err := kv.Save(); err != nil {
		utils.LogPrintError(err)
		myhttp.ServerError("failed to save db", w, req)
		return
	}

	renderPage(w, req)
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
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
		if !onlineSaveSS(ss, req, w) {
			return
		}
	}
	renderPage(w, req)
}

func handleGenerateYaml(w http.ResponseWriter, req *http.Request) {
	ret := RenderClashYaml(LoadSubscribe())
	err := os.WriteFile(clashYamlOutPath, []byte(ret), 0600)
	if err != nil {
		myhttp.ServerError("Failed to write yaml", w, req)
	}
	myhttp.ServerReply("Written", w)
}
