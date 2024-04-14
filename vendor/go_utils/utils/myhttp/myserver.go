package myhttp

import (
	"crypto/tls"
	"crypto/x509"
	"go_utils/utils"
	"net/http"
	"os"
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

func (s *MiniServer) ServeTLS(certFile, keyFile string) {
	s.ss = &http.Server{
		Addr:    s.address + ":" + s.port,
		Handler: s.mux,
	}
	utils.LogPrintInfo(s.name, "listening tls on", s.ss.Addr)
	utils.LogPrintError(s.ss.ListenAndServeTLS(certFile, keyFile))
}

func (s *MiniServer) ServeMutualTLS(certFile, keyFile, caFile string) {
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		utils.LogPrintError("Failed to read ca file:", err)
		return
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		utils.LogPrintError("Failed to add ca cert to pool: ", caFile)
		return
	}

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	// This is deprecated
	// tlsConfig.BuildNameToCertificate()

	s.ss = &http.Server{
		Addr:      s.address + ":" + s.port,
		Handler:   s.mux,
		TLSConfig: tlsConfig,
	}
	utils.LogPrintInfo(s.name, "listening mutual tls on", s.ss.Addr)
	utils.LogPrintError(s.ss.ListenAndServeTLS(certFile, keyFile))
}
