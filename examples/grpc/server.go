package main

import (
	"golang.org/x/net/context"

	"github.com/xxlixin1993/easyGo"
	"github.com/xxlixin1993/easyGo/examples/grpc/pb"
	"github.com/xxlixin1993/easyGo/rpc"
)

func main() {
	easyGo.InitFrame()
	InitGRPC()
	easyGo.WaitSignal()
}

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func InitGRPC() {
	rpcServer := rpc.NewServer()
	if err := rpc.InitGRPC(rpcServer); err != nil {
		panic(err)
	}
	pb.RegisterGreeterServer(rpcServer.GetServer(), &server{})
	go rpcServer.ListenAndServe()
}
