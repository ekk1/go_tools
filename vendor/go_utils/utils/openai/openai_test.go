package openai

import (
	"bufio"
	"context"
	"fmt"
	"go_utils/utils"
	"os"
	"strings"
	"testing"
)

func TestOpenAI(t *testing.T) {
	k, _ := os.ReadFile("key.txt")
	pp, _ := os.ReadFile("proxy.txt")
	e, _ := os.ReadFile("endpoint.txt")
	c := NewOpenAIClient(
		strings.TrimSuffix(string(e), "\n"),
		strings.TrimSuffix(string(k), "\n"),
	)
	if err := c.SetProxy(strings.TrimSuffix(string(pp), "\n")); err != nil {
		t.Fatal(err)
	}

	// utils.LogPrintInfo(c.ListModels())

	responses := make(chan ChatResponse)
	req := NewChatRequest("gpt-4o")
	req.AddUserMessage("how to disable all login from tty console, and only allow login from ssh")

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
