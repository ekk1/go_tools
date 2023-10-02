package main

import (
	_ "embed"
	"flag"
	"go_utils/utils"
	"html/template"
	"io"
	"net/http"
	"os"
	"sync"
)

//go:embed index.html
var indexHTML string

var indexTemplate = template.Must(template.New("index").Parse(indexHTML))

type Subscribe struct {
	URL         string
	LastUpdated string
	Content     string
}

func (s *Subscribe) Update() error {
	return nil
}

func (s *Subscribe) Save() error {
	return nil
}

func LoadSubscribe() ([]*Subscribe, error) {
	return nil, nil

}

type PageData struct {
	Title     string
	Msg       string
	Lister    []string
	Mapper    map[string]string
	Subscribe []*Subscribe
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	utils.ServerLog("root", req)

	// Deny req not to /
	if req.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(""))
	}

	// Print headers if debug
	for k, v := range req.Header {
		utils.LogPrintDebug(k, v)
	}
	if req.Method == http.MethodPost {
		postData := req.FormValue("name_list")
		if postData == "" {
			utils.LogPrintWarning("No postData")
		} else {
			utils.LogPrintInfo("POST DATA:", postData)
		}
		f, fHeader, err := req.FormFile("testfile")
		if err != nil {
			utils.LogPrintWarning("No file")
			utils.LogPrintWarning(err)
		} else {
			utils.LogPrintInfo(fHeader)
			defer f.Close()
			data, _ := io.ReadAll(f)
			os.WriteFile("tetset", data, 0666)
		}
	}

	// if err := indexTemplate.ExecuteTemplate(w, "index", &PageData{Title: "test", Info: "testinfo"}); err != nil {
	pageData := &PageData{
		Title:  "test",
		Msg:    "msg",
		Lister: []string{"1", "2"},
		Mapper: map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		},
	}
	if err := indexTemplate.ExecuteTemplate(w, "index", pageData); err != nil {
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
