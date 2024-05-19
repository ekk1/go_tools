package myhttp

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"go_utils/utils"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type HTTPClient struct {
	c                  *http.Client
	json               bool
	form               bool
	rawSend            bool
	headers            http.Header
	username           string
	password           string
	DownloadBufferSize int64
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
		c:                  c,
		headers:            http.Header{},
		DownloadBufferSize: 16 * 1024 * 1024,
	}
}

func countTrue(args ...bool) int64 {
	ret := int64(0)
	for _, v := range args {
		if v {
			ret++
		}
	}
	return ret
}

func (h *HTTPClient) SendGet(sendUrl string, body interface{}) (*HTTPResponse, error) {
	return h.SendReq(http.MethodGet, sendUrl, body)
}
func (h *HTTPClient) SendPost(sendUrl string, body interface{}) (*HTTPResponse, error) {
	return h.SendReq(http.MethodPost, sendUrl, body)
}
func (h *HTTPClient) SendPut(sendUrl string, body interface{}) (*HTTPResponse, error) {
	return h.SendReq(http.MethodPut, sendUrl, body)
}
func (h *HTTPClient) SendDelete(sendUrl string, body interface{}) (*HTTPResponse, error) {
	return h.SendReq(http.MethodDelete, sendUrl, body)
}

func (h *HTTPClient) SendReq(method, sendUrl string, body interface{}) (*HTTPResponse, error) {
	if countTrue(h.json, h.form, h.rawSend) > 1 {
		return nil, errors.New("can not use multiple body types")
	}
	var sendBody io.Reader = nil
	// 1. body is JSON
	if h.json && body != nil {
		sendJson, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		utils.LogPrintDebug3("Sending JSON: ", string(sendJson))
		sendBody = bytes.NewBuffer(sendJson)
	}
	// 2. body is form
	if h.form && body != nil {
		if f, ok := body.(url.Values); ok {
			sendBody = strings.NewReader(f.Encode())
		} else {
			return nil, errors.New("body needs to be url.Values")
		}
	}
	if h.rawSend && body != nil {
		if f, ok := body.([]byte); ok {
			sendBody = bytes.NewReader(f)
		} else {
			return nil, errors.New("body needs to be []byte")
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
		utils.LogPrintError("Failed to do request")
		return nil, err
	}
	defer ret.Body.Close()
	utils.LogPrintDebug2("HTTP Status: ", ret.Status)
	utils.LogPrintDebug2("HTTP Headers: ", ret.Header)

	data, err := io.ReadAll(ret.Body)
	if err != nil {
		return nil, err
	}
	utils.LogPrintDebug4("HTTP BodyBytes: ", data)
	utils.LogPrintDebug3("HTTP BodyString: ", string(data))

	return &HTTPResponse{Data: data, Status: ret.Status}, nil
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

func (h *HTTPClient) SetSendRawBody(isSendRawBody bool) bool {
	h.rawSend = isSendRawBody
	h.SetHeader("Content-Type", "application/octet-stream")
	return h.rawSend
}

func (h *HTTPClient) SetSendMultiPartForm(isSendForm bool) bool {
	h.form = isSendForm
	h.SetHeader("Content-Type", "multipart/form-data")
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
		return errors.New("failed to get transport")
	}
	return nil
}

func (h *HTTPClient) SetDisableCompress() error {
	if ts, ok := h.c.Transport.(*http.Transport); ok {
		ts.DisableCompression = true
	} else {
		return errors.New("failed to get transport")
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
			utils.LogPrintError("Failed to load cert: ", cert)
			continue
		}
		if ok := certPool.AppendCertsFromPEM(pem); !ok {
			utils.LogPrintError("Failed to apply cert: ", cert)
		}
		utils.LogPrintDebug("Adding cert: ", cert)
		didAddCert = true
	}
	if didAddCert {
		if ts, ok := h.c.Transport.(*http.Transport); ok {
			if ts.TLSClientConfig != nil {
				ts.TLSClientConfig.RootCAs = certPool
			} else {
				ts.TLSClientConfig = &tls.Config{RootCAs: certPool}
			}
		} else {
			return errors.New("failed to get transport")
		}
	}
	return nil
}

