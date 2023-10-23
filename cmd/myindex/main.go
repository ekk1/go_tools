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
	var listenPort = flag.String("p", "9900", "listen port")
	flag.Parse()

	kv = minikv.MustNewKV("ss", 0)
	kv.MustLoad()

	utils.SetLogLevelByVerboseFlag(*verboseFlag)

	ii := myhttp.NewServer("index", *listenAddr, *listenPort)
	ii.AddGet("/", handleRoot)
	ii.AddPost("/add", handleAdd)
	ii.AddPost("/delete", handleDelete)
	ii.AddPost("/move", handleMove)

	myhttp.RunServers(ii)
}
