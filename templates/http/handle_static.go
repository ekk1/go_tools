package main

import (
	_ "embed"
	"go_utils/utils/quickserver"
	"net/http"
)

//go:embed static/index.html
var indexHTML string

//go:embed static/js/base.js
var baseJSFile string

// and so on

var assetFileList map[string]string

func prepareAssetDict() {
	assetFileList = make(map[string]string)
	assetFileList["static/js/base.js"] = baseJSFile
}

func handleStatic(w http.ResponseWriter, req *http.Request) {
	quickserver.QuickServerLog(req, "handleStatic")

	dataBytes, ok := assetFileList[req.URL.Path[1:]]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(""))
		return
	}
	w.Write([]byte(dataBytes))
}
