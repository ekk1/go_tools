package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatUsage struct {
	PromptToken     int64 `json:"prompt_tokens"`
	CompletionToken int64 `json:"completion_tokens"`
	TotalToken      int64 `json:"total_tokens"`
}

type ChatChoice struct {
	Index          int64        `json:"index"`
	Messages       *ChatMessage `json:"messages"`
	Delta          *ChatMessage `json:"delta"`
	FinishedReason string       `json:"finish_reason"`
}

type ModelData struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type ModelResponse struct {
	Object string       `json:"object"`
	Data   []*ModelData `json:"data"`
}

type ChatRequest struct {
	Model    string         `json:"model"`
	Messages []*ChatMessage `json:"messages"`
	Stream   bool           `json:"stream"`
}

func NewChatRequest(model string) *ChatRequest {
	return &ChatRequest{
		Model:    model,
		Messages: []*ChatMessage{},
	}
}

func (c *ChatRequest) AddSystemMessage(msg string) {
	c.Messages = append(c.Messages, &ChatMessage{
		Role:    "system",
		Content: msg,
	})
}

func (c *ChatRequest) AddUserMessage(msg string) {
	c.Messages = append(c.Messages, &ChatMessage{
		Role:    "user",
		Content: msg,
	})
}

func (c *ChatRequest) AddAssistantMessage(msg string) {
	c.Messages = append(c.Messages, &ChatMessage{
		Role:    "assistant",
		Content: msg,
	})
}

type ChatResponse struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Model   string `json:"model"`

	Choices []*ChatChoice `json:"choices"`
	Error   string        `json:"error,omitempty"`
}

type OpenAIClient struct {
	endpoint string
	key      string
	client   *http.Client
}

func NewOpenAIClient(endpoint, key string) *OpenAIClient {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	ts := &http.Transport{
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c := &http.Client{
		Transport: ts,
	}

	return &OpenAIClient{
		endpoint: endpoint,
		key:      key,
		client:   c,
	}
}

func (c *OpenAIClient) SetProxy(proxyURL string) error {
	pURL, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}
	if ts, ok := c.client.Transport.(*http.Transport); ok {
		ts.Proxy = http.ProxyURL(pURL)
	} else {
		return errors.New("failed to get transport")
	}
	return nil
}

func (c *OpenAIClient) Chat(ctx context.Context, request *ChatRequest, responseChan chan<- ChatResponse) {
	defer close(responseChan)

	request.Stream = true
	data, err := json.Marshal(request)
	if err != nil {
		responseChan <- ChatResponse{Error: err.Error()}
		return
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost,
		c.endpoint+"/v1/chat/completions",
		bytes.NewBuffer(data),
	)
	if err != nil {
		responseChan <- ChatResponse{Error: err.Error()}
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.key)

	resp, err := c.client.Do(req)
	if err != nil {
		responseChan <- ChatResponse{Error: err.Error()}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseChan <- ChatResponse{
			Error: "non-200 status code received" + strconv.Itoa(
				resp.StatusCode,
			),
		}
		return
	}

	// decoder := json.NewDecoder(resp.Body)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			responseChan <- ChatResponse{Error: "request cancelled"}
			return
		default:
		}

		line := scanner.Bytes()
		if len(line) < 2 {
			continue
		}
		if !bytes.Contains(line, []byte("data:")) {
			continue
		}
		if bytes.Contains(line, []byte("[DONE]")) && !bytes.Contains(line, []byte("model")) {
			continue
		}
		// utils.LogPrintInfo(string(line[5:]))

		var response ChatResponse
		if err := json.Unmarshal(line[5:], &response); err != nil {
			responseChan <- ChatResponse{Error: err.Error()}
			return
		}

		responseChan <- response
	}
	if err := scanner.Err(); err != nil {
		responseChan <- ChatResponse{Error: err.Error()}
	}
}

func (c *OpenAIClient) ListModels() ([]string, error) {
	ret := []string{}

	req, err := http.NewRequest(
		http.MethodGet,
		c.endpoint+"/v1/models",
		nil,
	)
	if err != nil {
		return ret, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.key)

	resp, err := c.client.Do(req)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	// utils.LogPrintInfo(data)

	modelResponse := &ModelResponse{}
	if err := json.Unmarshal(data, modelResponse); err != nil {
		return ret, err
	}

	for _, v := range modelResponse.Data {
		ret = append(ret, v.ID)
	}

	return ret, nil
}
