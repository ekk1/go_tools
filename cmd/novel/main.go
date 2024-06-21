package main

import (
	"flag"
	"go_utils/utils"
	"os"
	"regexp"
	"strings"
)

func main() {
	fName := flag.String("f", "1.txt", "file to process")
	fChapter := flag.Int64("c", 0, "chapter to process")
	fOutput := flag.String("o", "io.txt", "file for openai")
	flag.Parse()

	fContent, err := os.ReadFile(*fName)
	if err != nil {
		utils.LogPrintError(err)
	}
	chapterRegex := regexp.MustCompile(`第\d+章`)
	chapterList := []string{}
	chapterBuffer := ""
	for _, v := range strings.Split(string(fContent), "\n") {
		if chapterRegex.MatchString(v) {
			if chapterBuffer != "" {
				chapterList = append(chapterList, chapterBuffer)
			}
			chapterBuffer = v
		} else {
			chapterBuffer += "\n" + v
		}
	}
	utils.LogPrintInfo(len(chapterList))
	wrwr := "PROMPT: \n<<>>++__--!!@@##--<<>>\nUSER: 帮我总结一下这一个章节的梗概\n" + chapterList[*fChapter] + "\n<<>>++__--!!@@##--<<>>\n"
	err = os.WriteFile(*fOutput, []byte(wrwr), 0644)
	if err != nil {
		utils.LogPrintError(err)
	}
}
