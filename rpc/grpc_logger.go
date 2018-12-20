package rpc

import (
	"errors"

	"github.com/xxlixin1993/easyGo/configure"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
	"go.uber.org/zap"
)

var (
	// Logger 全局Logger, 用于记录accesslog
	Logger *zap.Logger
	// 是否启用
	enabled bool
)

// InitAccessLog 初始化accesslog的zep
func initAccessLog() error {
	path := configure.DefaultString("grpc.accesslog.path", "/tmp/grpc.access.log")
	if path == "" {
		return errors.New("empty path")
	}

	enabled = configure.DefaultBool("grpc.accesslog.enabled", true)

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


// UnaryServerInterceptor Grpc服务端Log中间件
func LogUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	if !enabled {
		return emptyUnaryServerInterceptor
	}
	//, zapOpts...
	return grpc_zap.UnaryServerInterceptor(Logger)
}

// StreamServerInterceptor Grpc服务端Log中间件
func LogStreamServerInterceptor() grpc.StreamServerInterceptor {
	if !enabled {
		return emptyTraceStreamClientInterceptor
	}
	//, zapOpts...
	return grpc_zap.StreamServerInterceptor(Logger)
}

