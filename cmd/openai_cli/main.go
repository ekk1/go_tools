package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"go_utils/utils"
	"go_utils/utils/openai"
	"io"
	"os"
	"slices"
	"strings"
	"unicode/utf8"
)

const (
	SEPRATOR = "<<>>++__--!!@@##--<<>>\n"
)

var (
	QUICK_PROMPT map[string]string = map[string]string{
		"ff": "你是一个翻译专家，帮我翻译一下我发过来的内容到中文",
		"cc": "你是一个编程专家，擅长在尽量不使用外部依赖的前提下，使用尽可能原生的方法完成各类编程任务",
	}
)

func main() {
	verboseFlag := flag.Int("v", 0, "debug (max 4)")
	translateFlag := flag.Bool("t", false, "translate mode, only read from stdin")
	addFlag := flag.Bool("a", false, "add mode, add file content to context")
	nameFlag := flag.String("n", "default", "file name to add")
	flag.Parse()
	utils.SetLogLevelByVerboseFlag(*verboseFlag)

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

	if *translateFlag {
		inputBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			utils.LogPrintError("Failed to read from stdin: ", err)
			os.Exit(1)
		}
		if !utf8.Valid(inputBytes) {
			utils.LogPrintError("Illegal utf8 input detected")
			os.Exit(1)
		}
		data := string(inputBytes)
		chatReq := openai.NewChatRequest("gpt-4o-mini")
		chatReq.AddSystemMessage("你是一个翻译专家，帮我翻译一下我发过来的内容到中文，并将中文和原文对照的结果返回给我")
		chatReq.AddUserMessage(data)

		responseCh := make(chan openai.ChatResponse)
		go cc.Chat(context.Background(), chatReq, responseCh)
		for res := range responseCh {
			if res.Error != "" {
				fmt.Print("Error during request: " + res.Error)
			} else {
				if len(res.Choices) > 0 {
					fmt.Print(res.Choices[0].Delta.Content)
				} else {
					fmt.Print("No choice returned")
				}
			}
		}
	}

	if *addFlag {
		_, err := os.Stat("io.txt")
		if err != nil {
			utils.LogPrintError("Failed to load io file", err)
			os.Exit(1)
		}

		inputBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			utils.LogPrintError("Failed to read from stdin: ", err)
			os.Exit(1)
		}
		if !utf8.Valid(inputBytes) {
			utils.LogPrintError("Illegal utf8 input detected")
			os.Exit(1)
		}
		data := string(inputBytes)

		chatReq := openai.NewChatRequest("gpt-4o")
		if err := decodeIOFile(chatReq); err != nil {
			utils.LogPrintError("Failed to decode io file", err)
			os.Exit(1)
		}
		if len(chatReq.Messages) == 0 {
			utils.LogPrintError("Illegal io file")
			os.Exit(1)
		}

		if strings.HasSuffix(data, "\n") {
			data = *nameFlag + "\n```\n" + data + "```\n"
		} else {
			data = *nameFlag + "\n```\n" + data + "\n```\n"
		}
		if chatReq.Messages[len(chatReq.Messages)-1].Role != "user" {
			if err := writeIOFile("user", "\n"+data); err != nil {
				utils.LogPrintError("Failed to write file data", err)
				os.Exit(1)
			}
		} else {
			chatReq.Messages[len(chatReq.Messages)-1].Content += "\n" + data
			err := os.WriteFile("io.txt", []byte(""), 0600)
			if err != nil {
				utils.LogPrintError("Failed to clear io file", err)
				os.Exit(1)
			}
			for _, m := range chatReq.Messages {
				if err := writeIOFile(m.Role, m.Content); err != nil {
					utils.LogPrintError("Failed to write io file", err)
					os.Exit(1)
				}
			}
		}

	}

	if !*translateFlag && !*addFlag {
		_, err := os.Stat("io.txt")
		if err != nil {
			if os.IsNotExist(err) {
				if errInit := initIOFile(); errInit != nil {
					utils.LogPrintError("Failed to init io file: ", errInit)
					os.Exit(1)
				}
			} else {
				utils.LogPrintError("Failed to load io file", err)
				os.Exit(1)
			}
		}
		chatReq := openai.NewChatRequest("gpt-4o")
		if err := decodeIOFile(chatReq); err != nil {
			utils.LogPrintError("Failed to decode io file", err)
			os.Exit(1)
		}
		utils.LogPrintInfo("History Context:")
		for _, m := range chatReq.Messages {
			if len(m.Content) > 20 {
				utils.LogPrintInfo(m.Role, ": ", m.Content[0:19]+"[TRUNCATED]")
			} else {
				utils.LogPrintInfo(m.Role, ": ", m.Content)
			}
		}
		utils.LogPrintInfo("cc to clear context, sa to save dialog, ll to list dialog, lo to load dialog")
		for {
			if len(chatReq.Messages) == 0 {
				if errInit := initIOFile(); errInit != nil {
					utils.LogPrintError("Failed to init io file: ", errInit)
					os.Exit(1)
				}
				if err := decodeIOFile(chatReq); err != nil {
					utils.LogPrintError("Failed to decode io file", err)
					os.Exit(1)
				}
			}
			if chatReq.Messages[len(chatReq.Messages)-1].Role != "user" {
				ccFlag := false
				for {
					userInput, err := readUserInputMulti("[两次回车] User")
					if err != nil {
						utils.LogPrintError("Failed to read user input", err)
						os.Exit(1)
					}
					if len(userInput) == 0 {
						fmt.Println("User input cannot be empty")
					} else {
						if userInput == "cc" {
							chatReq.Messages = nil
							if errInit := initIOFile(); errInit != nil {
								utils.LogPrintError("Failed to init io file: ", errInit)
								os.Exit(1)
							}
							if err := decodeIOFile(chatReq); err != nil {
								utils.LogPrintError("Failed to decode io file", err)
								os.Exit(1)
							}
							utils.LogPrintInfo("Current context: ")
							for _, m := range chatReq.Messages {
								utils.LogPrintInfo(m.Role, m.Content)
							}
							ccFlag = true
							break
						}
						if strings.HasPrefix(userInput, "sa ") {
							saveName := strings.TrimPrefix(userInput, "sa ")
							saveName = strings.ReplaceAll(saveName, " ", "_")
							if slices.Contains([]string{"io", "key", "proxy", "endpoint"}, saveName) {
								utils.LogPrintError("Illegal save name")
								ccFlag = true
								break
							}
							data, err := os.ReadFile("io.txt")
							if err != nil {
								utils.LogPrintError("Failed to read io file", err)
								os.Exit(1)
							}
							err = os.WriteFile(saveName+".txt", data, 0600)
							if err != nil {
								utils.LogPrintError("Failed to read io file", err)
								os.Exit(1)
							}
							ccFlag = true
							break
						}
						if strings.HasPrefix(userInput, "lo ") {
							saveName := strings.TrimPrefix(userInput, "lo ")
							saveName = strings.ReplaceAll(saveName, " ", "_")
							if slices.Contains([]string{"io", "key", "proxy", "endpoint"}, saveName) {
								utils.LogPrintError("Illegal save name")
								ccFlag = true
								break
							}
							data, err := os.ReadFile(saveName + ".txt")
							if err != nil {
								utils.LogPrintError("Failed to read io file", err)
								os.Exit(1)
							}
							err = os.WriteFile("io.txt", data, 0600)
							if err != nil {
								utils.LogPrintError("Failed to read io file", err)
								os.Exit(1)
							}
							chatReq.Messages = nil
							if err := decodeIOFile(chatReq); err != nil {
								utils.LogPrintError("Failed to decode io file", err)
								os.Exit(1)
							}
							utils.LogPrintInfo("Current context: ")
							for _, m := range chatReq.Messages {
								utils.LogPrintInfo(m.Role, m.Content)
							}
							ccFlag = true
							break
						}
						if userInput == "ll" {
							ioFiles := []string{}
							files, err := os.ReadDir(".")
							if err != nil {
								utils.LogPrintError("Failed to list io files", err)
								os.Exit(1)
							}
							for _, f := range files {
								if f.IsDir() {
									continue
								}
								if strings.HasSuffix(f.Name(), ".txt") {
									if !slices.Contains([]string{"io.txt", "key.txt", "proxy.txt", "endpoint.txt"}, f.Name()) {
										ioFiles = append(ioFiles, f.Name())
									}
								}
							}
							if len(ioFiles) == 0 {
								fmt.Println("No saved dialogs")
							}
							for _, io := range ioFiles {
								fmt.Println(io)
							}
							ccFlag = true
							break
						}
						if err := writeIOFile("user", userInput); err != nil {
							utils.LogPrintError("Failed to write user input", err)
							os.Exit(1)
						}
						chatReq.AddUserMessage(userInput)
						break
					}
				}
				if ccFlag {
					continue
				}
			}
			response := make(chan openai.ChatResponse)
			utils.LogPrintDebug("Sending request")
			resultBuffer := ""
			go cc.Chat(context.Background(), chatReq, response)
			fmt.Print("AI: ")
			for res := range response {
				utils.LogPrintDebug4("res: ", res)
				if res.Error != "" {
					fmt.Print("Error during request: " + res.Error)
				} else {
					if len(res.Choices) > 0 {
						fmt.Print(res.Choices[0].Delta.Content)
						resultBuffer += res.Choices[0].Delta.Content
					} else {
						fmt.Print("No choice returned")
					}
				}
			}
			fmt.Print("\n")
			if err := writeIOFile("ai", resultBuffer); err != nil {
				utils.LogPrintError("Failed to write ai result", err)
				os.Exit(1)
			}
			chatReq.AddAssistantMessage(resultBuffer)
		}
	}
}

