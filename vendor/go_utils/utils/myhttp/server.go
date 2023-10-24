package myhttp

import (
	"fmt"
	"go_utils/utils"
	"net/http"
	"slices"
	"strings"
	"sync"
)

func ServerCheckPath(expect string, req *http.Request, w http.ResponseWriter) bool {
	if req.URL.Path != expect {
		ServerError("Not found", w, req)
		return false
	}
	return true
}

func ServerCheckParam(args ...string) bool {
	for _, v := range args {
		if v == "" {
			return false
		}
	}
	return true
}

func ServerDebugHeader(r *http.Request) {
	for k, v := range r.Header {
		utils.LogPrintDebug2("Header:", k, v)
	}
}

func ServerLog(caller string, r *http.Request) {
	utils.LogPrintInfo(fmt.Sprintf(
		"[%s]: got %s from %s for %s, Host: %s",
		caller, r.Method, r.RemoteAddr, r.URL.Path, r.Host,
	))
}

func ServerReply(msg string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html lang=\"en\"><head><link rel=\"icon\" href=\"data:;base64,iVBORw0KGgo=\"></head><body><p>" + msg + "</p>\n\n<a href=\"/\">index</a></body></html>"))
}

func ServerError(msg string, w http.ResponseWriter, r *http.Request) {
	utils.LogPrintError("Got", r.Method, "from", r.RemoteAddr, "to", r.URL.Path, "Failed")
	//w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html lang=\"en\"><head><link rel=\"icon\" href=\"data:;base64,iVBORw0KGgo=\"></head><body><p>" + msg + "</p>\n\n<a href=\"/\">index</a></body></html>"))
}

func HandlerMaker(method []string, path string, h http.HandlerFunc) (string, http.HandlerFunc) {
	return path, func(ww http.ResponseWriter, rr *http.Request) {
		ServerLog(path, rr)
		if !ServerCheckPath(path, rr, ww) {
			ServerError("Path error", ww, rr)
			return
		}
		if !slices.Contains(method, rr.Method) {
			ServerError("Expect "+strings.Join(method, ",")+" but got: "+rr.Method, ww, rr)
			return
		}
		h(ww, rr)
	}
}

func TokenChecker(tokenKey string, allowedTokens []string, h http.HandlerFunc) http.HandlerFunc {
	return func(ww http.ResponseWriter, rr *http.Request) {
		token := rr.Header.Get(tokenKey)
		if !slices.Contains(allowedTokens, token) {
			ServerError("Not Allowed", ww, rr)
			return
		}
		h(ww, rr)
	}
}

func RunServers(s ...*MiniServer) {
	var wg sync.WaitGroup

	wg.Add(len(s))

	for _, v := range s {
		go func(ss *MiniServer) {
			ss.Serve()
			wg.Done()
		}(v)
	}
	wg.Wait()
	utils.LogPrintWarning("All servers were shutdown")
}
