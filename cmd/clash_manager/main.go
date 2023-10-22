package main

import (
	"flag"
	"go_utils/utils"
	"go_utils/utils/minikv"
	"go_utils/utils/myhttp"
)

var kv *minikv.KV
var clashYamlOutPath string

func main() {
	var verboseFlag = flag.Int("v", 0, "debug (max 4)")
	var listenAddr = flag.String("l", "127.0.0.1", "listen address")
	var listenPort = flag.String("p", "8888", "listen port")
	var yamlPath = flag.String("o", "./config.yml", "clash out yaml")
	flag.Parse()

	clashYamlOutPath = *yamlPath

	kv = minikv.MustNewKV("ss", 0)
	kv.MustLoad()

	LoadClashRules()

	utils.SetLogLevelByVerboseFlag(*verboseFlag)

	ss := myhttp.NewServer("cc", *listenAddr, *listenPort)
	ss.AddRoute("/", handleRoot)
	ss.AddRoute("/delete", handleDelete)
	ss.AddRoute("/update", handleUpdete)
	ss.AddRoute("/selectproxy", handleSelectProxy)
	ss.AddRoute("/selectnode", handleSelectNode)
	ss.AddRoute("/rules", handleAddRules)
	ss.AddRoute("/deleterule", handleDeleteRules)
	ss.AddRoute("/gen", handleGenerateYaml)

	myhttp.RunServers(ss)
}