func initIOFile() error {
	f, err := os.OpenFile("io.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	fmt.Println("Quick prompts:")
	for k, v := range QUICK_PROMPT {
		fmt.Println(k + ": " + v)
	}
	var userPrompt string = ""
	userInput, err := readUserInputMulti("Enter prompt")
	if err != nil {
		return err
	}
	if quickMatch, ok := QUICK_PROMPT[userInput]; ok {
		userPrompt = quickMatch
	} else {
		userPrompt = userInput
	}
	_, err = f.Write([]byte("PROMPT: " + userPrompt + "\n"))
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(SEPRATOR))
	if err != nil {
		return err
	}
	return nil
}

func writeIOFile(role, content string) error {
	f, err := os.OpenFile("io.txt", os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	switch role {
	case "user":
		_, err = f.Write([]byte("USER: " + content + "\n"))
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(SEPRATOR))
		if err != nil {
			return err
		}
	case "ai":
		_, err = f.Write([]byte("AI: " + content + "\n"))
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(SEPRATOR))
		if err != nil {
			return err
		}
	case "assistant":
		_, err = f.Write([]byte("AI: " + content + "\n"))
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(SEPRATOR))
		if err != nil {
			return err
		}
	case "system":
		_, err = f.Write([]byte("PROMPT: " + content + "\n"))
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(SEPRATOR))
		if err != nil {
			return err
		}
	default:
		return errors.New("Undefined role")
	}

	return nil
}

