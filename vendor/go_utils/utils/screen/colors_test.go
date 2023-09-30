package screen

import (
	"testing"
)

func TestColors(t *testing.T) {
	ResetScreen()

	ColoredPrintln(ANSI_STYLE_BOLD, "Bold ", "Hello")
	ColoredPrintln(ANSI_STYLE_BLINK, "Blink ", "Hello")
	ColoredPrintln(ANSI_STYLE_DIM, "Dim ", "Hello")
	ColoredPrintln(ANSI_STYLE_ITALIC, "Italic ", "Hello")
	ColoredPrintln(ANSI_STYLE_STRIKE, "Strike ", "Hello")
	ColoredPrintln(ANSI_STYLE_UNDERLINE, "Underline ", "Hello")

	ColoredPrintln(ANSI_COLOR_BLACK, "Black ", "Hello")
	ColoredPrintln(ANSI_COLOR_RED, "Red ", "Hello")
	ColoredPrintln(ANSI_COLOR_BLUE, "Blue ", "Hello")
	ColoredPrintln(ANSI_COLOR_CYAN, "Cyan ", "Hello")
	ColoredPrintln(ANSI_COLOR_GREEN, "Green ", "Hello")
	ColoredPrintln(ANSI_COLOR_MAGENTA, "Magenta ", "Hello")
	ColoredPrintln(ANSI_COLOR_YELLOW, "Yellow ", "Hello")
	ColoredPrintln(ANSI_COLOR_WHITE, "White ", "Hello")

	ColoredPrintln(ANSI_COLOR_RED+ANSI_STYLE_BOLD, "Red Bold ", "Hello")
	ColoredPrintln(ANSI_COLOR_YELLOW+ANSI_STYLE_BOLD, "Yellow Bold ", "Hello")
}
