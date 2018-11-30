package configure

import "os"

var Pid = os.Getpid()

// 错误码
const (
	KInitConfigError = iota + 1
	KInitLogError
	KInitMySQLError
	KInitRedisError
)

// 读写两种模式
var Modes = [2]string{"read", "write"}

// 读写分离类型常量
const (
	KReadMode  uint8 = 1
	KWriteMode uint8 = 2
)
