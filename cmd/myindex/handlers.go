package main

import (
	"errors"
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"go_utils/utils/webui"
	"net/http"
	"slices"
	"strings"
	"time"
)

func prepareLinksData() map[string][]string {
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

	return linkDict
}

func renderPage(w http.ResponseWriter, _ *http.Request) {
	base := webui.NewNavBase("myindex")

	base.AddNavItem("Index", "/")
	base.CurrentNavItem = "Index"

	base.AddSection(
		"",
		webui.NewCardHalf(
			webui.NewHeader("Index", "h1"),
			webui.NewLinkBtn("Refresh", "/"),
		),
	)
	base.AddSection("Edit",
		webui.NewCardThird(
			webui.NewForm("/add", "Add & Update",
				webui.NewTextInput("name"),
				webui.NewTextInputWithValue("url", "https://"),
				webui.NewTextInputWithValue("folder", "default"),
				webui.NewSubmitBtn("submit", "submit1"),
			),
		),
		webui.NewCardThird(
			webui.NewForm("/delete", "Delete",
				webui.NewTextInput("name"),
				webui.NewTextInputWithValue("folder", "default"),
				webui.NewSubmitBtn("delete", "submit2"),
			),
		),
		webui.NewCardThird(
			webui.NewForm("/move", "Move",
				webui.NewTextInput("name"),
				webui.NewTextInput("name_new"),
				webui.NewTextInputWithValue("old_folder", "default"),
				webui.NewTextInputWithValue("new_folder", "default"),
				webui.NewSubmitBtn("move", "submit3"),
			),
		),
	)

	linkDict := prepareLinksData()
	for _, v := range utils.SortedMapKeys(linkDict) {
		folderDiv := webui.NewCardRest(
			webui.NewHeader(v, "h3"),
		)
		slices.Sort(linkDict[v])
		for _, name := range linkDict[v] {
			li := webui.NewLinkBtn(name, kv.Get("LINK::"+v+"::"+name))
			li.SetAttr("target", "_blank")
			folderDiv.AddChild(li)
		}
		base.AddSection("", folderDiv)
	}

	w.Write([]byte(base.Render()))
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	authCookie, err := req.Cookie("authToken")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Redirect(w, req, "/loginpage", http.StatusSeeOther)
			return
		default:
			utils.LogPrintError(err)
			myhttp.ServerError("server error", w, req)
			return
		}
	}
	if ts, ok := cookieCache[authCookie.Value]; !ok {
		http.Redirect(w, req, "/loginpage", http.StatusSeeOther)
	} else {
		timeCookie := time.Unix(ts, 0)
		timePassed := time.Now().Sub(timeCookie)
		if timePassed.Hours() > 1 {
			http.Redirect(w, req, "/loginpage", http.StatusSeeOther)
		}
	}
	renderPage(w, req)
}

func handleLoginPage(w http.ResponseWriter, _ *http.Request) {
	login := webui.NewLoginPage("/login", "MyIndex")
	w.Write([]byte(login.Render()))
}

func handleLogin(w http.ResponseWriter, _ *http.Request) {
	login := webui.NewLoginPage("/login", "MyIndex")
	w.Write([]byte(login.Render()))
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
	linkNewName := req.FormValue("name_new")
	linkOldFolder := req.FormValue("old_folder")
	linkNewFolder := req.FormValue("new_folder")
	utils.LogPrintInfo("Move link:", linkName, linkNewName, linkOldFolder, linkNewFolder)
	if !myhttp.ServerCheckParam(linkName, linkNewName, linkOldFolder, linkNewFolder) {
		myhttp.ServerError("Field can not be empty", w, req)
		return
	}
	kKey := "LINK::" + linkOldFolder + "::" + linkName
	if !kv.Exists(kKey) {
		myhttp.ServerError("Not exists", w, req)
		return
	}
	kKeyNew := "LINK::" + linkNewFolder + "::" + linkNewName
	kv.Set(kKeyNew, kv.Get(kKey))
	kv.Delete(kKey)
	if err := kv.Save(); err != nil {
		myhttp.ServerError("Failed to save DB", w, req)
		return
	}
	renderPage(w, req)
}
