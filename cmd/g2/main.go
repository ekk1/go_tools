package main

import (
	"embed"
	"go_tools/cmd/g2/config"
	"go_tools/cmd/g2/event"
	"go_tools/cmd/g2/player"
	"go_utils/utils"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"
)

//go:embed static
var staticFiles embed.FS

var (
	GlobalEventChannel = make(chan *event.PlayerEvent, 100)
	GlobalTicker       = time.NewTicker(10 * time.Second)

	ExitChannel = make(chan string)

	// GlobalCityList   = map[string]config.City{}
	GlobalPlayerList = map[string]*player.PlayerStruct{}
)

func UpdateWorld() {
	var wg sync.WaitGroup
	for n, p := range GlobalPlayerList {
		utils.LogPrintInfo("Updating player" + n)
		wg.Add(len(p.CityList))
		for cN, c := range p.CityList {
			utils.LogPrintInfo("Updating city" + cN)
			go func() {
				c.Update()
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func RenderLoop() {
	for {
		select {
		case e := <-GlobalEventChannel:
			select {
			case <-GlobalTicker.C:
				UpdateWorld()
			default:
			}
			utils.LogPrintInfo(e)
			<-e.Finished
		case <-GlobalTicker.C:
			UpdateWorld()
		}
	}

}

func main() {
	// Run main loop
	config.InitConfig()

	go RenderLoop()

	var staticFS = fs.FS(staticFiles)
	htmlContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(htmlContent))

	// cc := myhttp.NewServer("gas", "127.0.0.1", "9999")
	ccMux := http.NewServeMux()
	ccMux.Handle("/", fs)
	ccMux.HandleFunc("/post", handleUserInput)

	ss := &http.Server{
		Addr:    "127.0.0.1:9999",
		Handler: ccMux,
	}

	go ss.ListenAndServe()

	<-ExitChannel
}
