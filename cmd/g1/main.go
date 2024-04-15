package main

import (
	"fmt"
	"go_utils/utils"
	"os/exec"
	"time"
)

var (
	CurrentPosX      int64 = 0
	CurrentPosY      int64 = 0
	CurrentViewRange int64 = 10

	GlobalItemList  []Items
	GlobalBuildings []Building
	GlobalUnits     []Unit

	GlobalMap *Map

	GlobalTicker *time.Ticker
	EventChannel chan GameEvent

	ExitChannelInputLoop  chan byte
	ExitChannelRenderLoop chan byte
)

func main() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	GlobalMap = NewMap(1000, 1000, 1)

	if err := LoadConfigTable(); err != nil {
		panic(err)
	}

	fmt.Println("Finish generate map: ", GlobalMap.SizeX, GlobalMap.SizeY)

	GlobalTicker = time.NewTicker(100 * time.Millisecond)
	EventChannel = make(chan GameEvent)
	ExitChannelInputLoop = make(chan byte)
	ExitChannelRenderLoop = make(chan byte)

	go RenderLoop()
	go InputLoop()

	<-ExitChannelInputLoop
	<-ExitChannelRenderLoop

	utils.LogPrintInfo("All loops exited")
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}
