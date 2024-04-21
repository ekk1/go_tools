package main

import (
	"go_tools/cmd/g2/config"
	"go_tools/cmd/g2/event"
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"sync"
	"time"
)

var (
	GlobalEventChannel = make(chan *event.PlayerEvent, 100)
	GlobalTicker       = time.NewTicker(10 * time.Second)

	ExitChannel = make(chan string)

	GlobalCityList = map[string]config.City{}
)

func RenderLoop() {
	for {
		select {
		case e := <-GlobalEventChannel:
			utils.LogPrintInfo(e)
			<-e.Finished
		case <-GlobalTicker.C:
			var wg sync.WaitGroup
			wg.Add(len(GlobalCityList))
			for n, c := range GlobalCityList {
				utils.LogPrintInfo("Updating " + n)
				go func() {
					c.Update()
					wg.Done()
				}()
			}
			wg.Wait()
		}
	}

}

func main() {
	// Run main loop
	config.InitConfig()

	go RenderLoop()

	cc := myhttp.NewServer("gas", "127.0.0.1", "9999")
	cc.AddGet("/", handleServeIndex)
	cc.AddGet("/post", handleUserInput)
	go cc.Serve()

	<-ExitChannel
}
