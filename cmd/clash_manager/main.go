package main

import (
	"flag"
	"go_utils/utils"
	"go_utils/utils/minikv"
	"go_utils/utils/myhttp"
)

var kv *minikv.KV
var clashYamlOutPath string

var clashBinary string
var clashConfigDir string

func main() {
	var verboseFlag = flag.Int("v", 0, "debug (max 4)")
	var listenAddr = flag.String("l", "127.0.0.1", "listen address")
	var listenPort = flag.String("p", "8888", "listen port")
	var yamlPath = flag.String("o", "./config.yml", "clash out yaml")
	var clashBinaryFlag = flag.String("b", "./clash", "path to clash core")
	var clashConfigFlag = flag.String("c", "clashc", "path to clash config dir")
	flag.Parse()

	clashYamlOutPath = *yamlPath
	clashBinary = *clashBinaryFlag
	clashConfigDir = *clashConfigFlag

	kv = minikv.MustNewKV("ss", 0)
	kv.MustLoad()

	LoadClashRules()

	utils.SetLogLevelByVerboseFlag(*verboseFlag)

	ss := myhttp.NewServer("cc", *listenAddr, *listenPort)
	ss.AddGet("/", handleRoot)
	ss.AddPost("/add", handleAddSub)
	ss.AddPost("/rules", handleAddRules)
	ss.AddGet("/delete", handleDelete)
	ss.AddGet("/update", handleUpdete)
	ss.AddGet("/selectproxy", handleSelectProxy)
	ss.AddGet("/selectnode", handleSelectNode)
	ss.AddGet("/deleterule", handleDeleteRules)
	ss.AddGet("/gen", handleGenerateYaml)
	ss.AddGet("/start", handleRunCmd)
	ss.AddGet("/stop", handleTerminate)

	myhttp.RunServers(ss)
}
