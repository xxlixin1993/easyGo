package main

import (
	"github.com/xxlixin1993/easyGo"
	"github.com/xxlixin1993/easyGo/cache"
	"github.com/xxlixin1993/easyGo/configure"
	"fmt"
	"github.com/xxlixin1993/easyGo/logging"
)

func main() {
	easyGo.InitFrame()
	easyGo.InitMysql()
	easyGo.InitRedis()
	go testRedis()
	easyGo.WaitSignal()
}

func testRedis(){
	readClient, err:= cache.GetClient("redis_first", configure.KReadMode)
	if err != nil {
		logging.Fatal(err)
	}

	writeClient, err := cache.GetClient("redis_first", configure.KWriteMode)
	if err != nil {
		logging.Fatal(err)
	}

	fmt.Println(writeClient.Set("test", "easyGo1"))
	fmt.Println(readClient.Get("test"))
}