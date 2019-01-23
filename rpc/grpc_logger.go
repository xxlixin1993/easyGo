package rpc

import (
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/xxlixin1993/easyGo/configure"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var (
	// 用于记录客户端accesslog
	clientLogger *zap.Logger
	// 是否启用
	clientEnabled bool

	// 用于记录服务端accesslog
	serverLogger *zap.Logger
	// 是否启用
	serverEnabled bool
)

// 初始化grpc server access log
func initServerAccessLog() error {
	path := configure.DefaultString("grpc.server.accesslog.path", "/tmp/grpc_server.access.log")
	if path == "" {
		return errors.New("[grpc_logger] empty path")
	}

	serverEnabled = configure.DefaultBool("grpc.server.accesslog.enabled", true)

	var err error
	serverLogger, err = initLog(path)
	if err != nil {
		return err
	}

	return nil
}

// 初始化grpc client access log
func initClientAccessLog() error {
	path := configure.DefaultString("grpc.client.accesslog.path", "/tmp/grpc_client.access.log")
	if path == "" {
		return errors.New("[grpc_logger] empty path")
	}

	clientEnabled = configure.DefaultBool("grpc.client.accesslog.enabled", false)

	var err error
	clientLogger, err = initLog(path)
	if err != nil {
		return err
	}

	return nil
}

// 初始化zap logger
func initLog(path string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = []string{path}
	cfg.ErrorOutputPaths = []string{}

	cfg.EncoderConfig.LevelKey = ""
	cfg.EncoderConfig.TimeKey = ""
	cfg.EncoderConfig.NameKey = ""
	cfg.EncoderConfig.CallerKey = ""

	// info级别不输出code ok, 所以默认是debug
	cfg.Level.SetLevel(
		zapcore.Level(configure.DefaultInt("grpc.client.accesslog.level", int(zapcore.DebugLevel))),
	)

	return cfg.Build()
}

// StreamClientInterceptor Grpc客户端Stream Log中间件
func LogUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	if !clientEnabled {
		return emptyUnaryClientInterceptor
	}
	//, zapOpts...
	return grpc_zap.UnaryClientInterceptor(clientLogger)
}

// StreamClientInterceptor Grpc客户端Unary Log中间件
func LogSteamClientInterceptor() grpc.StreamClientInterceptor {
	if !clientEnabled {
		return emptyStreamClientInterceptor
	}
	//, zapOpts...
	return grpc_zap.StreamClientInterceptor(clientLogger)
}

// UnaryServerInterceptor Grpc服务端Unary Log中间件
func LogUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	if !serverEnabled {
		return emptyUnaryServerInterceptor
	}
	//, zapOpts...
	return grpc_zap.UnaryServerInterceptor(serverLogger)
}

// StreamServerInterceptor Grpc服务端Stream Log中间件
func LogStreamServerInterceptor() grpc.StreamServerInterceptor {
	if !serverEnabled {
		return emptyStreamServerInterceptor
	}
	//, zapOpts...
	return grpc_zap.StreamServerInterceptor(serverLogger)
}
