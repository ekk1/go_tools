package screen

import "fmt"

const (
	ANSI_COLOR_BLACK   string = "\033[30m"
	ANSI_COLOR_RED     string = "\033[31m"
	ANSI_COLOR_GREEN   string = "\033[32m"
	ANSI_COLOR_YELLOW  string = "\033[33m"
	ANSI_COLOR_BLUE    string = "\033[34m"
	ANSI_COLOR_MAGENTA string = "\033[35m"
	ANSI_COLOR_CYAN    string = "\033[36m"
	ANSI_COLOR_WHITE   string = "\033[37m"
	ANSI_COLOR_DEFAULT string = "\033[39m"
	ANSI_COLOR_RESET   string = "\033[0m"
)

const (
	ANSI_STYLE_BOLD      string = "\033[1m"
	ANSI_STYLE_DIM       string = "\033[2m"
	ANSI_STYLE_ITALIC    string = "\033[3m"
	ANSI_STYLE_UNDERLINE string = "\033[4m"
	ANSI_STYLE_BLINK     string = "\033[5m"
	ANSI_STYLE_STRIKE    string = "\033[9m"
)

func ColoredPrint(style string, msg ...any) {
	fmt.Print(style)
	fmt.Print(msg...)
	fmt.Print(ANSI_COLOR_RESET)
}

func ColoredPrintln(style string, msg ...any) {
	fmt.Print(style)
	fmt.Println(msg...)
	fmt.Print(ANSI_COLOR_RESET)
}
