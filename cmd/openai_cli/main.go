package main

import (
	"context"
	"flag"
	"go_utils/utils"
	"go_utils/utils/openai"
	"os"
	"strings"
)

func main() {
	translateFlag := flag.Bool("t", false, "translate mode")
	flag.Parse()

	if *translateFlag {
		utils.LogPrintDebug("Working in translate mode")
	}

	k, _ := os.ReadFile("key.txt")
	pp, _ := os.ReadFile("proxy.txt")
	e, _ := os.ReadFile("endpoint.txt")

	kStr := strings.TrimSuffix(string(k), "\n")
	ppStr := strings.TrimSuffix(string(pp), "\n")
	eStr := strings.TrimSuffix(string(e), "\n")

	cc := openai.NewOpenAIClient(eStr, kStr)
	if err := cc.SetProxy(ppStr); err != nil {
		utils.LogPrintError(err)
		os.Exit(1)
	}

	chatReq := openai.NewChatRequest("gpt-4o")
	chatReq.AddSystemMessage("")

	response := make(chan openai.ChatResponse)
	go cc.Chat(context.Background(), chatReq, response)
	for res := range response {
		if res.Error != "" {

		}
	}

}
