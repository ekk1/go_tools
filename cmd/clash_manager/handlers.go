package main

import (
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"net/http"
	"slices"
)

var (
	pageMsg      string
	CurrentProxy string   = "p2"
	AllProxies   []string = []string{"p2", "p3"}
	CurrentNode  string
	AllNodes     []string
)

type PageData struct {
	Title        string
	Msg          string
	Subscribe    []*Subscribe
	CurrentProxy string
	AllProxies   []string
	CurrentNode  string
	AllNodes     []string
}

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
	data := PageData{
		Title:        "test",
		Msg:          pageMsg,
		Subscribe:    LoadSubscribe(),
		CurrentProxy: CurrentProxy,
		AllProxies:   AllProxies,
		CurrentNode:  nowNode,
		AllNodes:     allNodes,
	}
	if err := indexTemplate.ExecuteTemplate(w, "index", data); err != nil {
		utils.LogPrintError(err)
		myhttp.ServerError("Failed to exec template", w, req)
	}
}

func handleSelectProxy(w http.ResponseWriter, req *http.Request) {
	myhttp.ServerLog("selectProxy", req)
	if !myhttp.ServerCheckPath("/selectproxy", req, w) {
		return
	}
	ppName := req.URL.Query().Get("name")
	if slices.Contains(AllProxies, ppName) {
		CurrentProxy = ppName
	} else {
		myhttp.ServerError("No proxy named "+ppName, w, req)
	}
	renderPage(w, req)
}

func handleSelectNode(w http.ResponseWriter, req *http.Request) {
	myhttp.ServerLog("selectNodes", req)
	if !myhttp.ServerCheckPath("/selectnode", req, w) {
		return
	}
	ppName := req.URL.Query().Get("name")
	if slices.Contains(AllNodes, ppName) {
		utils.LogPrintInfo("Selecting", ppName, "for", CurrentProxy)
		err := ChangeNodeForProxy(CurrentProxy, ppName)
		if err != nil {
			utils.LogPrintError(err)
			myhttp.ServerError("Failed to select node", w, req)
		}
	} else {
		myhttp.ServerError("No proxy named "+ppName, w, req)
	}
	renderPage(w, req)
}

func handleUpdete(w http.ResponseWriter, req *http.Request) {
	myhttp.ServerLog("update", req)
	if !myhttp.ServerCheckPath("/update", req, w) {
		return
	}

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
	myhttp.ServerLog("delete", req)
	if !myhttp.ServerCheckPath("/delete", req, w) {
		return
	}

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
	myhttp.ServerLog("root", req)
	if !myhttp.ServerCheckPath("/", req, w) {
		return
	}

	if req.Method == http.MethodPost {
		ssName := req.FormValue("name")
		ssAction := req.FormValue("action")
		ssURL := req.FormValue("url")
		utils.LogPrintInfo("Got sub:", ssAction, ssName)
		if !myhttp.ServerCheckParam(ssName, ssAction, ssURL) {
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
