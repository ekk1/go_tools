package openai

import (
	"bufio"
	"context"
	"fmt"
	"go_utils/utils"
	"os"
	"testing"
)

func TestOpenAI(t *testing.T) {
	k, _ := os.ReadFile("key.txt")
	pp, _ := os.ReadFile("proxy.txt")
	e, _ := os.ReadFile("endpoint.txt")
	c := NewOpenAIClient(
		string(e),
		string(k),
	)
	if err := c.SetProxy(string(pp)); err != nil {
		t.Fatal(err)
	}

	// utils.LogPrintInfo(c.ListModels())

	responses := make(chan ChatResponse)
	req := NewChatRequest("gpt-4o")
	req.AddSystemMessage("You are a experienced vim pro, and you can write vimscript and other things without using any third party plugins, just with pure vim")
	req.AddUserMessage("help me write a lsp client in vimscript for gopls, make it asynchronus, and use this as the omnifunc for golang, also make the completion list show as i input codes, and make the completion list more beautiful, like surrounded by borders like a balloon, and no background color, making it look like transparent")

	go c.Chat(context.Background(), req, responses)

	writer := bufio.NewWriter(os.Stdout)

	for res := range responses {
		if res.Error != "" {
			fmt.Fprintln(writer, "Error:", res.Error)
			writer.Flush()
		} else {
			if len(res.Choices) > 0 {
				fmt.Fprint(writer, res.Choices[0].Delta.Content)
				writer.Flush()
			} else {
				utils.LogPrintWarning("Empty response:?", res)
			}
		}
	}
}
