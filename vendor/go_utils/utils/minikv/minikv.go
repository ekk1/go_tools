package minikv

import (
	"encoding/json"
	"errors"
	"go_utils/utils"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
TODO:
1. Add LRU and limit for kv
2. Encrypt data when saving to disk
3. Count hit num
*/

const MAX_DB_SAVES = 3

type KV struct {
	name string
	// counter    map[string]int64
	kvs map[string]string
	// lastSaved  int64
	// lastSaveOK bool
	limit int64
	// current    int64
	// key        []byte
	lock sync.Mutex
}

func NewKV(name string, limit int64) (*KV, error) {
	var newLimit int64 = 0
	if limit <= 0 {
		newLimit = 128 * 1024 * 1024
	} else if limit > 0 && limit < 1*1024*1024 {
		return nil, errors.New("limit should be larger than 1MB")
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
	utils.ErrExit(err)
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
	// if _, ok := kv.kvs[key]; ok {
	delete(kv.kvs, key)
	// }
	return nil
}

func getDBFiles(kvName string) ([]string, error) {
	pwdItems, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	foundDBFiles := []string{}

	for _, f := range pwdItems {
		if f.IsDir() {
			continue
		}
		if m, err := regexp.MatchString(
			kvName+".snapshot."+"\\d+",
			f.Name(),
		); err != nil {
			return nil, err
		} else {
			if m {
				foundDBFiles = append(foundDBFiles, f.Name())
			}
		}
	}

	slices.Sort(foundDBFiles)
	return foundDBFiles, nil
}

func (kv *KV) Save() error {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	t1 := strconv.FormatInt(
		time.Now().UnixMilli(),
		10,
	)
	fileName := kv.name + ".snapshot." + t1
	jsonBytes, err := json.Marshal(kv.kvs)
	if err != nil {
		return err
	}
	if err := os.WriteFile(fileName, jsonBytes, 0600); err != nil {
		return err
	}

	if dbFiles, err := getDBFiles(kv.name); err != nil {
		return err
	} else {
		if len(dbFiles) > MAX_DB_SAVES {
			for _, v := range dbFiles[0 : len(dbFiles)-MAX_DB_SAVES] {
				if delErr := os.Remove(v); delErr != nil {
					return delErr
				}
			}
		}
	}

	return nil
}

func (kv *KV) Load() error {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	dbFiles, err := getDBFiles(kv.name)
	if err != nil {
		return err
	}

	if len(dbFiles) < 1 {
		utils.LogPrintWarning("no db file found")
		return nil
	}

	fileName := dbFiles[len(dbFiles)-1]
	utils.LogPrintDebug("Reading DB File:", fileName)
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
		utils.ErrExit(err2)
		return
	}
	utils.ErrExit(err)
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
