package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/xxlixin1993/easyGo"
	"github.com/xxlixin1993/easyGo/messageQueue"
	"log"
)

func main() {
	easyGo.InitFrame()
	easyGo.InitMq()
	rpcParam := messageQueue.NewRpcParam("rpctest", []byte("what fuck!"))
	err := messageQueue.RpcClient(rpcParam, func(corrId string, delivery amqp.Delivery) {
		if corrId == delivery.CorrelationId {
			fmt.Println(string(delivery.Body))
		}

	})

	if err != nil {
		log.Fatal(err)
	}

}
