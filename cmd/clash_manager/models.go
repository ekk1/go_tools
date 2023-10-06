package main

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"go_utils/utils"
	"net/http"
	"time"
)

const ssPrefix = "SS_"

type Subscribe struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	LastUpdated string `json:"last"`
	Content     string `json:"content"`
}

func (s *Subscribe) Update() error {
	c := utils.NewHTTPClient()
	if err := c.SetProxy("http://127.0.0.1:7890"); err != nil {
		return err
	}
	resp, err := c.SendReq(http.MethodGet, s.URL, nil)
	if err != nil {
		return err
	}
	s.Content = resp.Text()
	s.LastUpdated = time.Now().Format(time.RFC1123Z)
	return nil
}

func (s *Subscribe) Save() error {
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	jsonStr := base64.URLEncoding.EncodeToString(jsonBytes)
	kv.Set(ssPrefix+s.Name, jsonStr)
	return kv.Save()
}

func LoadSingleSubscribe(name string) *Subscribe {
	ssData := kv.Get(name)
	jsonBytes, err := base64.URLEncoding.DecodeString(ssData)
	if err != nil {
		utils.LogPrintWarning("Failed to decode ss data", ssData)
		return nil
	}
	ss := &Subscribe{}
	if err := json.Unmarshal(jsonBytes, ss); err != nil {
		utils.LogPrintWarning("Failed to decode json: ", jsonBytes)
		return nil
	}
	return ss
}

func LoadSubscribe() []*Subscribe {
	ssKeys := kv.Keys(ssPrefix)
	ret := []*Subscribe{}
	for _, v := range ssKeys {
		ss := LoadSingleSubscribe(v)
		if ss != nil {
			ret = append(ret, ss)
		}
	}
	return ret
}
