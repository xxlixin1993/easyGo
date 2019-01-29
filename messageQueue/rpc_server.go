package messageQueue

import "github.com/streadway/amqp"

type HandlerFunc func(channel *safeChannel, delivery amqp.Delivery)

func RpcServer(param *RpcParam, Handler HandlerFunc) error {
	shareConn, err := GetConnection()
	if err != nil {
		return err
	}
	channel, err := shareConn.Channel()
	if err != nil {
		return err
	}
	queue, err := channel.QueueDeclare(
		param.queueName, // rpcServer发送数据的Queue
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
		2,     //做多可以不确认的消息数，未确认消息若超过这个值，broker不会发送消息
		0,     // 对消息大小的控制
		false, //若为true: 表明作用于这个connection上的所有channel和消费者
	)

	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		queue.Name,
		"", //rpc 模式下, 不需要consumer标识
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for d := range msgs {
		Handler(channel, d) //业务处理
	}
	return nil
}
