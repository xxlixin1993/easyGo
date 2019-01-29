package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/xxlixin1993/easyGo"
	"github.com/xxlixin1993/easyGo/messageQueue"
)

func main() {
	easyGo.InitFrame()
	easyGo.InitMq()
	err := messageQueue.RpcServer("rpctest", func(channel *messageQueue.SafeChannel, delivery amqp.Delivery) {
		fmt.Println(delivery.ReplyTo)
		err := channel.Publish(
			"",               // exchange
			delivery.ReplyTo, // routing key
			false,            // mandatory
			false,            // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: delivery.CorrelationId,
				Body:          []byte(delivery.Body),
			})
		if err != nil {
			fmt.Println(err)
		}

		delivery.Ack(false)

	})
	if err != nil {
		fmt.Println(err)
	}

}
