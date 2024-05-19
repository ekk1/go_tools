package main

import (
	"encoding/json"
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

func LoadAllKeys() error {
	keyData, err := os.ReadFile("key.txt")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			utils.LogPrintWarning("No key file exists, creating...")
			initKeys := []*APIKey{{
				Name: "test",
				Key:  "",
			}}
			writeData, errJson := json.Marshal(initKeys)
			if errJson != nil {
				return errJson
			}
			if errWr := os.WriteFile("key.txt", writeData, 0400); errWr != nil {
				return errWr

			}
			return nil
		}
		return err
	}

	if err := json.Unmarshal(keyData, &KEYS); err != nil {
		return err
	}

	return nil
}
