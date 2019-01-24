package messageQueue

import (
	"github.com/streadway/amqp"
	"github.com/xxlixin1993/easyGo/configure"
	"github.com/xxlixin1993/easyGo/logging"
	"os"
)

type ERRORSTRING string

var ERR_NIL ERRORSTRING
var ERR_CLOSED ERRORSTRING = "Channel/Connection Closed!"
var ERR_FAILED_RECREATE ERRORSTRING = "Recreate Channel/Connection Failed!"
var ERR_FAILED_CREATE ERRORSTRING = "Create Channel Failed"

var ERR_NO_INIT_CONNECTION_POOL = "Need Initialize RabbitMq Connection Pool Firstly!"

type shareConn struct {
	position    int8 //连接在池中的索引, 便于连接失败时，及时清除
	conn        *amqp.Connection
	maxReConn   int  //最大重连次数
	isReachable bool // 连接是否可达
}

//共享连接是对底层amqp连接的封装, 记录连接在池中的位置，便于池中连接数的管理(剔除)
func newShareConn(id int8, conn *amqp.Connection) *shareConn {
	return &shareConn{position: id, conn: conn}
}

//获取连接
func GetConnection() (*shareConn, error) {
	if !rabbitPool.initialized {
		logging.Fatal("Get Connection Failed, Need Initialize RabbitMq Connection Pool Firstly!")
		os.Exit(configure.KInitRabbitMqError)
	}
	conn, index := rabbitPool.getConnection()

	return newShareConn(index, conn), nil
}

func (c *shareConn) ReConnect() (*amqp.Connection, error) {
	var errBack error
	dsn := configure.DefaultString("rabbitMq.dsn", "amqp://guest:guest@localhost:5672")
	for i := 0; i < c.maxReConn; i++ {
		connection, err := amqp.Dial(dsn)
		if err != nil {
			errBack = err
			continue
		}
		rabbitPool.removeByPos(c.position)
		c.position = rabbitPool.put(connection)
		return connection, nil
	}
	return nil, errBack
}
