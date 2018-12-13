package rpc

import (
	"net"

	"google.golang.org/grpc"
	"github.com/xxlixin1993/easyGo/gracefulExit"
	"github.com/xxlixin1993/easyGo/configure"
	"github.com/xxlixin1993/easyGo/logging"
)

type GRPCInterface interface {
	// 创建一个GRPC server
	initServer(network, addr string) error

	// 获取GRPC server
	GetServer() *grpc.Server

	// 监听
	ListenAndServe()

	// 返回server信息
	GetServiceInfo() map[string]grpc.ServiceInfo
}

type GRPCServerInterface interface {
	GRPCInterface
	gracefulExit.ExitInterface
}

type Server struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

// 初始化grpc server
func (s *Server) initServer(network, addr string) error {
	var err error
	s.listener, err = net.Listen(network, addr)
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer()

	return nil
}

// 获取GRPC server
func (s *Server) GetServer() *grpc.Server {
	return s.grpcServer
}

// 监听
func (s *Server) ListenAndServe() {
	if err := s.grpcServer.Serve(s.listener); err != nil {
		logging.ErrorF("[grpc] listenAndServe server err:(%s)", err)
	}
}

// 返回从服务名称到ServiceInfo的映射。服务名称包括包名称，格式为<package>.<service>.
func (s *Server) GetServiceInfo() map[string]grpc.ServiceInfo {
	return s.grpcServer.GetServiceInfo()
}

// Implement ExitInterface
func (s *Server) GetModuleName() string {
	return configure.KGRPCModuleName
}

// Implement ExitInterface
func (s *Server) Stop() error {
	s.grpcServer.GracefulStop()
	return nil
}

// 初始化GRPC并启动
func InitGRPC(grpcServer GRPCServerInterface) error {
	network := configure.DefaultString("network", "tcp")
	addr := configure.DefaultString("address", "0.0.0.0:50051")

	initErr := grpcServer.initServer(network, addr)
	if initErr != nil {
		return initErr
	}

	// 平滑退出
	exitErr := gracefulExit.GetExitList().UnShift(grpcServer)
	if exitErr != nil {
		return exitErr
	}

	return nil
}
