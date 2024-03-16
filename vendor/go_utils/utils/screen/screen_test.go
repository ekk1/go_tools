package screen

import (
	"testing"
)

func TestScreen(t *testing.T) {
	Clear()
	MoveTopLeft()
	ColoredPrintln(ANSI_COLOR_RED, "help")
}
