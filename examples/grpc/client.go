package main

import (
	"log"
	"golang.org/x/net/context"
	"github.com/xxlixin1993/easyGo/examples/grpc/pb"
	"github.com/xxlixin1993/easyGo"
	"github.com/xxlixin1993/easyGo/rpc"
)

func main() {
	easyGo.InitFrame()
	easyGo.InitMysql()
	easyGo.InitRedis()
	// TODO exit
	easyGo.InitHTTP(nil)

	testClient()
	easyGo.WaitSignal()
}

func testClient() {
	rpc.InitGRPCClient()
	conn := rpc.GetGRPCClientConn("first")
	if conn == nil {
		log.Fatal("conn is nil")
	}
	c := pb.NewGreeterClient(conn)

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "world test"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

}
