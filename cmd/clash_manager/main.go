package main

import (
	_ "embed"
	"flag"
	"go_utils/utils"
	"html/template"
	"io"
	"net/http"
	"sync"
)

//go:embed index.html
var indexHTML string

var indexTemplate *template.Template

type Page struct {
	Title string
	Info  string
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(""))
	}
	utils.LogPrintInfo(
		req.RemoteAddr,
		req.Method,
		req.URL.Path,
		req.Header,
	)
	data, err := io.ReadAll(req.Body)
	if err != nil {
		utils.LogPrintError(err)
	}
	utils.LogPrintInfo(string(data))
	if err := req.ParseForm(); err != nil {
		utils.LogPrintError("Failed to parse form")
		utils.LogPrintError(err)
	}
	utils.LogPrintInfo(
		req.Form,
		req.PostForm,
		req.FormValue("name_list"),
		req.FormValue("name_list2"),
	)
	if err := indexTemplate.ExecuteTemplate(w, "index", &Page{Title: "test", Info: "testinfo"}); err != nil {
		utils.LogPrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to execute template"))
	}
}

func main() {
	var verboseFlag = flag.Int("v", 0, "debug (max 4)")
	flag.Parse()
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

	t := template.New("index")
	indexTemplate = template.Must(t.Parse(indexHTML))
	utils.LogPrintInfo(indexTemplate)

	muxUser := http.NewServeMux()
	muxUser.HandleFunc("/", handleRoot)
	addrUser := "127.0.0.1:8888"
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
