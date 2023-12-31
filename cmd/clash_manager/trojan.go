package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Trojan struct {
	Name     string
	Address  string
	Port     string
	Password string
	Sni      string
}

func NewFlowerTrojan(input string) (*Trojan, error) {
	fields := strings.Split(input, ", ")
	if len(fields) < 5 {
		return nil, errors.New("Failed to decode " + input)
	}
	return &Trojan{
		Name:     strings.Split(fields[0], " =")[0],
		Address:  fields[1],
		Port:     fields[2],
		Password: strings.Split(fields[3], "=")[1],
		Sni:      strings.Split(fields[4], "=")[1],
	}, nil
}

func NewXianyuTrojan(input string) (*Trojan, error) {
	info := strings.Split(input, "?")
	basicInfo := strings.Split(info[0], "@")
	extraInfo := strings.Split(info[1], "#")
	decodedName, err := url.QueryUnescape(extraInfo[1])
	sniData := ""
	if err != nil {
		return nil, err
	}
	sniInfos := strings.Split(extraInfo[0], "&")
	for _, v := range sniInfos {
		if strings.Contains(v, "sni") {
			sniData = strings.Split(v, "=")[1]
		}
	}
	return &Trojan{
		Name:     decodedName,
		Password: strings.Split(basicInfo[0], "://")[1],
		Address:  strings.Split(basicInfo[1], ":")[0],
		Port:     strings.Split(basicInfo[1], ":")[1],
		Sni:      sniData,
	}, nil
}

func (t *Trojan) Render() string {
	ret := "\n"
	ret += "  - name: \"%s\"\n"
	ret += "    type: trojan\n"
	ret += "    server: %s\n"
	ret += "    port: %s\n"
	ret += "    password: %s\n"
	ret += "    sni: %s\n"
	ret += "    skip-cert-verify: true\n"

	return fmt.Sprintf(
		ret,
		t.Name, t.Address, t.Port,
		t.Password, t.Sni,
	)
}
