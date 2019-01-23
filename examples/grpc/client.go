package main

import (
	"log"
	"context"

	"github.com/xxlixin1993/easyGo"
	"github.com/xxlixin1993/easyGo/rpc"
	"github.com/xxlixin1993/easyGo/examples/grpc/pb"
)

func main() {
	easyGo.InitFrame()
	easyGo.InitGRPCClient()
	testClient()

	easyGo.WaitSignal()
}

func testClient() {
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
