package main

import (
	"flag"
	"go_utils/utils"
	"go_utils/utils/minikv"
	"go_utils/utils/myhttp"
	"runtime"
)

func main() {
	var verboseFlag = flag.Int("v", 0, "debug (max 4)")
	var listenAddr = flag.String("l", "127.0.0.1", "listen address")
	var listenPort = flag.String("p", "8888", "listen port")
	var yamlPath = flag.String("o", "config.yaml", "clash out yaml")
	var clashBinaryFlag = flag.String("b", "clash", "path to clash core")
	var clashConfigFlag = flag.String("c", "clashc", "path to clash config dir")
	flag.Parse()

	if runtime.GOOS == "windows" {
		clashYamlOutPath = ".\\" + *clashConfigFlag + "\\" + *yamlPath
	} else if runtime.GOOS == "linux" {
		clashYamlOutPath = "./" + *clashConfigFlag + "/" + *yamlPath
	}
	clashBinary = *clashBinaryFlag
	clashConfigDir = *clashConfigFlag

	utils.LogPrintInfo("Yaml path: ", clashYamlOutPath)

	kv = minikv.MustNewKV("ss", 0)
	kv.MustLoad()

	LoadClashRules()
	LoadClashProxies()

	utils.SetLogLevelByVerboseFlag(*verboseFlag)

	ss := myhttp.NewServer("cc", *listenAddr, *listenPort)
	ss.AddGet("/", handleIndex)
	ss.AddREST("/rules", handleRules)
	ss.AddGet("/rules/delete", handleRules)
	ss.AddREST("/subs", handleSubs)
	ss.AddGet("/subs/delete", handleSubs)
	ss.AddGet("/subs/update", handleSubs)
	ss.AddREST("/proxies", handleProxies)
	ss.AddGet("/proxies/delete", handleProxies)

	myhttp.RunServers(ss)
}
