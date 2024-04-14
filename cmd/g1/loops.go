package main

import (
	"fmt"
	"go_utils/utils"
	"go_utils/utils/screen"
	"os"
	"time"
)

func UpdateWorld() {
	for _, v := range GlobalItemList {
		v.Next()
	}
	for _, v := range GlobalBuildings {
		v.Next()
	}
	for _, v := range GlobalUnits {
		v.Next()
	}
}

func DrawWorld() {
	screen.ResetScreen()
	fmt.Printf("(q to quit) PosX: %d PosY: %d\n", CurrentPosX, CurrentPosY)
	fmt.Print(GlobalMap.Render(CurrentPosX, CurrentPosY, CurrentViewRange))

	b := GlobalMap.Blocks[0][CurrentPosY][CurrentPosX]
	fmt.Println("-------------")
	fmt.Println(b.Info())
}

func RenderLoop() {
	for {
		breakFlag := false
		select {
		case e := <-EventChannel:
			switch e {
			case EventCameraUP:
				if CurrentPosY > 0 {
					CurrentPosY -= 1
				}
			case EventCameraDown:
				if CurrentPosY < GlobalMap.SizeY-1 {
					CurrentPosY += 1
				}
			case EventCameraLeft:
				if CurrentPosX > 0 {
					CurrentPosX -= 1
				}
			case EventCameraRight:
				if CurrentPosX < GlobalMap.SizeX-1 {
					CurrentPosX += 1
				}

			case EventPlant:
				fmt.Println("Waiting:")

				PendingEventChannel <- PendingEventPlant
				DrawWorld()
				fmt.Println("Waiting:")
				selection := <-PendingEventSelection
				GlobalMap.Blocks[0][CurrentPosY][CurrentPosX].Plant(
					GlobalItemConfigTable.PlantList[selection-48],
				)
			case EventGameStop:
				ExitChannelRenderLoop <- 0
			}
		case <-GlobalTicker.C:
			UpdateWorld()
			DrawWorld()
		}
		if breakFlag {
			utils.LogPrintInfo("Render loop exits")
			break
		}
	}

}

var EscapePhase1 bool
var EscapePhase2 bool

func ResetEscapePhase() {
	EscapePhase1 = false
	EscapePhase2 = false
}
func IsNotInEscapePhase() bool {
	return (!EscapePhase1 && !EscapePhase2)
}
func IsEscapePhase1() bool {
	return (EscapePhase1 && !EscapePhase2)
}
func IsEscapePhase2() bool {
	return (EscapePhase1 && EscapePhase2)
}

func InputLoop() {
	var b []byte = make([]byte, 1)
	ResetEscapePhase()
	for {
		os.Stdin.Read(b)
		breakFlag := false

		switch {
		case b[0] == 104: // h
			EventChannel <- EventCameraLeft
		case b[0] == 106: // j
			EventChannel <- EventCameraDown
		case b[0] == 107: // k
			EventChannel <- EventCameraUP
		case b[0] == 108: // l
			EventChannel <- EventCameraRight
		case b[0] == 112: // p
			EventChannel <- EventPlant
		case b[0] == 27: // ESC started
			if IsNotInEscapePhase() {
				EscapePhase1 = true
			}
		case b[0] == 91: // [
			if IsEscapePhase1() {
				EscapePhase2 = true
			} else {
				ResetEscapePhase()
			}
		case b[0] == 65: // A mostly for escape up
			if IsEscapePhase2() {
				EventChannel <- EventCameraUP
			}
			ResetEscapePhase()
		case b[0] == 66: // B mostly for escape down
			if IsEscapePhase2() {
				EventChannel <- EventCameraDown
			}
			ResetEscapePhase()
		case b[0] == 67: // C mostly for escape right
			if IsEscapePhase2() {
				EventChannel <- EventCameraRight
			}
			ResetEscapePhase()
		case b[0] == 68: // D mostly for escape left
			if IsEscapePhase2() {
				EventChannel <- EventCameraLeft
			}
			ResetEscapePhase()
		case b[0] >= 48 && b[0] <= 57:
			PendingEventSelection <- b[0]
		case b[0] == 113: // q
			EventChannel <- EventGameStop
			breakFlag = true
		default:
			fmt.Println(b)
			time.Sleep(1 * time.Second)
		}

		if breakFlag {
			utils.LogPrintInfo("Input loop exits")
			ExitChannelInputLoop <- 0
			break
		}
	}
}
