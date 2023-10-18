package main

import (
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"go_utils/utils/webui"
	"net/http"
	"sort"
	"strings"
)

var (
	pageMsg string
)

func renderPage(w http.ResponseWriter, req *http.Request) {
	base := webui.NewBase("myindex")
	headDiv := webui.NewDiv(
		webui.NewHeader("Index", "h1"),
		webui.NewButton("Refresh", "/"),
	)

	submitDiv := webui.NewDiv(
		webui.NewForm("/", "Add",
			webui.NewTextInput("name"),
			webui.NewTextInputWithValue("url", "https://"),
			webui.NewTextInputWithValue("folder", "default"),
			webui.NewInput("submit", "submit", "submit", "submit"),
		),
	)

	links := kv.Keys("LINK::")
	linkDict := map[string][]string{}
	for _, v := range links {
		fields := strings.Split(v, "::")
		if _, ok := linkDict[fields[1]]; !ok {
			linkDict[fields[1]] = []string{}
		}
		linkDict[fields[1]] = append(
			linkDict[fields[1]], fields[2],
		)
	}
	infoDiv := webui.NewDiv()
	for k, v := range linkDict {
		folderDiv := webui.NewDiv(
			webui.NewHeader(k, "h3"),
		)
		sort.Strings(v)
		for _, name := range v {
			li := webui.NewLink(name, kv.Get("LINK::"+k+"::"+name))
			li.SetAttr("target", "_blank")
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
