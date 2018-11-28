package mysql

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/xxlixin1993/easyGo/configure"
	"math/rand"
	"strconv"
	"time"
)

var AllDB *allConnSet

// 读写两种模式
var modes = [2]string{"read", "write"}

// 读写分离类型常量
const (
	KReadMode  uint8 = 1
	KWriteMode uint8 = 2
)

type (
	allConnSet struct {
		// key是数据库名 ex:first
		readSet  map[string]*connSet
		writeSet map[string]*connSet
	}

	connSet struct {
		// ex:first
		name string
		set  []*connection
		// set的长度
		count int
	}

	connection struct {
		// ex:first.read.2
		name string
		dsn  string
		// 最大闲置的连接数
		maxIdleConns int
		// 最大打开的连接数
		maxOpenConns int
		// 读写分离类型 1读 2写
		mode   uint8
		dbConn *gorm.DB
	}
)

// 初始化数据库
func InitDB() error {
	AllDB = newConnAllSet()
	mysqlNames := configure.DefaultStrings("mysql.names", []string{})
	if len(mysqlNames) > 0 {
		for _, mysqlName := range mysqlNames {
			for _, mode := range modes {
				err := AllDB.connect(mysqlName, mode)
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func newConnAllSet() *allConnSet {
	return &allConnSet{
		readSet:  make(map[string]*connSet),
		writeSet: make(map[string]*connSet),
	}
}

// 连接配置文件中的数据库
func (allSet *allConnSet) connect(mysqlName string, mode string) error {
	modeCount := configure.DefaultInt(mysqlName+"."+mode+".count", 0)
	if modeCount == 0 {
		return errors.New("[orm] mysql mode count is empty")
	}

	modeType := KReadMode
	if mode == modes[1] {
		modeType = KWriteMode
	}

	connReadSet := newConnSet(mysqlName, modeCount)
	connWriteSet := newConnSet(mysqlName, modeCount)
	for index := 1; index <= modeCount; index++ {
		indexStr := strconv.Itoa(index)
		dsn := configure.DefaultString(mysqlName+"."+mode+"."+indexStr+".dsn", "")
		if dsn == "" {
			return errors.New("[orm] mysql dsn is empty")
		}

		maxIdleConns := configure.DefaultInt(mysqlName+"."+mode+"."+indexStr+".max_idle_conn", 10)
		maxOpenConns := configure.DefaultInt(mysqlName+"."+mode+"."+indexStr+".max_open_conn", 100)
		maxLifetime := configure.DefaultInt(mysqlName+"."+mode+"."+indexStr+".max_life_time", 300)

		conn, err := gorm.Open("mysql", dsn)
		if err != nil {
			return err
		}
		// 设置闲置的连接数
		conn.DB().SetMaxIdleConns(maxIdleConns)
		// 设置最大打开的连接数
		conn.DB().SetMaxOpenConns(maxOpenConns)
		// 设置连接可被重新使用的最大时间间隔。如果超时，则连接会在重新使用前被关闭。
		conn.DB().SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

		// connection.LogMode(true)

		mysqlConn := &connection{
			name:         mysqlName + ":" + mode + ":" + indexStr,
			dsn:          dsn,
			maxIdleConns: maxIdleConns,
			maxOpenConns: maxOpenConns,
			mode:         modeType,
			dbConn:       conn,
		}

		if modeType == KReadMode {
			connReadSet.set = append(connReadSet.set, mysqlConn)
		} else {
			connWriteSet.set = append(connWriteSet.set, mysqlConn)
		}
	}

	if modeType == KReadMode {
		allSet.readSet[mysqlName] = connReadSet
	} else {
		allSet.writeSet[mysqlName] = connWriteSet
	}

	return nil
}

// 获取一个slave的链接
func (allSet *allConnSet) GetSlaveConn(mysqlName string) (*gorm.DB, error) {
	if connSet, ok := allSet.readSet[mysqlName]; ok {
		return connSet.set[rand.Intn(connSet.count)].dbConn, nil
	}

	return nil, errors.New("[orm] cant find read connection")
}

// 获取一个master的链接
func (allSet *allConnSet) GetMasterConn(mysqlName string) (*gorm.DB, error) {
	if connSet, ok := allSet.writeSet[mysqlName]; ok {
		return connSet.set[rand.Intn(connSet.count)].dbConn, nil
	}

	return nil, errors.New("[orm] cant find write connection")
}

func newConnSet(name string, modeCount int) *connSet {
	return &connSet{
		name:  name,
		set:   []*connection{},
		count: modeCount,
	}
}
