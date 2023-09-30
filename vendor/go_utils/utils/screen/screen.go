package screen

import "fmt"

type ScreenUIEvent interface {
}

type ScreenUIComponent interface {
	Add()
	Draw()
}

type ScreenUIMainFrame struct {
	ChildNodes []ScreenUIComponent
}

func (m *ScreenUIMainFrame) Draw() {
	fmt.Print()
	for _, node := range m.ChildNodes {
		node.Draw()
	}
}

func (m *ScreenUIMainFrame) Run() {

}

type ScreenUIBox struct {
}

// Clear clears the screen
func Clear() {
	fmt.Print("\033[2J")
}

// MoveTopLeft moves the cursor to the top left position of the screen
func MoveTopLeft() {
	fmt.Print("\033[H")
}

// ResetScreen clears screen and move to top left position
func ResetScreen() {
	fmt.Print("\033[2J\033[H")
}
