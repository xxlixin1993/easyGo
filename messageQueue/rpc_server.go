package messageQueue

import (
	"github.com/streadway/amqp"
)

type HandlerServer func(channel *SafeChannel, delivery amqp.Delivery)
type HandlerClient func(corrId string, delivery amqp.Delivery) interface{}

func RpcServer(queueName string, Handler HandlerServer) error {
	shareConn, err := GetConnection()
	if err != nil {
		return err
	}
	channel, err := shareConn.Channel()
	if err != nil {
		return err
	}
	queue, err := channel.QueueDeclare(
		queueName, // rpcServer发送数据的Queue
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.Qos(
		2,     //最多可以不确认的消息数，未确认消息若超过这个值，broker不会发送消息
		0,     // 对消息大小的控制
		false, //若为true: 表明作用于这个connection上的所有channel和消费者
	)

	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		queue.Name,
		"", //rpc 模式下, 不需要consumer标识
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for delivery := range msgs {
		Handler(channel, delivery) //业务处理
	}
	return nil
}
