package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"go_utils/utils/openai"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var QUICK_PROMPT = map[string]string{
	"ff": "帮我翻译一下我发过来的内容到中文",
}

var aiClient *openai.OpenAIClient

func main() {
	k, _ := os.ReadFile("key.txt")
	pp, _ := os.ReadFile("proxy.txt")
	e, _ := os.ReadFile("endpoint.txt")
	aiClient = openai.NewOpenAIClient(
		string(e),
		string(k),
	)
	if err := aiClient.SetProxy(string(pp)); err != nil {
		utils.LogPrintError(err)
		os.Exit(1)
	}

	s := myhttp.NewServer("ai", "127.0.0.1", "7777")
	s.AddGet("/", handleIndex)
	s.AddGet("/script.js", handleJS)
	s.AddGet("/w3.css", handleCSS)
	s.AddGet("/favicon.ico", faviconHandler)
	s.AddPost("/chat", handleChat)
	s.AddPost("/save", handleSave)
	s.AddPost("/load", handleLoad)
	s.AddGet("/list", handleList)
	s.Serve()
}

func handleChat(w http.ResponseWriter, req *http.Request) {
	utils.LogPrintInfo(fmt.Sprintf("Request from %s: %s %s", req.RemoteAddr, req.Method, req.URL.Path))
	// req.ParseForm()
	// data, _ := io.ReadAll(req.Body)
	// utils.LogPrintInfo(string(data))
	var (
		inputPrompt     = req.FormValue("prompt")
		inputStatements = req.FormValue("chat")
		model           = req.FormValue("model")
	)
	// utils.LogPrintInfo(inputPrompt)
	// utils.LogPrintInfo(inputStatements)
	// utils.LogPrintInfo(model)

	if len(inputStatements) == 0 || len(model) == 0 {
		w.Write([]byte("No content\n"))
		return
	}

	chatReq := openai.NewChatRequest(model)
	if len(inputPrompt) != 0 {
		utils.LogPrintInfo("Adding system: ", inputPrompt)
		chatReq.AddSystemMessage(inputPrompt)
	}

	states := strings.Split(inputStatements, "<<>>++--==--++<<>>")
	for _, v := range states {
		// utils.LogPrintInfo("Adding user: ", v)
		chatReq.AddUserMessage(v)
	}

	flush := func() {
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	responses := make(chan openai.ChatResponse)
	go aiClient.Chat(context.Background(), chatReq, responses)

	for res := range responses {
		if res.Error != "" {
			w.Write([]byte("\n Error during processing: " + res.Error + "\n"))
			flush()
		} else {
			if len(res.Choices) > 0 {
				// utils.LogPrintInfo(res.Choices[0].Delta.Content)
				w.Write([]byte(res.Choices[0].Delta.Content))
				flush()
			} else {
				utils.LogPrintWarning("Empty response:?", res)
			}
		}
	}
	utils.LogPrintInfo("DONE")
	w.Write([]byte("\n[[[DONE]]]\n"))
	flush()
}

func handleSave(w http.ResponseWriter, req *http.Request) {
	utils.LogPrintInfo(fmt.Sprintf("Request from %s: %s %s", req.RemoteAddr, req.Method, req.URL.Path))
	var (
		inputStatements = req.FormValue("chat")
		saveName        = req.FormValue("name")
	)
	if match, err := regexp.MatchString(`[a-zA-Z0-9]+`, saveName); err != nil {
		utils.LogPrintError(err)
		w.Write([]byte("Failed\n"))
		return
	} else {
		if !match {
			w.Write([]byte("Illegal save name\n"))
			return
		}
	}
	if len(inputStatements) == 0 {
		w.Write([]byte("No content\n"))
		return
	}
	saveName = "save_" + saveName
	if err := os.WriteFile(saveName, []byte(inputStatements), 0600); err != nil {
		w.Write([]byte("Failed write\n"))
		return
	}

	w.Write([]byte("OK"))
}

func handleLoad(w http.ResponseWriter, req *http.Request) {
	utils.LogPrintInfo(fmt.Sprintf("Request from %s: %s %s", req.RemoteAddr, req.Method, req.URL.Path))
	var (
		saveName = req.FormValue("name")
	)
	if match, err := regexp.MatchString(`[a-zA-Z0-9]+`, saveName); err != nil {
		utils.LogPrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed\n"))
		return
	} else {
		if !match {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Illegal save name\n"))
			return
		}
	}
	data, err := os.ReadFile(saveName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed read data\n"))
		return
	}
	w.Write(data)
}

func handleList(w http.ResponseWriter, req *http.Request) {
	utils.LogPrintInfo(fmt.Sprintf("Request from %s: %s %s", req.RemoteAddr, req.Method, req.URL.Path))
	files, err := os.ReadDir(".")
	if err != nil {
		utils.LogPrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to list files\n"))
		return
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "save_") { // 过滤前缀
			fileList = append(fileList, file.Name())
		}
	}

	fileListJSON, err := json.Marshal(fileList)
	if err != nil {
		utils.LogPrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode file list\n"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(fileListJSON)
}
