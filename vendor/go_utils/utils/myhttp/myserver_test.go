package myhttp

import (
	"fmt"
	"go_utils/utils"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	s := NewServer("ss", "127.0.0.1", "10001")
	s.mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("Hello")) })
	go RunServers(s)

	time.Sleep(1 * time.Second)

	totalReq := 1000000

	ch := make(chan int, 1)
	go func() {
		c := NewHTTPClient()
		t1 := time.Now()
		for ii := 0; ii < totalReq; ii++ {
			_, err := c.SendGet("http://127.0.0.1:10001/", nil)
			if err != nil {
				utils.LogPrintError(err)
				t.Fail()
				ch <- 1
				return
			}
		}
		t2 := time.Now().Sub(t1)
		utils.LogPrintInfo(fmt.Sprintf("HTTP Speed: %.2f/s", float64(totalReq)/(t2.Seconds())))
		ch <- 1
	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
		//s.ss.Close()
	}()

	<-ch

	time.Sleep(5 * time.Second)

	s.ss.Close()
}

func TestMutualTLSServer(t *testing.T) {
	s := NewServer("ss", "127.0.0.1", "10001")
	s.mux.HandleFunc("/", ClientTLSChecker(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello"))
	}))
	go s.ServeMutualTLS("server.crt", "server.key", "ca.crt")

	time.Sleep(1 * time.Second)

	totalReq := 1000000

	ch := make(chan int, 1)
	go func() {
		c := NewHTTPClient()
		c.SetCustomCert([]string{"ca.crt"})
		c.SetClientCert("client.crt", "client.key")
		t1 := time.Now()
		for ii := 0; ii < totalReq; ii++ {
			ret, err := c.SendGet("https://127.0.0.1:10001/", nil)
			if err != nil {
				utils.LogPrintError(err)
				t.Fail()
				ch <- 1
				return
			}
			if totalReq == 1 {
				utils.LogPrintInfo(ret.Status)
				utils.LogPrintInfo(ret.Text())
			}
		}
		t2 := time.Now().Sub(t1)
		utils.LogPrintInfo(fmt.Sprintf("HTTP Speed: %.2f/s", float64(totalReq)/(t2.Seconds())))
		ch <- 1
	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
		//s.ss.Close()
	}()

	<-ch

	time.Sleep(1 * time.Second)

	s.ss.Close()
}

func TestGenerateCerts(t *testing.T) {
	caPrivKey, err := GenerateRSAKey()
	if err != nil {
		t.Fatal(err)
	}
	serverPrivKey, err := GenerateRSAKey()
	if err != nil {
		t.Fatal(err)
	}
	clientPrivKey, err := GenerateRSAKey()
	if err != nil {
		t.Fatal(err)
	}

	caCert, err := GenerateCACert(caPrivKey, 10)
	if err != nil {
		t.Fatal(err)
	}
	serverCert, err := GenerateCertWithCA(
		CertTypeServer,
		1,
		caPrivKey,
		&serverPrivKey.PublicKey,
		caCert,
		[]string{"127.0.0.1"},
	)
	if err != nil {
		t.Fatal(err)
	}
	clientCert, err := GenerateCertWithCA(
		CertTypeClient,
		1,
		caPrivKey,
		&clientPrivKey.PublicKey,
		caCert,
		[]string{"127.0.0.2", "127.0.0.1"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := WriteCertAndKey("ca", caCert, caPrivKey); err != nil {
		t.Fatal(err)
	}
	if err := WriteCertAndKey("server", serverCert, serverPrivKey); err != nil {
		t.Fatal(err)
	}
	if err := WriteCertAndKey("client", clientCert, clientPrivKey); err != nil {
		t.Fatal(err)
	}

}
