package mysql

import (
	"math/rand"
	"strconv"
	"time"
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/xxlixin1993/easyGo/configure"
	"github.com/xxlixin1993/easyGo/gracefulExit"
)

var allDB *mysqlPool



type (
	mysqlPool struct {
		// key是数据库名 ex:mysql_first
		readSet  map[string]*connSet
		writeSet map[string]*connSet
	}

	connSet struct {
		// ex:mysql_first
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
	allDB = newMysqlPool()
	mysqlNames := configure.DefaultStrings("mysql.names", []string{})
	if len(mysqlNames) > 0 {
		for _, mysqlName := range mysqlNames {
			for _, mode := range configure.Modes {
				err := allDB.connect(mysqlName, mode)
				if err != nil {
					return err
				}
			}

		}
	}

	// 平滑退出
	gracefulExit.GetExitList().Push(allDB)
	return nil
}

// Implement ExitInterface
func (mp *mysqlPool) GetModuleName() string {
	return configure.KMysqlModuleName
}

// Implement ExitInterface
func (mp *mysqlPool) Stop() (err error) {
	err = mp.closeConnSet(mp.writeSet)
	err = mp.closeConnSet(mp.readSet)
	return
}

// 关闭一个ConnSet的所有链接
func (mp *mysqlPool) closeConnSet(set map[string]*connSet) error {
	for _, connSet := range set {
		for _, conn := range connSet.set {
			err := conn.dbConn.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 连接配置文件中的数据库
func (mp *mysqlPool) connect(mysqlName string, mode string) error {
	modeCount := configure.DefaultInt(mysqlName+"."+mode+".count", 0)
	if modeCount == 0 {
		return errors.New("[orm] mysql mode count is empty")
	}

	modeType := configure.KReadMode
	if mode == configure.Modes[1] {
		modeType = configure.KWriteMode
	}

	connReadSet := newConnSet(mysqlName, modeCount)
	connWriteSet := newConnSet(mysqlName, modeCount)
	for index := 1; index <= modeCount; index++ {
		indexStr := strconv.Itoa(index)

		mysqlConn, err := newConnection(mysqlName, mode, indexStr, modeType)
		if err != nil {
			return err
		}

		if modeType == configure.KReadMode {
			connReadSet.set = append(connReadSet.set, mysqlConn)
		} else {
			connWriteSet.set = append(connWriteSet.set, mysqlConn)
		}
	}

	if modeType == configure.KReadMode {
		mp.readSet[mysqlName] = connReadSet
	} else {
		mp.writeSet[mysqlName] = connWriteSet
	}

	return nil
}

// 获取一个slave的链接
func (mp *mysqlPool) getSlaveConn(mysqlName string) (*gorm.DB, error) {
	if connSet, ok := mp.readSet[mysqlName]; ok {
		return connSet.set[rand.Intn(connSet.count)].dbConn, nil
	}

	return nil, errors.New("[orm] cant find read connection")
}

// 获取一个master的链接
func (mp *mysqlPool) getMasterConn(mysqlName string) (*gorm.DB, error) {
	if connSet, ok := mp.writeSet[mysqlName]; ok {
		return connSet.set[rand.Intn(connSet.count)].dbConn, nil
	}

	return nil, errors.New("[orm] cant find write connection")
}

func newMysqlPool() *mysqlPool {
	return &mysqlPool{
		readSet:  make(map[string]*connSet),
		writeSet: make(map[string]*connSet),
	}
}

func newConnection(mysqlName string, mode string, indexStr string, modeType uint8) (*connection, error) {
	dsn := configure.DefaultString(mysqlName+"."+mode+"."+indexStr+".dsn", "")
	if dsn == "" {
		return nil, errors.New("[orm] mysql dsn is empty")
	}

	maxIdleConns := configure.DefaultInt(mysqlName+"."+mode+"."+indexStr+".max_idle_conn", 10)
	maxOpenConns := configure.DefaultInt(mysqlName+"."+mode+"."+indexStr+".max_open_conn", 100)
	maxLifetime := configure.DefaultInt(mysqlName+"."+mode+"."+indexStr+".max_life_time", 300)

	conn, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
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

	return mysqlConn, nil
}

func newConnSet(name string, modeCount int) *connSet {
	return &connSet{
		name:  name,
		set:   []*connection{},
		count: modeCount,
	}
}

// 随机获取从库链接
func GetSlaveConn(mysqlName string) (*gorm.DB, error) {
	return allDB.getSlaveConn(mysqlName)
}

// 随机获取主库链接
func GetMasterConn(mysqlName string) (*gorm.DB, error) {
	return allDB.getMasterConn(mysqlName)
}
