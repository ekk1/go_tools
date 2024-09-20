package redis

import (
	"bufio"
	"errors"
	"go_utils/utils"
	"net"
	"strconv"
	"strings"
	"sync"
)

// RedisClient is a safe for concurrent use redis client, support some of the commands
type RedisClient struct {
	Address  string
	Port     string
	Password string
	DB       string
	conn     *net.TCPConn
	rdBuf    *bufio.Reader
	mu       sync.Mutex
	Cluster  bool
	// ClusterNodeMap map[string]{int, int}
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

// GenerateRedisClient returns a RedisClient, which needs to make a connection first
func GenerateRedisClientWithAuth(addr, port, password string) *RedisClient {
	return &RedisClient{Address: addr, Port: port, Password: password}
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
	if r.Password != "" {
		utils.LogPrintDebug2("Authenticating redis connection")
		_, err := r.AuthNoLock(r.Password)
		if err != nil {
			return err
		}
	}
	if r.DB != "" {
		utils.LogPrintDebug2("Selecting DB", r.DB)
		_, err := r.SelectNoLock(r.DB)
		if err != nil {
			return err
		}
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

func (r *RedisClient) Auth(p string) (string, error) {
	return r.SimpleCommandSingleReturn("AUTH", p)
}

func (r *RedisClient) AuthNoLock(p string) (string, error) {
	return r.SimpleCommandSingleReturnNoLock("AUTH", p)
}

func (r *RedisClient) Select(d string) (string, error) {
	if r.Cluster {
		return "", errors.New("Can not select when using cluster")
	}
	if dbNum, err := strconv.Atoi(d); err != nil {
		return "", err
	} else {
		if dbNum < 0 || dbNum > 16 {
			return "", errors.New("DB index out of range: " + d)
		}
	}
	r.DB = d
	return r.SimpleCommandSingleReturn("SELECT", d)
}

func (r *RedisClient) SelectNoLock(d string) (string, error) {
	if r.Cluster {
		return "", errors.New("Can not select when using cluster")
	}
	if dbNum, err := strconv.Atoi(d); err != nil {
		return "", err
	} else {
		if dbNum < 0 || dbNum > 16 {
			return "", errors.New("DB index out of range: " + d)
		}
	}
	r.DB = d
	return r.SimpleCommandSingleReturnNoLock("SELECT", d)
}

// time is seconds
func (r *RedisClient) Expire(k string, time int) (string, error) {
	return r.SimpleCommandSingleReturn("EXPIRE", k, strconv.Itoa(time))
}

func (r *RedisClient) TTL(k string) (string, error) {
	return r.SimpleCommandSingleReturn("TTL", k)
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

func (r *RedisClient) ClusterNodes() ([]string, error) {
	nodes := []string{}
	if !r.Cluster {
		return nil, errors.New("This is not a cluster")
	}
	ret, err := r.SimpleCommandSingleReturn("CLUSTER", "NODES")
	if err != nil {
		return nil, err
	}
	nodesStr := strings.Split(ret, "\n")
	for _, v := range nodesStr {
		if len(v) < 1 {
			continue
		}
		nodes = append(nodes, v)
	}
	return nodes, err
}

func ParseNodesAddr(nodes []string, onlyMaster bool) []string {
	ret := []string{}

	for _, v := range nodes {
		fields := strings.Split(v, " ")
		if len(fields) < 8 {
			continue
		}
		if onlyMaster {
			if strings.Contains(fields[2], "master") {
				ret = append(ret, strings.Split(fields[1], "@")[0])
			}
		} else {
			ret = append(ret, strings.Split(fields[1], "@")[0])
		}
	}
	return ret
}
