package messageQueue

import (
	"errors"
	"github.com/streadway/amqp"
)

//安全Channel会做重连，更新连接池，记录日志等


type SafeChannel struct {
	originChannel *amqp.Channel
	position      int8
}

func (c *shareConn) Channel() (*SafeChannel, error) {
	channel, err := c.conn.Channel()
	if err == amqp.ErrClosed {
		//若为连接错误，重试3次
		connection, err := c.ReConnect()
		if err == amqp.ErrClosed {
			return nil, errors.New(ErrFailedRecreateConnection)
		}
		channel, err = connection.Channel()
		if err == amqp.ErrClosed {
			return nil, errors.New(ErrFailedRecreateChannel)
		}
	}

	if err != nil {
		return nil, errors.New(ErrFailedCreate)
	}
	return &SafeChannel{originChannel: channel, position: c.position}, nil
}

func (c *SafeChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return c.originChannel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (c *SafeChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return c.originChannel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (c *SafeChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return c.originChannel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (c *SafeChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return c.originChannel.QueueBind(name, key, exchange, noWait, args)
}

func (c *SafeChannel) QueueUnbind(name, key, exchange string, args amqp.Table) error {
	return c.originChannel.QueueUnbind(name, key, exchange, args)
}

func (c *SafeChannel) QueueInspect(name string) (amqp.Queue, error) {
	return c.originChannel.QueueInspect(name)
}

func (c *SafeChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return c.originChannel.Publish(exchange, key, mandatory, immediate, msg)
}

func (c *SafeChannel) Confirm(noWait bool) error {
	return c.originChannel.Confirm(noWait)
}
func (c *SafeChannel) NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64) {
	return c.originChannel.NotifyConfirm(ack, nack)
}

func (c *SafeChannel) NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation {
	return c.originChannel.NotifyPublish(confirm)
}

func (c *SafeChannel) Qos(prefetchCount, prefetchSize int, global bool) error {
	return c.originChannel.Qos(prefetchCount, prefetchSize, global)
}

func (c *SafeChannel) Close() error {
	return c.originChannel.Close()
}
