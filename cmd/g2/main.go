package main

import (
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"sync"
	"time"
)

type PlayerEvent struct {
}

var (
	GlobalEventChannel = make(chan *PlayerEvent, 100)
	GlobalTicker       = time.NewTicker(10 * time.Second)

	ExitChannel = make(chan string)

	GlobalCityList = map[string]*City{}
)

func RenderLoop() {
	for {
		select {
		case e := <-GlobalEventChannel:
			utils.LogPrintInfo(e)
		case <-GlobalTicker.C:
			var wg sync.WaitGroup
			wg.Add(len(GlobalCityList))
			for n, c := range GlobalCityList {
				utils.LogPrintInfo("Updating " + n)
				go func() {
					c.Next()
					wg.Done()
				}()
			}
			wg.Wait()
		}
	}

}

func main() {
	// Run main loop
	InitConfig()

	go RenderLoop()

	cc := myhttp.NewServer("gas", "127.0.0.1", "9999")
	cc.AddGet("/", handleServeIndex)
	cc.AddGet("/post", handleUserInput)
	go cc.Serve()

	<-ExitChannel
}
