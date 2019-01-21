package cache

import (
	"errors"

	"github.com/xxlixin1993/easyGo/configure"
	redigo "github.com/gomodule/redigo/redis"
)

// Client redis handler
type Client struct {
	rc       *redigo.Pool
	name     string
	modeType uint8
}

// GetClient get redis-client
func GetClient(redisName string, modeType uint8) (*Client, error) {
	var rc *redigo.Pool
	var err error

	if modeType == configure.KReadMode {
		rc, err = GetSlaveConn(redisName)
	} else {
		rc, err = GetMasterConn(redisName)
	}

	if err != nil {
		return nil, err
	}

	if rc == nil {
		return nil, errors.New("[redis] the name given is not in the client pool")
	}
	return &Client{
		rc:       rc,
		name:     redisName,
		modeType: modeType,
	}, nil
}

// Del redis-del command
func (c *Client) Del(key string) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

// Get redis-get command
func (c *Client) Get(key string) (string, error) {
	conn := c.rc.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", key)
	return redigo.String(reply, err)
}

// MGetWithInt64 redis-mget command with returning valye as int64s
func (c *Client) MGetWithInt64(key ...interface{}) ([]int64, error) {
	conn := c.rc.Get()
	defer conn.Close()
	reply, err := conn.Do("MGET", key...)
	return redigo.Int64s(reply, err)
}

// GetWithUint64 redis-get command with returning value as uint64
func (c *Client) GetWithUint64(key string) (uint64, error) {
	conn := c.rc.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", key)
	return redigo.Uint64(reply, err)
}

// GetWithBytes redis-get command with returning value as byte
func (c *Client) GetWithBytes(key string) ([]byte, error) {
	conn := c.rc.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", key)
	return redigo.Bytes(reply, err)
}

// Set redis-set command
func (c *Client) Set(key string, value interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	return err
}

// SetEx redis-setex command
func (c *Client) SetEx(key string, expire int, value interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", key, expire, value)
	return err
}

// Incr redis-incr command
func (c *Client) Incr(key string, value interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("INCRBY", key, value)
	return err
}

// SRandMember redis-smember command
func (c *Client) SRandMember(key string, count interface{}) ([]int64, error) {
	conn := c.rc.Get()
	defer conn.Close()
	member, err := conn.Do("SRANDMEMBER", key, count)
	return redigo.Int64s(member, err)
}

// HMSet redis-hmset command
func (c *Client) HMSet(data ...interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("HMSET", data...)
	return err
}

// HSet redis-hset command
func (c *Client) HSet(key string, field string, value interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", key, field, value)
	return err
}

// HGet redis-hget command
func (c *Client) HGet(key string, field string) (string, error) {
	conn := c.rc.Get()
	defer conn.Close()
	reply, err := conn.Do("HGET", key, field)
	return redigo.String(reply, err)
}

// HGetAll redis-hgetall command
func (c *Client) HGetAll(key string) (map[string]int, error) {
	conn := c.rc.Get()
	defer conn.Close()
	reply, err := conn.Do("HGETALL", key)
	return redigo.IntMap(reply, err)
}

// HDel redis-hdel command
func (c *Client) HDel(key string, field string) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", key, field)
	return err
}

// HExists redis-hexists command
func (c *Client) HExists(key string, field string) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("HEXISTS", key, field)
	return err
}

// Expire redis-expire command
func (c *Client) Expire(key string, value interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("EXPIRE", key, value)
	return err
}

// ZRangeByScore redis-zrangebyyscore command
func (c *Client) ZRangeByScore(key string, low, high interface{}) ([]string, error) {
	conn := c.rc.Get()
	defer conn.Close()
	reply, err := conn.Do("ZRANGEBYSCORE", key, low, high)
	return redigo.Strings(reply, err)
}

// ZAdd redis-zadd command
// every score is always in front of its member
func (c *Client) ZAdd(key string, score, member interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("zadd", key, score, member)
	return err
}

// SMembers redis-smembers command
func (c *Client) SMembers(key string) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("SMEMBERS", key)
	return err
}

// SAdd redis-sadd command
func (c *Client) SAdd(data ...interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", data...)
	return err
}

// SRem redis-srem command
func (c *Client) SRem(data ...interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", data...)
	return err
}

// ZAddBatch redis-zadd command
// key must be the first param
// every score is always in front of its member
func (c *Client) ZAddBatch(keyScoreMember ...interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("zadd", keyScoreMember...)
	return err
}

// ZRem redis-zrem command
func (c *Client) ZRem(key string, member interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("ZREM", key, member)
	return err
}

// ZRemBatch redis-zrem command
// key must be the first param
func (c *Client) ZRemBatch(keyMember ...interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("ZREM", keyMember...)
	return err
}

// LLen redis-llen command
func (c *Client) LLen(key string) (int64, error) {
	conn := c.rc.Get()
	defer conn.Close()
	reply, err := conn.Do("LLEN", key)
	return redigo.Int64(reply, err)
}

// LPop redis-lpop command
func (c *Client) LPop(key string) (string, error) {
	conn := c.rc.Get()
	defer conn.Close()
	reply, err := conn.Do("LPOP", key)
	if err != nil || reply == nil {
		return "", nil
	}
	return redigo.String(reply, err)
}

// LPopWithBytes redis-lpop command with returning value as byte
func (c *Client) LPopWithBytes(key string) ([]byte, error) {
	conn := c.rc.Get()
	defer conn.Close()
	reply, err := conn.Do("LPOP", key)
	if err != nil || reply == nil {
		return []byte{}, nil
	}
	return redigo.Bytes(reply, err)
}

// RPush redis-rpush command
func (c *Client) RPush(key string, data interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", key, data)
	return err
}

// RPushBatch redis-rpush command
// key must be the first param
func (c *Client) RPushBatch(keyData ...interface{}) error {
	conn := c.rc.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", keyData...)
	return err
}

// GEORadius redis-georadius command
// key must be the first param
func (c *Client) GEORadius(keyData ...interface{}) ([]int64, error) {
	conn := c.rc.Get()
	defer conn.Close()
	reply, err := conn.Do("GEORADIUS", keyData...)
	return redigo.Int64s(reply, err)
}