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
	easyGo.InitHTTP(nil)
	InitGRPC()
	easyGo.WaitSignal()
}

type server struct {}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func InitGRPC(){
	rpcServer := &rpc.Server{}
	rpc.InitGRPC(rpcServer)
	pb.RegisterGreeterServer(rpcServer.GetServer(), &server{})
	go rpcServer.ListenAndServe()
}

