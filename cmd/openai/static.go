package main

import (
	_ "embed"
	"encoding/base64"
	"net/http"
)

//go:embed index.html
var indexHTML []byte

//go:embed script.js
var js []byte

//go:embed w3.css
var css []byte

func handleIndex(w http.ResponseWriter, req *http.Request) {
	w.Write(indexHTML)
}
func handleJS(w http.ResponseWriter, req *http.Request) {
	// utils.LogPrintInfo(req.Header)
	w.Write(js)
}
func handleCSS(w http.ResponseWriter, req *http.Request) {
	// utils.LogPrintInfo(req.Header)
	w.Write(css)
}

func faviconHandler(w http.ResponseWriter, req *http.Request) {
	// 空白 PNG 图片的 base64 编码
	base64Image := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/wcAAgUBAAD/1POZAAAAAElFTkSuQmCC"

	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		http.Error(w, "Failed to decode image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(imageData)
}
