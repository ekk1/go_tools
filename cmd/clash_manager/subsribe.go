package main

import (
	"encoding/base64"
	"encoding/json"
	"go_utils/utils"
	"go_utils/utils/myhttp"
	"net/http"
	"slices"
	"time"
)

const (
	ssPrefix      = "SS_"
	customNodeKey = "CC_ExtraNodes"
)

var (
	CustomNodes map[string]*Trojan = make(map[string]*Trojan)

	SelectedExtraNode string = ""
)

type Subscribe struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	LastUpdated string `json:"last"`
	Content     string `json:"content"`
}

func (s *Subscribe) Update(useProxy bool) error {
	c := myhttp.NewHTTPClient()
	if useProxy {
		if err := c.SetProxy("http://127.0.0.1:7890"); err != nil {
			return err
		}
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
	slices.Sort(ssKeys)
	slices.Reverse(ssKeys)
	ret := []*Subscribe{}
	for _, v := range ssKeys {
		ss := LoadSingleSubscribe(v)
		if ss != nil {
			ret = append(ret, ss)
		}
	}
	return ret
}

func AddCustomNode(name string, t *Trojan) error {
	CustomNodes[name] = t
	return SaveCustomNodes()
}

func DeleteCustomNode(name string) error {
	if _, ok := CustomNodes[name]; ok {
		delete(CustomNodes, name)
	}
	return SaveCustomNodes()
}

func SaveCustomNodes() error {
	customNodesBytes, err := json.Marshal(CustomNodes)
	if err != nil {
		return err
	}
	customNodesString := base64.StdEncoding.EncodeToString(customNodesBytes)
	if err := kv.Set(customNodeKey, string(customNodesString)); err != nil {
		return nil
	}
	return kv.Save()
}

func LoadCustomNodes() error {
	if !kv.Exists(customNodeKey) {
		return nil
	}
	customNodesString := kv.Get(customNodeKey)
	customNodesBytes, err := base64.StdEncoding.DecodeString(customNodesString)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(customNodesBytes, &CustomNodes); err != nil {
		return err
	}
	return nil
}
