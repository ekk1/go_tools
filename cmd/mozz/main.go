package main

import (
	"encoding/json"
	"flag"
	"go_utils/utils"
	"go_utils/utils/minikv"
	"net/http"
	"path"
)

const (
	AccountPrefix = "ACC::"
	RecordPrefix  = "REC::"
)

var (
	kv      *minikv.KV
	baseDir string

	Accounts []*Account = []*Account{}
)

func main() {
	var (
		verboseFlag = flag.Int("v", 0, "debug (max 4)")
		listenAddr  = flag.String("l", "127.0.0.1", "listen address")
		listenPort  = flag.String("p", "5011", "listen port")
		baseDirFlag = flag.String("d", ".", "base dir")
	)
	flag.Parse()
	baseDir = *baseDirFlag

	kv = minikv.MustNewKV(path.Join(baseDir, "index"), 0)
	kv.MustLoad()

	utils.SetLogLevelByVerboseFlag(*verboseFlag)

	LoadAccounts()

	// s := myhttp.NewServer("mozz", *listenAddr, *listenPort)
	mux := http.NewServeMux()
	mux.HandleFunc("/accounts", handleAccounts)
	// ("/accounts", handleAccounts)
	// s.Serve()
	utils.LogPrintError(http.ListenAndServe(*listenAddr+":"+*listenPort, mux))
}

func LoadAccounts() {
	a := kv.Keys(AccountPrefix)
	for _, k := range a {
		accString := kv.Get(k)
		acc := &Account{}
		json.Unmarshal([]byte(accString), acc)
		acc.Balance = acc.InitialBalance
		Accounts = append(Accounts, acc)
	}
}

func LoadRecords() {
	a := kv.Keys(RecordPrefix)
	for _, k := range a {
		recString := kv.Get(k)
		rec := &Record{}
		json.Unmarshal([]byte(recString), rec)
		for _, acc := range Accounts {
			if acc.Name == rec.FromAccount {
				acc.Records = append(acc.Records, rec)
				acc.Balance -= rec.Amount
			}
		}
	}
}
