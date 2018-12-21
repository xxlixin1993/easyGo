package rpc

import (
	"context"

	"github.com/xxlixin1993/easyGo/configure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	clientMap map[string]*grpc.ClientConn
}

var grpcClient *GRPCClient

// Implement ExitInterface
func (client *GRPCClient) GetModuleName() string {
	return configure.KGRPCClientModule
}

// Implement ExitInterface
func (client *GRPCClient) Stop() error {
	var err error
	if len(client.clientMap) > 0 {
		for _, conn := range client.clientMap {
			err = conn.Close()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// 创建GRPCClient
func newGRPCClient() *GRPCClient {
	return &GRPCClient{
		clientMap: make(map[string]*grpc.ClientConn),
	}
}

// 获取一个GRPCClient connection
func GetGRPCClientConn(clientName string) *grpc.ClientConn {
	if conn, ok := grpcClient.clientMap[clientName]; ok {
		return conn
	}

	return nil
}

// 初始化grpc client
func InitGRPCClient() error {
	clientsName := configure.DefaultStrings("grpc.client", []string{})
	if len(clientsName) == 0 {
		return errors.New("[grpc_client] empty client config")
	}

	// 初始化grpc access log
	initAccessLog(kClientLog)

	// 创建GRPCClient
	grpcClient = newGRPCClient()

	for _, clientName := range clientsName {
		clientAddress := configure.DefaultString("grpc.client_"+clientName+".address", "")
		if clientAddress == "" {
			return errors.New("[grpc_client] client config grpc.client." + clientName + "is empty")
		}

		// TODO load balance
		conn, err := grpc.Dial(clientAddress,
			grpc.WithInsecure(),
			grpc.WithStreamInterceptor(LogSteamClientInterceptor()),
			grpc.WithUnaryInterceptor(LogUnaryClientInterceptor()),
		)
		if err != nil {
			return err
		}

		grpcClient.clientMap[clientName] = conn

	}

	return nil
}

// UnaryClientInterceptor 空的客户端Unary 中间件
func emptyUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoker(ctx, method, req, reply, cc, opts...)
}

// StreamClientInterceptor 空的客户端Stream 中间件
func emptyStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
	streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return streamer(ctx, desc, cc, method, opts...)
}