// Usually when using client auth, server auth is also applied, so this function forces server tls auth
func (h *HTTPClient) SetClientCert(certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		utils.LogPrintError("Falied to load client cert: ", err)
		return err
	}
	if ts, ok := h.c.Transport.(*http.Transport); ok {
		if ts.TLSClientConfig != nil {
			ts.TLSClientConfig.Certificates = []tls.Certificate{cert}
		} else {
			ts.TLSClientConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
		}
	} else {
		return errors.New("failed to get transport")
	}

	return nil
}

func (h *HTTPClient) DownloadFile(url, downloadName string) (int64, error) {
	var downloadBuffer []byte = make([]byte, h.DownloadBufferSize)
	totalRead := int64(0)
	totalWritten := int64(0)

	f, err := os.OpenFile(downloadName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	// Check file info, and do not download if alread finished
	fInfo, err := f.Stat()
	if err != nil {
		return 0, err
	}
	totalWritten = fInfo.Size()
	if totalWritten != 0 {
		h.SetHeader("Range", fmt.Sprintf("bytes=%d-", totalWritten))
	}

	// Prepare HTTP request and get the response
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	for k, v := range h.headers {
		req.Header.Set(k, v[0])
	}
	if h.username != "" && h.password != "" {
		req.SetBasicAuth(h.username, h.password)
	}
	utils.LogPrintDebug(req.Header)
	res, err := h.c.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusRequestedRangeNotSatisfiable {
		utils.LogPrintWarning("Incorrect range, file may have already been downloaded")
		return 0, nil
	}

	if totalWritten != 0 && res.StatusCode == http.StatusOK {
		utils.LogPrintWarning("Server may not support partial download, redownloading all contents")
		f.Close()
		f, err = os.OpenFile(downloadName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
		if err != nil {
			return 0, err
		}
		defer f.Close()
		totalWritten = 0
	}

	// 1. range not set, ContentLength == expected size
	// 2. range set, and is resume: ContentLength == (expected size) - downloaded
	// 3. range set, and overflow: ContentLength == 196 (for nginx), and should not download
	// 4. range set, and not accepted by server, ContentLength == expected size
	totalSize := res.ContentLength
	if totalSize == -1 {
		utils.LogPrintWarning("Cannot get correct content size")
	}
	totalWritten = 0

	for {
		time.Sleep(3 * time.Second)
		num, errRead := io.ReadFull(res.Body, downloadBuffer)
		utils.LogPrintDebug(num, errRead)
		if num > 0 {
			totalRead += int64(num)
			written, errWrite := f.Write(downloadBuffer[:num])
			totalWritten += int64(written)
			if errWrite != nil {
				return totalWritten, errWrite
			}
			utils.LogPrintInfo(
				fmt.Sprintf(
					"Downloading...  %dMB / %dMB",
					int64(totalRead/1024/1024),
					int64(totalSize/1024/1024),
				),
			)
		}
		if errRead != nil {
			if errors.Is(errRead, io.EOF) ||
				errors.Is(errRead, io.ErrUnexpectedEOF) {
				if totalRead == totalSize {
					break
				} else if totalSize == -1 {
					utils.LogPrintError("Cannot determine total size, and EOF is met, please check download manually: ", errRead)
					return totalWritten, nil
				} else {
					utils.LogPrintError("Incompleted download: ", errRead)
					utils.LogPrintWarning("Total read: ", totalRead, " Total size: ", totalSize)
					return totalWritten, errors.New("incompleted download")
				}
			} else {
				utils.LogPrintError("Error during download: ", errRead)
				return totalWritten, errRead
			}
		}
	}

	if totalWritten == totalSize {
		return totalWritten, nil
	}
	return totalWritten, errors.New(fmt.Sprintf(
		"Unexpected totalWritten not equal totalSize: %d / %d", totalWritten, totalSize,
	))
}

type HTTPResponse struct {
	Status string
	Data   []byte
}

func (r *HTTPResponse) JSON(recvStruct interface{}) error {
	if recvStruct == nil {
		return errors.New("receive struct is nil")
	}
	return json.Unmarshal(r.Data, recvStruct)
}

func (r *HTTPResponse) Text() string {
	return string(r.Data)
}
