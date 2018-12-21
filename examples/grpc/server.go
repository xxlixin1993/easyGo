package main

import (
	"golang.org/x/net/context"

	"github.com/xxlixin1993/easyGo"
	"github.com/xxlixin1993/easyGo/examples/grpc/pb"
	"github.com/xxlixin1993/easyGo/rpc"
)

func main() {
	easyGo.InitFrame()
	easyGo.InitMysql()
	easyGo.InitRedis()

	InitGRPC()
	easyGo.InitGRPCClient()
	easyGo.WaitSignal()
}

type server struct {}


func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func InitGRPC(){
	rpcServer := rpc.NewServer()
	rpc.InitGRPC(rpcServer)
	pb.RegisterGreeterServer(rpcServer.GetServer(), &server{})
	go rpcServer.ListenAndServe()
}

