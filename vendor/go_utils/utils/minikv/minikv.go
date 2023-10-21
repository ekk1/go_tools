package minikv

import (
	"encoding/json"
	"errors"
	"go_utils/utils"
	"os"
	"strings"
	"sync"
)

/*
TODO:
1. Add LRU and limit for kv
2. Encrypt data when saving to disk
3. Count hit num
*/

type KV struct {
	name       string
	counter    map[string]int64
	kvs        map[string]string
	lastSaved  int64
	lastSaveOK bool
	limit      int64
	current    int64
	key        []byte
	lock       sync.Mutex
}

func NewKV(name string, limit int64) (*KV, error) {
	var newLimit int64 = 0
	if limit <= 0 {
		newLimit = 128 * 1024 * 1024
	} else if limit > 0 && limit < 1*1024*1024 {
		return nil, errors.New("Limit should be larger than 1MB")
	} else {
		newLimit = limit
	}
	kv := &KV{
		name:  name,
		kvs:   make(map[string]string),
		limit: newLimit,
		lock:  sync.Mutex{},
	}
	return kv, nil
}

func MustNewKV(name string, limit int64) *KV {
	kv, err := NewKV(name, limit)
	utils.ErrExit(err, 1)
	return kv
}

func (kv *KV) Get(key string) string {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	if v, ok := kv.kvs[key]; ok {
		return v
	}
	return ""
}

func (kv *KV) Set(key, value string) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	kv.kvs[key] = value
	return nil
}

func (kv *KV) Keys(keyword string) []string {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	retList := []string{}
	for k := range kv.kvs {
		if strings.Contains(k, keyword) {
			retList = append(retList, k)
		}
	}
	return retList
}

func (kv *KV) Exists(key string) bool {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	if _, ok := kv.kvs[key]; ok {
		return true
	}
	return false
}

func (kv *KV) Delete(key string) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	if _, ok := kv.kvs[key]; ok {
		delete(kv.kvs, key)
	}
	return nil
}

func (kv *KV) Save() error {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	fileName := kv.name + ".snapshot"
	jsonBytes, err := json.Marshal(kv.kvs)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, jsonBytes, 0600)
}

func (kv *KV) Load() error {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	fileName := kv.name + ".snapshot"
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &kv.kvs)
}

func (kv *KV) MustLoad() {
	err := kv.Load()
	if errors.Is(err, os.ErrNotExist) {
		utils.LogPrintWarning("DB not exists, creating...")
		err2 := kv.Save()
		utils.ErrExit(err2, 2)
		return
	}
	utils.ErrExit(err, 2)
}

func (kv *KV) DumpJSON() (string, error) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	jsonBytes, err := json.Marshal(kv.kvs)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
