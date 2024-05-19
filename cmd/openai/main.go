package main

import (
	"fmt"
	"go_utils/utils"
)

const SEPRATOR = "<<>>++__--!!@@##--<<>>\n"

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
	fmt.Println("Choose Model: ")
	for n, m := range MODELS {
		fmt.Printf("[%d]: %s", n, m)
	}
	modelNoStr, err := utils.ReadUserInput("Enter: ")
	if err != nil {
		utils.LogPrintError(err)
	}

}
