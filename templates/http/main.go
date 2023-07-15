package main

import (
	"errors"
	"fmt"
	"go_utils/utils"
	"go_utils/utils/quickserver"
	"net/http"
	"sync"
)

// main.go
func denyUnDefinedResouce(w http.ResponseWriter, req *http.Request) error {
	if req.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(""))
		return errors.New(fmt.Sprintf("Resource not defined for %s on %s", req.URL.Path, req.Host))
	}
	return nil
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	quickserver.QuickServerLog(req, "handleRoot")

	if err := denyUnDefinedResouce(w, req); err != nil {
		utils.LogPrintError(err)
		return
	}
}

func main() {
	prepareAssetDict()

	muxUser := http.NewServeMux()
	muxUser.HandleFunc("/", handleRoot)
	muxUser.HandleFunc("/static/", handleStatic)
	addrUser := "127.0.0.1:8080"
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
