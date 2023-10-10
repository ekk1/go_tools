package main

import (
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"net/http"
)

const (
	clashManagerAddr = "http://127.0.0.1:9090"
)

type Proxy struct {
	All  []string `json:"all"`
	Name string   `json:"name"`
	Type string   `json:"type"`
	Now  string   `json:"now"`
}

func GetNodesByProxy(p string) ([]string, string, error) {
	c := myhttp.NewHTTPClient()
	ret, err := c.SendReq(http.MethodGet, clashManagerAddr+"/proxies/"+p, nil)
	if err != nil {
		return nil, "", err
	}
	utils.LogPrintDebug(ret.Text())
	utils.LogPrintDebug(ret)
	var pp *Proxy = &Proxy{}
	if err := ret.JSON(pp); err != nil {
		return nil, "", err
	}

	return pp.All, pp.Now, nil
}

func ChangeNodeForProxy(p, n string) error {
	c := myhttp.NewHTTPClient()
	c.SetSendJSON(true)
	sendDict := map[string]string{"name": n}
	ret, err := c.SendReq(http.MethodPut, clashManagerAddr+"/proxies/"+p, sendDict)
	if err != nil {
		return err
	}
	utils.LogPrintInfo(ret.Text())
	return nil
}
