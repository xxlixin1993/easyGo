package logging

import "fmt"

type LogStdout struct {
}

func NewStdoutLog() ILog {
	stdout := &LogStdout{}
	return stdout
}

// 初始化
func (s *LogStdout) Init() error {
	return nil
}

// 用标准输出流输出
func (s *LogStdout) OutputLogMsg(msg []byte) error {
	fmt.Print(string(msg))
	return nil
}

// 用标准输出流输出
func (s *LogStdout) Flush() {

}
