package main

import (
	"flag"
	"go_utils/utils"
	"go_utils/utils/minikv"
	"go_utils/utils/myhttp"
)

var kv *minikv.KV

func main() {
	var verboseFlag = flag.Int("v", 0, "debug (max 4)")
	var listenAddr = flag.String("l", "127.0.0.1", "listen address")
	var listenPort = flag.String("p", "8888", "listen port")
	flag.Parse()

	kv = minikv.MustNewKV("ss", 0)
	kv.MustLoad()

	utils.SetLogLevelByVerboseFlag(*verboseFlag)

	ss := myhttp.NewServer("cc", *listenAddr, *listenPort)
	ss.AddRoute("/", handleRoot)
	ss.AddRoute("/delete", handleDelete)
	ss.AddRoute("/update", handleUpdete)
	ss.AddRoute("/selectproxy", handleSelectProxy)
	ss.AddRoute("/selectnode", handleSelectNode)

	myhttp.RunServers(ss)
}
