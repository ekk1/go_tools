package main

import (
	"go_utils/utils"
	"go_utils/utils/redis"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const INJECT_NUM = 1000000

func main() {
	r := redis.GenerateRedisClientWithAuth("127.0.0.1", "7000", "123")
	if len(os.Args) > 2 {
		utils.LogPrintInfo("Running injection")
		r.Cluster = true
		nodes, err := r.ClusterNodes()
		if err != nil {
			utils.LogPrintError(err)
			os.Exit(1)
		}
		utils.LogPrintInfo(nodes)
		wg := sync.WaitGroup{}
		for _, n := range redis.ParseNodesAddr(nodes, true) {
			wg.Add(1)
			fields := strings.Split(n, ":")
			host := fields[0]
			port := fields[1]
			utils.LogPrintInfo("Staring:", host, port)
			go func() {
				rr := redis.GenerateRedisClientWithAuth(host, port, "123")
				t1 := time.Now()
				for i := 0; i < INJECT_NUM; i++ {
					data := strconv.Itoa(i)
					//rr.Set(data, data)
					rr.Del(data)
				}
				utils.LogPrintInfo(host, port, "finished, speed: ", float64(INJECT_NUM)/float64(time.Now().Sub(t1).Seconds()))
				wg.Done()
			}()
		}
		wg.Wait()
		utils.LogPrintInfo("Done injection")
		os.Exit(0)
	}
}
