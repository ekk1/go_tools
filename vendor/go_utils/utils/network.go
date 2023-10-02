package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type HTTPClient struct {
	c        *http.Client
	json     bool
	form     bool
	headers  http.Header
	username string
	password string
}

func NewHTTPClient() *HTTPClient {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	ts := &http.Transport{
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c := &http.Client{
		Transport: ts,
	}
	return &HTTPClient{
		c:       c,
		headers: http.Header{},
	}
}

func (h *HTTPClient) SendReq(method, sendUrl string, body interface{}) (*HTTPResponse, error) {
	if h.json && h.form {
		return nil, errors.New("json and form cannot both be set")
	}
	var sendBody io.Reader = nil
	// 1. body is JSON
	if h.json && body != nil {
		sendJson, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		LogPrintDebug3("Sending JSON: ", string(sendJson))
		sendBody = bytes.NewBuffer(sendJson)
	}
	// 2. body is form
	if h.form && body != nil {
		if f, ok := body.(url.Values); ok {
			sendBody = strings.NewReader(f.Encode())
		} else {
			return nil, errors.New("Body needs to be url.Values")
		}
	}
	// 3.
	req, err := http.NewRequest(method, sendUrl, sendBody)
	if err != nil {
		return nil, err
	}

	for k, v := range h.headers {
		req.Header.Set(k, v[0])
	}

	if h.username != "" && h.password != "" {
		req.SetBasicAuth(h.username, h.password)
	}

	ret, err := h.c.Do(req)
	if err != nil {
		return nil, err
	}
	LogPrintDebug2("HTTP Status: ", ret.Status)
	LogPrintDebug2("HTTP Headers: ", ret.Header)
	defer ret.Body.Close()

	data, err := io.ReadAll(ret.Body)
	if err != nil {
		return nil, err
	}
	LogPrintDebug4("HTTP BodyBytes: ", data)
	LogPrintDebug3("HTTP BodyString: ", string(data))

	return &HTTPResponse{data: data}, nil
}

func (h *HTTPClient) SetSendJSON(isSendJSON bool) bool {
	h.json = isSendJSON
	h.SetHeader("Content-Type", "application/json; charset=UTF-8")
	return h.json
}

func (h *HTTPClient) SetSendForm(isSendForm bool) bool {
	h.form = isSendForm
	h.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	return h.form
}

func (h *HTTPClient) SetHeader(key, value string) {
	h.headers.Set(key, value)
}

func (h *HTTPClient) SetProxy(proxyURL string) error {
	pURL, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}
	if ts, ok := h.c.Transport.(*http.Transport); ok {
		ts.Proxy = http.ProxyURL(pURL)
	} else {
		return errors.New("Failed to get transport")
	}
	return nil
}

func (h *HTTPClient) SetBasicAuth(username, password string) {
	h.username = username
	h.password = password
}

func (h *HTTPClient) SetCustomCert(certPath []string) error {
	certPool := x509.NewCertPool()
	didAddCert := false
	for _, cert := range certPath {
		pem, err := os.ReadFile(cert)
		if err != nil {
			LogPrintError("Failed to load cert: ", cert)
			continue
		}
		if ok := certPool.AppendCertsFromPEM(pem); !ok {
			LogPrintError("Failed to apply cert: ", cert)
		}
		LogPrintDebug("Adding cert: ", cert)
		didAddCert = true
	}
	if didAddCert {
		if ts, ok := h.c.Transport.(*http.Transport); ok {
			ts.TLSClientConfig = &tls.Config{RootCAs: certPool}
		} else {
			return errors.New("Failed to get transport")
		}
	}
	return nil
}

type HTTPResponse struct {
	data []byte
}

func (r *HTTPResponse) JSON(recvStruct interface{}) error {
	if recvStruct == nil {
		return errors.New("Receive struct is nil!!!")
	}
	return json.Unmarshal(r.data, recvStruct)
}

func (r *HTTPResponse) Text() string {
	return string(r.data)
}
