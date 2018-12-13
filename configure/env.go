package configure

import "os"

var Pid = os.Getpid()

// 错误码
const (
	KInitConfigError = iota + 1
	KInitLogError
	KInitMySQLError
	KInitRedisError
	KInitHTTPError
	KInitGRPCError
)

const (
	// mysql模块名
	KMysqlModuleName = "mysqlModule"

	// 日志模块名
	KLogModuleName = "logModule"

	// redis模块名
	KRedisModuleName = "redisModule"

	// http server模块名
	KHTTPModuleName = "httpModule"

	// GRPC 模块名
	KGRPCModuleName = "grpcModule"
)

// 读写两种模式
var Modes = [2]string{"read", "write"}

// 读写分离类型常量
const (
	KReadMode  uint8 = 1
	KWriteMode uint8 = 2
)
