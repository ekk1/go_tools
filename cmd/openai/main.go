package main

import (
	_ "embed"
	"fmt"
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"net/http"
	"time"
)

//go:embed index.html
var indexHTML []byte

//go:embed script.js
var js []byte

var QUICK_PROMPT = map[string]string{
	"ff": "帮我翻译一下我发过来的内容到中文",
}

var MODELS = []string{
	"gpt-4o",
	"gpt-3.5-turbo-1106",
	"dall-e-3",
}

var KEYS = []*APIKey{}
var UsingKey *APIKey
var UsingModel string

func main() {
	LoadAllKeys()

	s := myhttp.NewServer("ai", "127.0.0.1", "7777")
	s.AddGet("/", handleIndex)
	s.AddGet("/script.js", handleJS)
	s.AddPost("/chat", handleChat)
	s.Serve()

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handleIndex)
	// mux.HandleFunc()
	// mux.HandleFunc("/chat", handleChat)
	// utils.LogPrintError(http.ListenAndServe("127.0.0.1:7777", mux))
}

func handleChat(w http.ResponseWriter, req *http.Request) {
	// inputPrompt := req.FormValue("chat")
	// utils.LogPrintInfo(inputPrompt)
	// w.Write([]byte("OK"))
	utils.LogPrintInfo(fmt.Sprintf("Request from %s: %s %s", req.RemoteAddr, req.Method, req.URL.Path))
	flush := func() {
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
	for ii := 0; ii < 10; ii++ {
		w.Write([]byte("test\n"))
		flush()
		time.Sleep(100 * time.Millisecond)
	}
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	w.Write(indexHTML)
}
func handleJS(w http.ResponseWriter, req *http.Request) {
	utils.LogPrintInfo(req.Header)
	w.Write(js)
}
