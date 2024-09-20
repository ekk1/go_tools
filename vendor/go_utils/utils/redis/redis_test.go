package redis

import (
	"go_utils/utils"
	"sync"
	"testing"
	"time"
)

func TestRedisBenchmark(t *testing.T) {
	redisTestBatch := 50000
	r := GenerateRedisClient("127.0.0.1", "6379")
	if err := r.DialRedis(); err != nil {
		t.Fatal(err)
	}
	utils.LogLevel = utils.LOG_LEVEL_INFO
	t1 := time.Now()
	for i := 0; i < redisTestBatch; i++ {
		if _, err := r.Set("test", "test"); err != nil {
			t.Fatal(err)
		}
	}
	tDelta := time.Since(t1)
	t.Log("Redis Set bench: ", float64(redisTestBatch)/float64(tDelta.Seconds()), "IOPS")

	t1 = time.Now()
	for i := 0; i < redisTestBatch; i++ {
		if _, err := r.Get("test"); err != nil {
			t.Fatal(err)
		}
	}
	tDelta = time.Since(t1)
	t.Log("Redis Get bench: ", float64(redisTestBatch)/float64(tDelta.Seconds()), "IOPS")

	t.Log(r.FlushAll())
}

func TestRedisConnection(t *testing.T) {
	// r := GenerateRedisClient("127.0.0.1", "6379")
	r := GenerateRedisClientWithAuth("127.0.0.1", "6379", "123")
	if err := r.DialRedis(); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Select("1"); err != nil {
		t.Fatal(err)
	}
	utils.LogLevel = utils.LOG_LEVEL_DEBUG4
	if _, err := r.Set("test", "test"); err != nil {
		t.Fatal(err)
	}
	if ret, err := r.Get("test"); err != nil {
		t.Fatal(err)
	} else {
		utils.LogPrintInfo(ret)
	}

	utils.LogPrintInfo("Turn off redis")
	time.Sleep(5 * time.Second)

	if _, err := r.Expire("test", 3); err != nil {
		t.Error(err)
	}
	time.Sleep(1 * time.Second)
	if ret, err := r.TTL("test"); err != nil {
		t.Error(err)
	} else {
		utils.LogPrintInfo(ret)
	}
	for li := 0; li < 20; li++ {
		time.Sleep(1 * time.Second)
		if ret, err := r.Get("test"); err != nil {
			t.Error(err)
		} else {
			utils.LogPrintInfo(ret)
		}
	}

	utils.LogPrintInfo("Turn on redis")
	time.Sleep(5 * time.Second)

	if ret, err := r.Get("test"); err != nil {
		t.Fatal(err)
	} else {
		utils.LogPrintInfo(ret)
	}

	if _, err := r.Expire("test", 3); err != nil {
		t.Fatal(err)
	}

	if ret, err := r.TTL("test"); err != nil {
		t.Fatal(err)
	} else {
		utils.LogPrintInfo(ret)
	}
	time.Sleep(1 * time.Second)

	if ret, err := r.TTL("test"); err != nil {
		t.Fatal(err)
	} else {
		utils.LogPrintInfo(ret)
	}
	time.Sleep(2 * time.Second)
	if ret, err := r.TTL("test"); err != nil {
		t.Fatal(err)
	} else {
		utils.LogPrintInfo(ret)
	}
	// t.Log(r.FlushAll())
}

func TestRedisCmd(t *testing.T) {
	r := GenerateRedisClient("127.0.0.1", "6379")
	if err := r.DialRedis(); err != nil {
		t.Fatal(err)
	}

	utils.LogLevel = utils.LOG_LEVEL_INFO

	t.Log("Redis FlushAll: ")
	t.Log(r.FlushAll())

	t.Log("Redis Ping: ")
	t.Log(r.Ping())

	t.Log("Redis Keys: ")
	keys, err := r.Keys("*")
	t.Log(keys, err, len(keys))

	t.Log("Redis Set: ")
	t.Log(r.Set("test", "testout"))

	t.Log("Redis Get: ")
	t.Log(r.Get("test"))

	t.Log("Redis Set empty: ")
	t.Log(r.Set("test", ""))

	t.Log("Redis Get empty: ")
	t.Log(r.Get("test"))

	t.Log("Redis LRange: ")
	values, err := r.LRange("test_list", "0", "-1")
	t.Log(values, err, len(values))

	t.Log("Redis RPush: ")
	t.Log(r.RPush("test_list", "test"))
	t.Log(r.RPush("test_list", "test"))

	t.Log("Redis LLEN: ")
	t.Log(r.LLEN("test_list"))

	t.Log("Redis LRange: ")
	values, err = r.LRange("test_list", "0", "-1")
	t.Log(values, err, len(values))

	t.Log("Redis LPOP: ")
	t.Log(r.LPOP("test_list"))

	t.Log("Redis LRange: ")
	values, err = r.LRange("test_list", "0", "-1")
	t.Log(values, err, len(values))

	t.Log("Redis HSet: ")
	t.Log(r.HSet("test_hash", "ffff", "test1"))

	t.Log("Redis HGet: ")
	t.Log(r.HGet("test_hash", "ffff"))

	t.Log("Redis HKeys: ")
	t.Log(r.HKeys("test_hash"))

	t.Log("Redis Keys: ")
	keys, err = r.Keys("*")
	t.Log(keys, err, len(keys))

	t.Log("Redis Exists: ")
	t.Log(r.Exists("test_hash"))

	t.Log("Redis Del: ")
	t.Log(r.Del("test_hash"))

	t.Log("Redis Exists after del: ")
	t.Log(r.Exists("test_hash"))

	t.Log(r.FlushAll())
}

func TestRedisClientRace(t *testing.T) {
	r := GenerateRedisClient("127.0.0.1", "6379")
	if err := r.DialRedis(); err != nil {
		t.Fatal(err)
	}
	utils.LogLevel = utils.LOG_LEVEL_INFO
	r.FlushAll()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Incr("testcounter")
		}()
	}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Decr("testcounter")
		}()
	}
	wg.Wait()

	t.Log(r.Get("testcounter"))
	t.Log(r.FlushAll())
}

func TestGenerateCommandBytesFromRawCommand(t *testing.T) {
	ret := generateCommandBytesFromRawCommand("PING")
	t.Log(ret)
}

func TestSingle(t *testing.T) {
	r := GenerateRedisClientWithAuth("127.0.0.1", "7000", "123")
	r.Cluster = true
	ret, err := r.ClusterNodes()
	if err != nil {
		t.Fatal(err)
	}
	utils.LogPrintInfo(ParseNodesAddr(ret, false))
	utils.LogPrintInfo(ParseNodesAddr(ret, true))
}
