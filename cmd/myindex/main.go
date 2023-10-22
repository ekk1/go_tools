package main

import (
	"errors"
	"flag"
	"go_utils/utils"
	"go_utils/utils/minikv"
	"net/http"
	"os"
	"sync"
)

var kv *minikv.KV

func main() {
	var verboseFlag = flag.Int("v", 0, "debug (max 4)")
	var listenAddr = flag.String("l", "127.0.0.1", "listen address")
	var listenPort = flag.String("p", "9900", "listen port")
	flag.Parse()

	kvv, err := minikv.NewKV("ss", 0)
	if err != nil {
		utils.LogPrintError(err)
		os.Exit(1)
	}
	kv = kvv
	if err := kv.Load(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			utils.LogPrintWarning("DB not exist, creating...")
			if err2 := kv.Save(); err2 != nil {
				utils.LogPrintError(err2)
				os.Exit(2)
			}
		} else {
			utils.LogPrintError(err)
			utils.LogPrintError("Failed to load DB, exiting")
			os.Exit(2)
		}
	}

	switch *verboseFlag {
	case 0:
		utils.SetLogLevelInfo()
	case 1:
		utils.SetLogLevelDebug()
	case 2:
		utils.SetLogLevelDebug2()
	case 3:
		utils.SetLogLevelDebug3()
	case 4:
		utils.SetLogLevelDebug4()
	}

	muxUser := http.NewServeMux()
	muxUser.HandleFunc("/", handleRoot)
	addrUser := *listenAddr + ":" + *listenPort
	serverUser := http.Server{
		Addr:    addrUser,
		Handler: muxUser,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		utils.LogPrintInfo("User page listening on " + addrUser)
		utils.LogPrintError(serverUser.ListenAndServe())
		wg.Done()
	}()

	wg.Wait()
}