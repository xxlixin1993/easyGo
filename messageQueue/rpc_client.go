package messageQueue

import (
	"errors"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type rpcParam struct {
	queueName string
	body      []byte
}
type HandlerClient func(corrId string, delivery amqp.Delivery) (interface{}, error)

// queueName 指定rpc模式下,client
func NewRpcParam(queueName string, body []byte) *rpcParam {
	return &rpcParam{queueName: queueName, body: body}
}

func RpcClient(param *rpcParam, handler HandlerClient) (interface{}, error) {
	shareConn, err := GetConnection()
	if err != nil {
		return nil, err
	}
	channel, err := shareConn.Channel()
	if err != nil {
		return nil, err
	}
	queue, err := channel.QueueDeclare(
		"", // rpcClient 模式，系统会默认生成队列的唯一标识名称
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := channel.Consume(
		queue.Name,
		"", // 消费者为空
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	corrId := uuid.NewV4()

	channel.Publish(
		"", // rpc模式下, 不需要指定交换器，使用默认的即可
		param.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId.String(), // 请求标识
			ReplyTo:       queue.Name,
			Body:          param.body,
		},
	)
	for delivery := range msgs {
		return handler(corrId.String(), delivery)
	}

	return nil, errors.New(ErrRpcClientConsumeFailed)
}
