package myhttp

import (
	"go_utils/utils"
	"net/http"
)

type MiniServer struct {
	name    string
	address string
	port    string
	mux     *http.ServeMux
	ss      *http.Server
}

func NewServer(name, addr, port string) *MiniServer {
	return &MiniServer{
		name:    name,
		address: addr,
		port:    port,
		mux:     http.NewServeMux(),
	}
}

func (s *MiniServer) AddRoute(path string, method []string, h http.HandlerFunc) {
	s.mux.HandleFunc(HandlerMaker(method, path, h))
}

func (s *MiniServer) AddREST(path string, h http.HandlerFunc) {
	s.mux.HandleFunc(
		HandlerMaker(
			[]string{
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
				http.MethodPut,
			},
			path, h,
		),
	)
}

func (s *MiniServer) AddGet(path string, h http.HandlerFunc) {
	s.mux.HandleFunc(HandlerMaker([]string{http.MethodGet}, path, h))
}

func (s *MiniServer) AddPost(path string, h http.HandlerFunc) {
	s.mux.HandleFunc(HandlerMaker([]string{http.MethodPost}, path, h))
}

func (s *MiniServer) Serve() {
	s.ss = &http.Server{
		Addr:    s.address + ":" + s.port,
		Handler: s.mux,
	}
	utils.LogPrintInfo(s.name, "listening on", s.ss.Addr)
	utils.LogPrintError(s.ss.ListenAndServe())
}
