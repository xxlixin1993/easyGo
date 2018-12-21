package rpc

import (
	"errors"

	"github.com/xxlixin1993/easyGo/configure"
	"github.com/xxlixin1993/easyGo/utils/slice"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
	"go.uber.org/zap"
)

var (
	// Logger 全局Logger, 用于记录accesslog
	Logger *zap.Logger
	// 是否启用
	enabled bool
	// 日志客户端 服务端模式
	logCSMode = []string{kServerLog, kClientLog}
)

const (
	kServerLog = "server"
	kClientLog = "client"
)

// InitAccessLog 初始化accesslog的zep
func initAccessLog(cSMode string) error {
	// 判断日志格式是否合法
	if !slice.StrInSlice(cSMode, logCSMode) {
		errors.New("[grpc_logger] log cs mode error")
	}

	path := configure.DefaultString("grpc."+cSMode+".accesslog.path", "/tmp/grpc"+cSMode+".access.log")
	if path == "" {
		return errors.New("[grpc_logger] empty path")
	}

	enabled = configure.DefaultBool("grpc."+cSMode+".accesslog.enabled", false)

	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = []string{path}
	cfg.ErrorOutputPaths = []string{}

	cfg.EncoderConfig.LevelKey = ""
	cfg.EncoderConfig.TimeKey = ""
	cfg.EncoderConfig.NameKey = ""
	cfg.EncoderConfig.CallerKey = ""

	var err error
	Logger, err = cfg.Build()
	if err != nil {
		return err
	}

	return nil
}

// StreamClientInterceptor Grpc客户端Stream Log中间件
func LogUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	if !enabled {
		return emptyUnaryClientInterceptor
	}
	//, zapOpts...
	return grpc_zap.UnaryClientInterceptor(Logger)
}

// StreamClientInterceptor Grpc客户端Unary Log中间件
func LogSteamClientInterceptor() grpc.StreamClientInterceptor {
	if !enabled {
		return emptyStreamClientInterceptor
	}
	//, zapOpts...
	return grpc_zap.StreamClientInterceptor(Logger)
}

// UnaryServerInterceptor Grpc服务端Unary Log中间件
func LogUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	if !enabled {
		return emptyUnaryServerInterceptor
	}
	//, zapOpts...
	return grpc_zap.UnaryServerInterceptor(Logger)
}

// StreamServerInterceptor Grpc服务端Stream Log中间件
func LogStreamServerInterceptor() grpc.StreamServerInterceptor {
	if !enabled {
		return emptyStreamServerInterceptor
	}
	//, zapOpts...
	return grpc_zap.StreamServerInterceptor(Logger)
}
