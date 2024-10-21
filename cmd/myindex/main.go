package main

import (
	"flag"
	"go_utils/utils"
	"go_utils/utils/minikv"
	"go_utils/utils/myhttp"
	"time"
)

var kv *minikv.KV

var cookieCache = map[string]int64{}
var authDB = map[string]string{}

func cookieCheck(value string) (string, bool) {
	if ts, ok := cookieCache[value]; !ok {
		return "Not found", false
	} else {
		tsTime := time.Unix(ts, 0)
		if time.Now().Sub(tsTime).Seconds() > 3600 {
			return "Expired", false
		}
	}
	return "ok", true
}

func main() {
	var verboseFlag = flag.Int("v", 0, "debug (max 4)")
	var listenAddr = flag.String("l", "127.0.0.1", "listen address")
	var listenPort = flag.String("p", "9900", "listen port")
	flag.Parse()

	kv = minikv.MustNewKV("ss", 0)
	kv.MustLoad()

	authDB["admin"] = "3c9909afec25354d551dae21590bb26e38d53f2173b8d3dc3eee4c047e7ab1c1eb8b85103e3be7ba613b31bb5c9c36214dc9f14a42fd7a2fdb84856bca5c44c2"

	utils.SetLogLevelByVerboseFlag(*verboseFlag)

	ii := myhttp.NewServer("index", *listenAddr, *listenPort)
	ii.AddGet("/", myhttp.CookieChecker("authToken", "/loginpage", cookieCheck, handleRoot))
	ii.AddGet("/loginpage", handleLoginPage)
	ii.AddPost("/login", handleLogin)
	ii.AddPost("/add", myhttp.CookieChecker("authToken", "/loginpage", cookieCheck, handleAdd))
	ii.AddPost("/delete", myhttp.CookieChecker("authToken", "/loginpage", cookieCheck, handleDelete))
	ii.AddPost("/move", myhttp.CookieChecker("authToken", "/loginpage", cookieCheck, handleMove))

	myhttp.RunServers(ii)
}
