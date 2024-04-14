package main

import (
	"fmt"
	"go_utils/utils/screen"
	"os"
	"os/exec"
)

var (
	CurrentPosX      int64 = 0
	CurrentPosY      int64 = 0
	CurrentViewRange int64 = 10
)

func main() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	m := NewMap(1000, 1000, 1)

	LoadConfigTable()

	// m.Blocks[0][0][0].Build(&Farm{})
	fmt.Println("Finish Generate: ", m.SizeX)

	var b []byte = make([]byte, 1)

	for {
		screen.ResetScreen()
		fmt.Printf("(q to quit) PosX: %d PosY: %d\n", CurrentPosX, CurrentPosY)
		fmt.Print(m.Render(CurrentPosX, CurrentPosY, CurrentViewRange))

		os.Stdin.Read(b)

		fmt.Println(b)

		// break
		quitFlag := false

		switch b[0] {
		case 104:
			if CurrentPosX > 0 {
				CurrentPosX -= 1
			}
		case 106:
			if CurrentPosY < m.SizeY-1 {
				CurrentPosY += 1
			}
		case 107:
			if CurrentPosY > 0 {
				CurrentPosY -= 1
			}
		case 108:
			if CurrentPosX < m.SizeX-1 {
				CurrentPosX += 1
			}
		case 113:
			quitFlag = true
		}

		if quitFlag {
			break
		}
	}
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}
