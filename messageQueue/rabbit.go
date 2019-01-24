package messageQueue

import (
	"github.com/streadway/amqp"
	"github.com/xxlixin1993/easyGo/logging"
)

//安全Channel会对消费和创建队列和创建新交换器时发生错误时做处理，重连，更新连接池，记录日志等
type safeChannel struct {
	originChannel *amqp.Channel
	maxTry        int
	position      int8
}

func (c *shareConn) Channel() (*safeChannel, ERRORSTRING) {
	channel, err := c.conn.Channel()
	if err != nil {
		if err == amqp.ErrClosed {
			//若为连接错误，重试3次
			connection, err := c.reConnect()
			if err == amqp.ErrClosed {
				logging.WarningF("ReConnect Failed, Because: %v", err)
				return nil, ERR_FAILED_RECREATE
			}
			channel, err = connection.Channel()
			if err == amqp.ErrClosed {
				logging.WarningF("Recreate Channel Failed, Because: %v", err)
				return nil, ERR_FAILED_RECREATE
			}
		} else {
			logging.FatalF("Create Channel Failed, Because: %v", err)
			return nil, ERR_FAILED_CREATE
		}

	}

	return &safeChannel{originChannel: channel, maxTry: c.maxReConn, position: c.position}, nil
}

func (c *safeChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return c.originChannel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (c *safeChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return c.originChannel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (c *safeChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return c.originChannel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (c *safeChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return c.originChannel.QueueBind(name, key, exchange, noWait, args)
}

func (c *safeChannel) QueueUnbind(name, key, exchange string, args amqp.Table) error {
	return c.originChannel.QueueUnbind(name, key, exchange, args)
}

func (c *safeChannel) QueueInspect(name string) (amqp.Queue, error) {
	return c.originChannel.QueueInspect(name)
}

func (c *safeChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return c.originChannel.Publish(exchange, key, mandatory, immediate, msg)
}

func (c *safeChannel) Confirm(noWait bool) error {
	return c.originChannel.Confirm(noWait)
}
func (c *safeChannel) NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64) {
	return c.originChannel.NotifyConfirm(ack, nack)
}

func (c *safeChannel) NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation {
	return c.originChannel.NotifyPublish(confirm)
}

func (c *safeChannel) Qos(prefetchCount, prefetchSize int, global bool) error {
	return c.originChannel.Qos(prefetchCount, prefetchSize, global)
}

func (c *safeChannel) Close() error {
	return c.originChannel.Close()
}
