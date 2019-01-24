package main

import (
	"github.com/xxlixin1993/easyGo"
	"github.com/xxlixin1993/easyGo/logging"
	"github.com/xxlixin1993/easyGo/messageQueue"
)

func main() {
	easyGo.InitFrame()
	easyGo.InitMq()
	go testRabbitMq()
	easyGo.WaitSignal()
}

func testRabbitMq() {

	var exchange = "test1"
	var exchangeType = "fanout"
	var queueName = "queueX"
	var bindingKey = ""

	consumer, err := messageQueue.NewConsumer("consumer1")

	if err != nil {
		logging.Warning(err)
	}
	param := messageQueue.NewConsumerParam(exchange, exchangeType, queueName, bindingKey)
	deliveries, err := consumer.Consume(param)

	if err != nil {
		logging.Warning(err)
	}

	for delivery := range deliveries {
		logging.InfoF(
			"got %q", delivery.Body,
		)
		delivery.Ack(false)
	}

}
