package messageQueue

import (
	"errors"
	"github.com/streadway/amqp"
	"github.com/xxlixin1993/easyGo/configure"
)



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
		return  nil, errors.New(string(ErrNoInitConnectionPool))
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
