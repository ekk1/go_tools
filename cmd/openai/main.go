package main

import (
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"go_utils/utils/webui"
	"net/http"
)

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

	s := myhttp.NewServer("ai", "0.0.0.0", "7777")
	s.AddGet("/", handleIndex)
	s.AddPost("/chat", handleChat)
	s.Serve()
}

func renderPage() []byte {
	base := webui.NewNavBase("ChatGPT")
	base.AddNavItem("Index", "/")
	base.CurrentNavItem = "Index"

	base.AddSection(
		"",
		webui.NewCardHalf(
			webui.NewHeader("ChatGPT", "h1"),
			webui.NewLinkBtn("Refresh", "/"),
		),
	)

	base.AddSection("Input",
		webui.NewCardHalf(
			webui.NewForm("/chat", "Chat",
				webui.NewTextAreaInput("chat"),
				webui.NewSubmitBtn("Submit", "submit2"),
			),
		),
	)
	return []byte(base.Render())
}

func handleChat(w http.ResponseWriter, req *http.Request) {
	inputPrompt := req.FormValue("chat")
	utils.LogPrintInfo(inputPrompt)
	w.Write(renderPage())
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	w.Write(renderPage())
}
