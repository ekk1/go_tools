package utils

import (
	"fmt"
	"net/http"
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
		LogPrintDebug2("Header:", k, v)
	}
}

func ServerLog(caller string, r *http.Request) {
	LogPrintInfo(fmt.Sprintf(
		"[%s]: got %s from %s for %s, Host: %s",
		caller, r.Method, r.RemoteAddr, r.URL.Path, r.Host,
	))
}

func ServerReply(msg string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func ServerError(msg string, w http.ResponseWriter, r *http.Request) {
	LogPrintInfo("Got", r.Method, "from", r.RemoteAddr, "to", r.URL.Path, "Failed")
	//w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html lang=\"en\"><body><p>" + msg + "</p>\n\n<a href=\"/\">index</a></body></html>"))
}
