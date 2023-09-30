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
	tDelta := time.Now().Sub(t1)
	t.Log("Redis Set bench: ", float64(redisTestBatch)/float64(tDelta.Seconds()), "IOPS")

	t1 = time.Now()
	for i := 0; i < redisTestBatch; i++ {
		if _, err := r.Get("test"); err != nil {
			t.Fatal(err)
		}
	}
	tDelta = time.Now().Sub(t1)
	t.Log("Redis Get bench: ", float64(redisTestBatch)/float64(tDelta.Seconds()), "IOPS")

	t.Log(r.FlushAll())
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

	t.Log("Redis Set: ")
	t.Log(r.Set("test", "testout"))

	t.Log("Redis Get: ")
	t.Log(r.Get("test"))

	t.Log("Redis Set empty: ")
	t.Log(r.Set("test", ""))

	t.Log("Redis Get empty: ")
	t.Log(r.Get("test"))

	t.Log("Redis RPush: ")
	t.Log(r.RPush("test_list", "test"))
	t.Log(r.RPush("test_list", "test"))

	t.Log("Redis LLEN: ")
	t.Log(r.LLEN("test_list"))

	t.Log("Redis LPOP: ")
	t.Log(r.LPOP("test_list"))

	t.Log("Redis LRange: ")
	t.Log(r.LRange("test_list", "0", "-1"))

	t.Log("Redis HSet: ")
	t.Log(r.HSet("test_hash", "ffff", "test1"))

	t.Log("Redis HGet: ")
	t.Log(r.HGet("test_hash", "ffff"))

	t.Log("Redis HKeys: ")
	t.Log(r.HKeys("test_hash"))

	t.Log("Redis Keys: ")
	t.Log(r.Keys("*"))

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
