package cache

import (
	"errors"
	"time"
	"strconv"
	"math/rand"

	"github.com/xxlixin1993/easyGo/configure"
	"github.com/xxlixin1993/easyGo/gracefulExit"
	redigo "github.com/gomodule/redigo/redis"
)

var pool *redisPool

type (
	redisPool struct {
		// key是redis名 ex:redis_first
		readPool  map[string]*connPoolSet
		writePool map[string]*connPoolSet
	}

	connPoolSet struct {
		// ex:redis_first
		name string
		set  []*connPool
		// set的长度
		count int
	}

	connPool struct {
		maxIdleConns    int
		maxActiveConns  int
		idleTimeout     time.Duration
		maxConnLifeTime time.Duration
		password        string
		modeType        uint8
		pool            *redigo.Pool
	}
)

// 初始化redis
func InitRedis() error {
	pool = newRedisPool()
	redisNames := configure.DefaultStrings("redis.names", []string{})

	if len(redisNames) > 0 {
		for _, redisName := range redisNames {
			for _, mode := range configure.Modes {
				err := pool.connect(redisName, mode)
				if err != nil {
					return err
				}
			}
		}
	}

	// 平滑退出
	gracefulExit.GetExitList().UnShift(pool)

	return nil
}

// Implement ExitInterface
func (rp *redisPool) GetModuleName() string {
	return configure.KRedisModuleName
}

// Implement ExitInterface
func (rp *redisPool) Stop() (err error) {
	err = rp.closeConnSet(rp.writePool)
	err = rp.closeConnSet(rp.readPool)
	return
}

// 关闭一个connPoolSet的所有链接
func (rp *redisPool) closeConnSet(set map[string]*connPoolSet) error {
	for _, connSet := range set {
		for _, conn := range connSet.set {
			err := conn.pool.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 连接配置文件中的redis
func (rp *redisPool) connect(redisName string, mode string) error {
	modeCount := configure.DefaultInt(redisName+"."+mode+".count", 0)
	if modeCount == 0 {
		return errors.New("[cache] redis mode count is empty")
	}

	modeType := configure.KReadMode
	if mode == configure.Modes[1] {
		modeType = configure.KWriteMode
	}

	connReadSet := newConnPoolSet(redisName, modeCount)
	connWriteSet := newConnPoolSet(redisName, modeCount)
	for index := 1; index <= modeCount; index++ {
		indexStr := strconv.Itoa(index)

		conn, err := newConnPool(redisName, mode, indexStr, modeType)
		if err != nil {
			return err
		}

		if modeType == configure.KReadMode {
			connReadSet.set = append(connReadSet.set, conn)
		} else {
			connWriteSet.set = append(connWriteSet.set, conn)
		}
	}

	if modeType == configure.KReadMode {
		rp.readPool[redisName] = connReadSet
	} else {
		rp.writePool[redisName] = connWriteSet
	}

	return nil
}

// 获取一个slave的链接
func (rp *redisPool) getSlaveConn(mysqlName string) (*redigo.Pool, error) {
	if connSet, ok := rp.readPool[mysqlName]; ok {
		return connSet.set[rand.Intn(connSet.count)].pool, nil
	}

	return nil, errors.New("[orm] cant find read connection")
}

// 获取一个master的链接
func (rp *redisPool) getMasterConn(mysqlName string) (*redigo.Pool, error) {
	if connSet, ok := rp.writePool[mysqlName]; ok {
		return connSet.set[rand.Intn(connSet.count)].pool, nil
	}

	return nil, errors.New("[orm] cant find write connection")
}

func newRedisPool() *redisPool {
	return &redisPool{
		readPool:  make(map[string]*connPoolSet),
		writePool: make(map[string]*connPoolSet),
	}
}

func newConnPool(redisName string, mode string, indexStr string, modeType uint8) (*connPool, error) {
	maxIdleConns := configure.DefaultInt(redisName+"."+mode+"."+indexStr+".max_idle_conn", 10)
	maxOpenConns := configure.DefaultInt(redisName+"."+mode+"."+indexStr+".max_open_conn", 100)
	idleTimeout := time.Second * time.Duration(
		configure.DefaultInt(redisName+"."+mode+"."+indexStr+".max_idle_timeout", 10))
	maxLifetime := time.Second * time.Duration(
		configure.DefaultInt(redisName+"."+mode+"."+indexStr+".max_life_time", 300))
	password := configure.DefaultString(redisName+"."+mode+"."+indexStr+".password", "")
	address := configure.DefaultString(redisName+"."+mode+"."+indexStr+".addr", "")
	if address == "" {
		return nil, errors.New("[cache] redis address is empty")
	}

	pool := &redigo.Pool{
		MaxIdle:         maxIdleConns,
		MaxActive:       maxOpenConns,
		IdleTimeout:     idleTimeout,
		MaxConnLifetime: maxLifetime,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", address, redigo.DialPassword(password))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < 3*time.Second {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	connErr := pool.Get().Err()
	if connErr != nil {
		return nil, connErr
	}


	return &connPool{
		maxIdleConns:    maxIdleConns,
		maxActiveConns:  maxOpenConns,
		idleTimeout:     idleTimeout,
		maxConnLifeTime: maxLifetime,
		password:        password,
		pool:            pool,
		modeType:        modeType,
	}, nil
}

func newConnPoolSet(name string, modeCount int) *connPoolSet {
	return &connPoolSet{
		name:  name,
		set:   []*connPool{},
		count: modeCount,
	}
}

// 随机获取从库链接
func GetSlaveConn(redisName string) (*redigo.Pool, error) {
	return pool.getSlaveConn(redisName)
}

// 随机获取主库链接
func GetMasterConn(redisName string) (*redigo.Pool, error) {
	return pool.getMasterConn(redisName)
}
