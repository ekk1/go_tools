package redis

import (
	"bufio"
	"go_utils/utils"
	"net"
	"sync"
)

// RedisClient is a safe for concurrent use redis client, support some of the commands
type RedisClient struct {
	Address string
	Port    string
	conn    *net.TCPConn
	rdBuf   *bufio.Reader
	mu      sync.Mutex
}

type RedisNotConnectedError struct{}

func (e *RedisNotConnectedError) Error() string {
	return "Redis is not connected, try DialRedis() first"
}

type RedisConnectError struct{}

func (e *RedisConnectError) Error() string {
	return "Can not connect to redis, check config and network connection"
}

// GenerateRedisClient returns a RedisClient, which needs to make a connection first
func GenerateRedisClient(addr, port string) *RedisClient {
	return &RedisClient{Address: addr, Port: port}
}

func (r *RedisClient) DialRedis() error {
	if tcpAddr, err := net.ResolveTCPAddr("tcp", r.Address+":"+r.Port); err == nil {
		if conn, err := net.DialTCP("tcp", nil, tcpAddr); err == nil {
			r.conn = conn
			r.rdBuf = bufio.NewReader(r.conn)
		} else {
			utils.LogPrintError("Failed to dial redis...", err.Error())
			r.conn = nil
			return &RedisConnectError{}
		}
	} else {
		utils.LogPrintError("Failed to resolve redis address...", err.Error())
		r.conn = nil
		return &RedisConnectError{}
	}
	return nil
}

func (r *RedisClient) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.conn.Close()
}

// General commands
func (r *RedisClient) Ping() (string, error) {
	return r.SimpleCommandSingleReturn("ping")
}

func (r *RedisClient) Set(k, v string) (string, error) {
	return r.SimpleCommandSingleReturn("SET", k, v)
}

func (r *RedisClient) Get(k string) (string, error) {
	return r.SimpleCommandSingleReturn("GET", k)
}

func (r *RedisClient) Keys(k string) ([]string, error) {
	return r.SimpleCommandMultiReturn("KEYS", k)
}

func (r *RedisClient) Exists(k string) (string, error) {
	return r.SimpleCommandSingleReturn("EXISTS", k)
}

func (r *RedisClient) Del(k string) (string, error) {
	return r.SimpleCommandSingleReturn("DEL", k)
}

func (r *RedisClient) FlushAll() (string, error) {
	return r.SimpleCommandSingleReturn("FLUSHALL")
}

func (r *RedisClient) Incr(k string) (string, error) {
	return r.SimpleCommandSingleReturn("INCR", k)
}

func (r *RedisClient) Decr(k string) (string, error) {
	return r.SimpleCommandSingleReturn("DECR", k)
}

// List commands
func (r *RedisClient) RPush(k, v string) (string, error) {
	return r.SimpleCommandSingleReturn("RPUSH", k, v)
}

func (r *RedisClient) LPOP(k string) (string, error) {
	return r.SimpleCommandSingleReturn("LPOP", k)
}

func (r *RedisClient) LLEN(k string) (string, error) {
	return r.SimpleCommandSingleReturn("LLEN", k)
}

// 暂时不处理嵌套 List 的情况，如果有的话应该自动变成一层的结构
func (r *RedisClient) LRange(k, start, end string) ([]string, error) {
	return r.SimpleCommandMultiReturn("LRANGE", k, start, end)
}

// Hash commands
func (r *RedisClient) HSet(k, f, v string) (string, error) {
	return r.SimpleCommandSingleReturn("HSET", k, f, v)
}

func (r *RedisClient) HGet(k, f string) (string, error) {
	return r.SimpleCommandSingleReturn("HGET", k, f)
}

func (r *RedisClient) HKeys(k string) ([]string, error) {
	return r.SimpleCommandMultiReturn("HKEYS", k)
}
