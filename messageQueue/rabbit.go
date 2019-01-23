package messageQueue

import (
	"github.com/streadway/amqp"
	"github.com/xxlixin1993/easyGo/configure"
	"github.com/xxlixin1993/easyGo/logging"
)

//安全Channel会对消费和创建队列和创建新交换器时发生错误时做处理，重连，更新连接池，记录日志等
type safeChannel struct {
	originChannel *amqp.Channel
	maxTry        int
	position      int8
}

func (c *shareConn) Channel() (*safeChannel, error) {
	channel, err := c.conn.Channel()
	if err != nil {
		//若发生错误，重试3次
		//logging.FatalF("get Channel failed, because: %v", err)
		return nil, err
	}
	return &safeChannel{originChannel: channel, maxTry: c.maxReConn, position: c.position}, nil
}

func (c *safeChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	declareErr := c.originChannel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
	if declareErr != nil {
		//若发生错误，重试3次
		err := c.reConnect()
		if err != nil {
			//logging.FatalF("ReCreate Channel  failed in `ExchangeDeclare step`, because: %v", err)
			return err
		}
		return c.originChannel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
	}
	return nil

}

func (c *safeChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return c.originChannel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (c *safeChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {

	queue, err := c.originChannel.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
	if err != nil {
		//若发生错误，重试3次
		err := c.reConnect()
		if err != nil {

			//logging.FatalF("ReCreate Channel  failed in `QueueDeclare step`, because: %v", err)
		}
		return c.originChannel.QueueDeclare(
			name,
			durable,
			autoDelete,
			exclusive,
			noWait,
			args,
		)
	}
	return queue, nil
}

func (c *safeChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	err := c.originChannel.QueueBind(
		name,
		key,
		exchange,
		noWait,
		args,
	)
	if err != nil {
		//若发生错误，重试3次
		err := c.reConnect()
		if err != nil {
			return err
		}
		err = c.originChannel.QueueBind(
			name,
			key,
			exchange,
			noWait,
			args,
		)
	}
	return err
}

func (c *safeChannel) QueueUnbind(name, key, exchange string, args amqp.Table) error {
	err := c.originChannel.QueueUnbind(
		name,
		key,
		exchange,
		args,
	)
	if err != nil {
		//若发生错误，重试3次
		err := c.reConnect()
		if err != nil {
			return err
		}
		err = c.originChannel.QueueUnbind(
			name,
			key,
			exchange,
			args,
		)
	}
	return err
}

func (c *safeChannel) QueueInspect(name string) (amqp.Queue, error) {
	return  c.originChannel.QueueInspect(name)
}

func (c *safeChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return c.originChannel.Publish(
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)
}

func (c *safeChannel) Confirm(noWait bool) error {
	err := c.originChannel.Confirm(noWait)
	if err != nil {
		//若发生错误，重试3次
		err := c.reConnect()
		if err != nil {
			return err
		}
		err = c.originChannel.Confirm(noWait)
	}
	return err
}
func (c *safeChannel) NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64) {
	return c.originChannel.NotifyConfirm(ack, nack)
}

func (c *safeChannel) NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation {
	return c.originChannel.NotifyPublish(confirm)
}

func (c *safeChannel) Qos(prefetchCount, prefetchSize int, global bool) error {
	err := c.originChannel.Qos(prefetchCount, prefetchSize, global)
	if err != nil {
		//若发生错误，重试3次
		err := c.reConnect()
		if err != nil {
			return err
		}
		err = c.originChannel.Qos(prefetchCount, prefetchSize, global)
	}
	return err
}

func (c *safeChannel) Close() error {
	return c.originChannel.Close()
}

func (c *safeChannel) reConnect() error {
	var errBack error
	dsn := configure.DefaultString("rabbitMq.dsn", "amqp://guest:guest@localhost:5672")
	for i := 0; i < c.maxTry; i++ {
		connection, err := amqp.Dial(dsn)
		if err != nil {
			errBack = err
			continue
		}
		rabbitPool.removeByPos(c.position)
		c.position = rabbitPool.put(connection)

		c.originChannel, err = connection.Channel()
		if err != nil {
			logging.Error("Recreate Channel  failed in `reConnect step`, because: %v", err)
		}
		return nil
	}
	return errBack
}
