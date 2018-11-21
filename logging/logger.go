package logging

import (
	"github.com/xxlixin1993/easyGo/configure"
	"github.com/xxlixin1993/easyGo/gracefulExit"
	"github.com/xxlixin1993/easyGo/utils"
	"errors"
	"fmt"
	"path"
	"runtime"
	"sync"
)

// Log message level
const (
	KLevelFatal = iota
	KLevelError
	KLevelWarnning
	KLevelNotice
	KLevelInfo
	KLevelTrace
	KLevelDebug
)

// 日志模块名
const KLogModuleName = "logModule"

// 日志输出等级
var LevelName = [7]string{"F", "E", "W", "N", "I", "T", "D"}

// 日志实例
var loggerInstance *LogBase

// 日志输出类型
const (
	KOutputFile   = "file"
	KOutputStdout = "stdout"
)

// 日志接口
type ILog interface {
	// 初始化
	Init() error

	// 输出
	OutputLogMsg(msg []byte) error

	Flush()
}

// 日志基础结构体
type LogBase struct {
	mu sync.Mutex
	sync.WaitGroup
	handle  ILog
	message chan []byte
	skip    int
	level   int
}

// Implement ExitInterface
func (l *LogBase) GetModuleName() string {
	return KLogModuleName
}

// Implement ExitInterface
func (l *LogBase) Stop() error {
	close(loggerInstance.message)
	loggerInstance.Wait()
	return nil
}

// 初始化
func InitLog() error {
	outputType := configure.DefaultString("log.output", KOutputStdout)
	level := configure.DefaultInt("log.level", KLevelDebug)

	logger, err := createLogger(outputType, level)
	if err != nil {
		return err
	}

	logger.handle.Init()
	gracefulExit.GetExitList().UnShift(logger)

	go logger.Run()

	return err
}

// 创建一个日志实例
func createLogger(outputType string, level int) (*LogBase, error) {
	switch outputType {
	case KOutputStdout:
		loggerInstance = &LogBase{
			handle:  NewStdoutLog(),
			message: make(chan []byte, 1000),
			skip:    3,
			level:   level,
		}
		return loggerInstance, nil
	case KOutputFile:
		loggerInstance = &LogBase{
			handle:  NewFileLog(),
			message: make(chan []byte, 1000),
			skip:    3,
			level:   level,
		}
		return loggerInstance, nil
	default:
		return nil, errors.New(configure.KUnknownTypeMsg)
	}
}

// 获取一个日志实例
func GetLogger() *LogBase {
	return loggerInstance
}

// 开启一个协程 等待信息
func (l *LogBase) Run() {
	loggerInstance.Add(1)

	for {
		msg, ok := <-l.message
		if !ok {
			l.Done()
			l.handle.Flush()
			break
		}
		err := l.handle.OutputLogMsg(msg)
		if err != nil {
			fmt.Printf("Log: Output handle fail, err:%v\n", err.Error())
		}
	}
}

// 输出
func (l *LogBase) Output(nowLevel int, msg string) {
	now := utils.GetMicTimeFormat()

	l.mu.Lock()
	defer l.mu.Unlock()

	if nowLevel <= l.level {
		_, file, line, ok := runtime.Caller(l.skip)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		msg = fmt.Sprintf("[%s] [%s %s:%d] %s\n", LevelName[nowLevel], now, filename, line, msg)
	}

	l.message <- []byte(msg)
}

func Debug(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelDebug, msg)
}

func DebugF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelDebug, msg)
}

func Trace(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelTrace, msg)
}

func TraceF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelTrace, msg)
}

func Info(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelInfo, msg)
}

func InfoF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelInfo, msg)
}
func Notice(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelNotice, msg)
}

func NoticeF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelNotice, msg)
}

func Warning(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelWarnning, msg)
}

func WarningF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelWarnning, msg)
}

func Error(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelError, msg)
}

func ErrorF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelError, msg)
}

func Fatal(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelFatal, msg)
}

func FatalF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelFatal, msg)
}
