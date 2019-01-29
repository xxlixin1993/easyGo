package messageQueue

import (
	"github.com/streadway/amqp"
	"math/rand"
)

type RpcParam struct {
	queueName string
	body      []byte
}



func RpcClient(param *RpcParam, Handler HandlerFunc) error {
	shareConn, err := GetConnection()
	if err != nil {
		return err
	}
	channel, err := shareConn.Channel()
	if err != nil {
		return err
	}
	queue, err := channel.QueueDeclare(
		"", //rpcClient 模式，系统会默认生成队列的唯一标识名称
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		queue.Name,
		"", //消费者为空
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	corId := randomString(32)

	channel.Publish(
		"", // rpc模式下, 不需要指定交换器，使用默认的即可
		param.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corId, // 请求标识
			ReplyTo:       queue.Name,
			Body:          param.body,
		},
	)
	for delivery := range msgs {
		Handler(channel, delivery)
	}
	return nil
}

//生成随机的请求标识
func randomString(len int) string {
	bytes := make([]byte, len)

	for i := 0; i < len; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}
