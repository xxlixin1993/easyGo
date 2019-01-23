package rpc

import (
	"net"
	"context"

	"github.com/xxlixin1993/easyGo/configure"
	"github.com/xxlixin1993/easyGo/logging"
	"github.com/xxlixin1993/easyGo/gracefulExit"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
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

// 创建Server
func NewServer() *Server {
	return &Server{}
}

// 初始化grpc server
func (s *Server) initServer(network, addr string) error {
	var err error
	s.listener, err = net.Listen(network, addr)
	if err != nil {
		return err
	}

	// 初始化grpc server access log
	initServerAccessLog()

	s.grpcServer = grpc.NewServer(
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				// TODO trace
				LogStreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// TODO trace
				LogUnaryServerInterceptor(),
			),
		))

	return gracefulExit.GetExitList().UnShift(s)
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

// UnaryServerInterceptor 空的服务端Unary 中间件
func emptyUnaryServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	return handler(ctx, req)
}

// StreamServerInterceptor 空的服务端Stream 中间件
func emptyStreamServerInterceptor(srv interface{}, stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) (err error) {
	return handler(srv, stream)
}

// 初始化GRPC并启动
func InitGRPC(grpcServer GRPCServerInterface) error {
	network := configure.DefaultString("grpc.server.network", "tcp")
	addr := configure.DefaultString("grpc.server.address", "0.0.0.0:50051")

	initErr := grpcServer.initServer(network, addr)
	if initErr != nil {
		return initErr
	}

	return nil
}
