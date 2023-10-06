package main

import (
	_ "embed"
	"go_utils/utils"
	"net/http"
)

type PageData struct {
	Title     string
	Msg       string
	Lister    []string
	Mapper    map[string]string
	Subscribe []*Subscribe
}

func PreparePageData() *PageData {
	return &PageData{
		Title:     "test",
		Msg:       "msg",
		Subscribe: LoadSubscribe(),
	}
}

func handleUpdete(w http.ResponseWriter, req *http.Request) {
	utils.ServerLog("update", req)
	ssName := req.URL.Query().Get("name")
	utils.LogPrintInfo("Updating: ", ssName)

	ss := LoadSingleSubscribe(ssPrefix + ssName)
	if ss == nil {
		utils.ServerError("failed to get ss", w, req)
		return
	}
	if err := ss.Update(); err != nil {
		utils.LogPrintError(err)
		utils.ServerError("failed to update", w, req)
		return
	}
	if err := ss.Save(); err != nil {
		utils.LogPrintError(err)
		utils.ServerError("failed to save db", w, req)
		return
	}

	if err := indexTemplate.ExecuteTemplate(w, "index", PreparePageData()); err != nil {
		utils.LogPrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to execute template"))
	}
}

func handleDelete(w http.ResponseWriter, req *http.Request) {
	utils.ServerLog("delete", req)
	ssName := req.URL.Query().Get("name")
	utils.LogPrintInfo("Deleting: ", ssName)
	kv.Delete(ssPrefix + ssName)
	if err := kv.Save(); err != nil {
		utils.LogPrintError(err)
		utils.ServerError("failed to save db", w, req)
		return
	}

	if err := indexTemplate.ExecuteTemplate(w, "index", PreparePageData()); err != nil {
		utils.LogPrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to execute template"))
	}
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	utils.ServerLog("root", req)

	// Deny req not to /
	if req.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(""))
	}

	// Print headers if debug
	for k, v := range req.Header {
		utils.LogPrintDebug(k, v)
	}
	if req.Method == http.MethodPost {
		ssName := req.FormValue("name")
		ssAction := req.FormValue("action")
		ssURL := req.FormValue("url")
		utils.LogPrintInfo(ssAction, ssName, ssURL)
		if ssAction == "" {
			utils.ServerError("Failed to decide ssAction", w, req)
			return
		}
		if ssName == "" || ssURL == "" {
			utils.ServerError("Field can not be empty", w, req)
			return
		}
		ss := &Subscribe{
			Name: ssName,
			URL:  ssURL,
		}
		if err := ss.Save(); err != nil {
			utils.ServerError("Failed to save ss", w, req)
			utils.LogPrintError(err)
			return
		}
	}

	if err := indexTemplate.ExecuteTemplate(w, "index", PreparePageData()); err != nil {
		utils.LogPrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to execute template"))
	}
}
