package main

import (
	"errors"
	"go_utils/utils"
	"io/fs"
	"os"
)

type APIKey struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Org      string `json:"org"`
	Endpoint string `json:"endpoint"`
}

// def load_all_keys():
//     """load all keys from db file"""
//     with open("key.txt", encoding="utf8") as f:
//         for k in json.loads(f.read()):
//             key_list.append(k)

func LoadAllKeys() error {
	keyData, err := os.ReadFile("key.txt")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			utils.LogPrintWarning("No key file exists, creating...")
			return nil
		}
		return err
	}

}