func decodeIOFile(ai *openai.ChatRequest) error {
	data, err := os.ReadFile("io.txt")
	if err != nil {
		return err
	}

	if !utf8.Valid(data) {
		return errors.New("Invalid io file")
	}

	parts := strings.Split(string(data), SEPRATOR)
	for _, p := range parts {
		if strings.HasPrefix(p, "PROMPT: ") {
			promptString := strings.TrimPrefix(p, "PROMPT: ")
			ai.AddSystemMessage(strings.TrimSuffix(promptString, "\n"))
		}
		if strings.HasPrefix(p, "USER: ") {
			userString := strings.TrimPrefix(p, "USER: ")
			ai.AddUserMessage(strings.TrimSuffix(userString, "\n"))
		}
		if strings.HasPrefix(p, "AI: ") {
			aiString := strings.TrimPrefix(p, "AI: ")
			ai.AddAssistantMessage(strings.TrimSuffix(aiString, "\n"))
		}
	}

	utils.LogPrintDebug("Decoded IO file: ")
	for _, m := range ai.Messages {
		utils.LogPrintDebug(m.Role, m.Content)
	}

	return nil
}

func readUserInputMulti(prompt string) (string, error) {
	fmt.Print(prompt + ": ")
	reader := bufio.NewReader(os.Stdin)
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		if len(strings.TrimSpace(line)) == 0 {
			break
		}
		lines = append(lines, strings.TrimSuffix(line, "\n"))
	}
	return strings.Join(lines, "\n"), nil
}
