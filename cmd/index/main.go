package main

import (
	"flag"
	"go_utils/utils"
	"go_utils/utils/minikv"
	"go_utils/utils/myhttp"
	"path"
	"strings"
)

var (
	kv      *minikv.KV
	baseDir string
)

func main() {
	var (
		verboseFlag = flag.Int("v", 0, "debug (max 4)")
		listenAddr  = flag.String("l", "127.0.0.1:1111", "listen address")
		baseDirFlag = flag.String("d", ".", "base dir")
	)

	flag.Parse()

	kv = minikv.MustNewKV(path.Join(baseDir, "index"), 0)
	kv.MustLoad()

	baseDir = *baseDirFlag
	utils.SetLogLevelByVerboseFlag(*verboseFlag)

	listenVars := strings.Split(*listenAddr, ":")
	if len(listenVars) < 2 {
		utils.LogPrintError("Failed to parse listen address:", listenAddr)
		panic("Failde")
	}
	ii := myhttp.NewServer("index", listenVars[0], listenVars[1])
	ii.AddPost("/get", handleLoad)
	ii.AddPost("/set", handleSave)
	ii.AddPost("/rebalance", handleRebalance)
	ii.AddPost("/scrub", handleScrub)
	ii.AddPost("/stat", handleStat)
	ii.AddPost("/debug", handleDebug)

	myhttp.RunServers(ii)
}
