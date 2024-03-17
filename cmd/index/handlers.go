package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"go_utils/utils/myhttp"
	"net/http"
	"os"
	"path"
)

func handleSave(w http.ResponseWriter, req *http.Request) {
	storeName := req.FormValue("name")
	storeIndex := req.FormValue("index")
	storeData := req.FormValue("data")
	storeChecksum := req.FormValue("check")

	if !myhttp.ServerCheckParam(
		storeName,
		storeIndex,
		storeChecksum,
		storeData,
	) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing parameters"))
		return
	}

	dataBytes, err := base64.URLEncoding.DecodeString(storeData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed"))
		return
	}

	h := sha256.New()
	sum := h.Sum(dataBytes)
	if base64.URLEncoding.EncodeToString(sum) != storeChecksum {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad checksum"))
		return
	}

	saveDir := path.Join(baseDir, storeName)

	_, err = os.Stat(saveDir)
	if os.IsNotExist(err) {
		err2 := os.MkdirAll(saveDir, 0o700)
		if err2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed creating dir"))
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed read dir"))
		return
	}

	saveObject := path.Join(saveDir, storeName+"_"+storeIndex)
	if err := os.WriteFile(saveObject, bytes.Join([][]byte{sum, dataBytes}, []byte{}), 0600); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed writing data"))
		return
	}
	w.Write([]byte("OK"))
}

func handleLoad(w http.ResponseWriter, req *http.Request) {
	storeName := req.FormValue("name")
	storeIndex := req.FormValue("index")

	if !myhttp.ServerCheckParam(
		storeName,
		storeIndex,
	) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing parameters"))
		return
	}
	saveDir := path.Join(baseDir, storeName)
	saveObject := path.Join(saveDir, storeName+"_"+storeIndex)

	if dataBytes, err := os.ReadFile(saveObject); err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Missing object"))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed read"))
			return
		}
	} else {
		dataSum := dataBytes[:32]
		readData := dataBytes[32:]
		h := sha256.New()
		sum := h.Sum(readData)
		if !bytes.Equal(dataSum, sum) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Corrupt data"))
			return
		}
		w.Write([]byte(base64.URLEncoding.EncodeToString(readData)))
		return
	}

}

func handleRebalance(w http.ResponseWriter, req *http.Request) {

}

func handleScrub(w http.ResponseWriter, req *http.Request) {

}

func handleStat(w http.ResponseWriter, req *http.Request) {

}

func handleDebug(w http.ResponseWriter, req *http.Request) {

}
