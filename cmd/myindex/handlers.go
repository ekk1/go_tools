package main

import (
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"go_utils/utils/webui"
	"net/http"
	"slices"
	"sort"
	"strings"
)

var (
	pageMsg string
)

func renderPage(w http.ResponseWriter, req *http.Request) {
	base := webui.NewBase("myindex")
	refreshBtn := webui.NewLink("Refresh", "/")
	refreshBtn.SetClass("btn")
	headDiv := webui.NewDiv(
		webui.NewHeader("Index", "h1"),
		refreshBtn,
	)

	submitBtn := webui.NewInput("submit", "submit", "submit", "submit")
	submitBtn.SetClass("btn")
	submitDiv := webui.NewDiv(
		webui.NewForm("/", "Add",
			webui.NewTextInput("name"),
			webui.NewTextInputWithValue("url", "https://"),
			submitBtn,
		),
	)

	links := kv.Keys("LINK::")
	linkDict := map[string][]string{}
	dictKeys := []string{}
	for _, v := range links {
		fields := strings.Split(v, "::")
		if _, ok := linkDict[fields[1]]; !ok {
			linkDict[fields[1]] = []string{}
		}
		linkDict[fields[1]] = append(
			linkDict[fields[1]], fields[2],
		)
		if !slices.Contains(dictKeys, fields[1]) {
			dictKeys = append(dictKeys, fields[1])
		}
	}
	sort.Strings(dictKeys)
	infoDiv := webui.NewDiv()
	for _, v := range dictKeys {
		folderDiv := webui.NewDiv(
			webui.NewHeader(v, "h3"),
		)
		sort.Strings(linkDict[v])
		for _, name := range linkDict[v] {
			li := webui.NewLink(name, kv.Get("LINK::"+v+"::"+name))
			li.SetAttr("target", "_blank")
			li.SetClass("btn")
			folderDiv.AddChild(li, webui.NewBR())
		}
		infoDiv.AddChild(folderDiv)
	}

	base.AddChild(headDiv, submitDiv, infoDiv)

	w.Write([]byte(base.Render()))
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	myhttp.ServerLog("root", req)
	if !myhttp.ServerCheckPath("/", req, w) {
		return
	}

	if req.Method == http.MethodPost {
		linkName := req.FormValue("name")
		linkURL := req.FormValue("url")
		linkFolder := req.FormValue("folder")
		utils.LogPrintInfo("Got link:", linkName, linkURL, linkFolder)
		if !myhttp.ServerCheckParam(linkName, linkURL, linkFolder) {
			myhttp.ServerError("Field can not be empty", w, req)
			return
		}
		kKey := "LINK::" + linkFolder + "::" + linkName
		kv.Set(kKey, linkURL)
		if err := kv.Save(); err != nil {
			myhttp.ServerError("Failed to save DB", w, req)
			return
		}
	}
	renderPage(w, req)
}
