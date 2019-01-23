package main

import "github.com/xxlixin1993/easyGo"

func main() {
	easyGo.InitFrame()
	easyGo.InitMq()
	go testRabbitMq()
	easyGo.WaitSignal()
}

func testRabbitMq() {

}