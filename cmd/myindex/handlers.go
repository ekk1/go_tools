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

func prepareLinksData() ([]string, map[string][]string) {
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
	return dictKeys, linkDict
}

func renderPage(w http.ResponseWriter, req *http.Request) {
	base := webui.NewBase("myindex")

	headDiv := webui.NewDiv(
		webui.NewHeader("Index", "h1"),
		webui.NewLinkBtn("Refresh", "/"),
	)

	submitDiv := webui.NewDiv(
		webui.NewForm("/add", "Add & Update",
			webui.NewTextInput("name"),
			webui.NewTextInputWithValue("url", "https://"),
			webui.NewTextInputWithValue("folder", "default"),
			webui.NewSubmitBtn("submit", "submit1"),
		),
	)

	deleteDiv := webui.NewDiv(
		webui.NewForm("/delete", "Delete",
			webui.NewTextInput("name"),
			webui.NewTextInputWithValue("folder", "default"),
			webui.NewSubmitBtn("delete", "submit2"),
		),
	)

	moveDiv := webui.NewDiv(
		webui.NewForm("/move", "Move",
			webui.NewTextInput("name"),
			webui.NewTextInputWithValue("old_folder", "default"),
			webui.NewTextInputWithValue("new_folder", "default"),
			webui.NewSubmitBtn("move", "submit3"),
		),
	)

	dictKeys, linkDict := prepareLinksData()
	infoDiv := webui.NewDiv()
	for _, v := range dictKeys {
		folderDiv := webui.NewDiv(
			webui.NewHeader(v, "h3"),
		)
		sort.Strings(linkDict[v])
		for _, name := range linkDict[v] {
			li := webui.NewLinkBtn(name, kv.Get("LINK::"+v+"::"+name))
			li.SetAttr("target", "_blank")
			folderDiv.AddChild(li)
		}
		infoDiv.AddChild(folderDiv)
	}

	formDiv := webui.NewColumnDiv(submitDiv, deleteDiv, moveDiv, webui.NewDiv())
	base.AddChild(headDiv, formDiv, infoDiv)

	w.Write([]byte(base.Render()))
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	renderPage(w, req)
}

func handleAdd(w http.ResponseWriter, req *http.Request) {
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
	renderPage(w, req)
}

func handleDelete(w http.ResponseWriter, req *http.Request) {
	linkName := req.FormValue("name")
	linkFolder := req.FormValue("folder")
	utils.LogPrintInfo("Delete link:", linkName, linkFolder)
	if !myhttp.ServerCheckParam(linkName, linkFolder) {
		myhttp.ServerError("Field can not be empty", w, req)
		return
	}
	kKey := "LINK::" + linkFolder + "::" + linkName
	kv.Delete(kKey)
	if err := kv.Save(); err != nil {
		myhttp.ServerError("Failed to save DB", w, req)
		return
	}
	renderPage(w, req)
}

func handleMove(w http.ResponseWriter, req *http.Request) {
	linkName := req.FormValue("name")
	linkOldFolder := req.FormValue("old_folder")
	linkNewFolder := req.FormValue("new_folder")
	utils.LogPrintInfo("Move link:", linkName, linkOldFolder, linkNewFolder)
	if !myhttp.ServerCheckParam(linkName, linkOldFolder, linkNewFolder) {
		myhttp.ServerError("Field can not be empty", w, req)
		return
	}
	kKey := "LINK::" + linkOldFolder + "::" + linkName
	kKeyNew := "LINK::" + linkNewFolder + "::" + linkName
	kv.Set(kKeyNew, kv.Get(kKey))
	kv.Delete(kKey)
	if err := kv.Save(); err != nil {
		myhttp.ServerError("Failed to save DB", w, req)
		return
	}
	renderPage(w, req)
}
