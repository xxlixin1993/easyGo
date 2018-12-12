package main

import (
	"github.com/xxlixin1993/easyGo"
)

func main() {
	easyGo.InitFrame()
	easyGo.InitMysql()
	easyGo.InitRedis()
	easyGo.InitHTTP(nil)
	easyGo.WaitSignal()
}

