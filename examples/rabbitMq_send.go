package main

import (
	"github.com/streadway/amqp"
	"github.com/xxlixin1993/easyGo"
	"github.com/xxlixin1993/easyGo/logging"
	"github.com/xxlixin1993/easyGo/messageQueue"
)

func main() {
	easyGo.InitFrame()
	easyGo.InitMq()
	producer, err := messageQueue.NewProducer()
	if err != nil {
		logging.Warning(err)
	}
	var exchange = "test1"
	var exchangeType = "fanout"
	var routingKey = ""
	var reliable = true
	body := "This is a test message"
	var publishing = amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	}

	param := messageQueue.NewProducerParam(exchange, exchangeType, routingKey, reliable, publishing)
	err = producer.Publish(param)
	if err != nil {
		logging.Warning(err)
	}

}
