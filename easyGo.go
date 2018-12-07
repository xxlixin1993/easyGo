package easyGo

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/xxlixin1993/easyGo/configure"
	"github.com/xxlixin1993/easyGo/gracefulExit"
	"github.com/xxlixin1993/easyGo/logging"
	"github.com/xxlixin1993/easyGo/orm/mysql"
	"github.com/xxlixin1993/easyGo/cache"
	"github.com/xxlixin1993/easyGo/server"
)

const (
	KVersion = "0.0.1"
)

// 初始化框架
func InitFrame() {
	runMode := flag.String("m", "local", "Use -m <config mode>")
	configFile := flag.String("c", "./app.ini", "use -c <config file>")
	version := flag.Bool("v", false, "Use -v <current version>")
	flag.Parse()

	if *version {
		fmt.Println("Version", KVersion, runtime.GOOS+"/"+runtime.GOARCH)
		os.Exit(0)
	}

	gracefulExit.InitExitList()

	configErr := configure.InitConfig(*configFile, *runMode)
	if configErr != nil {
		fmt.Printf("Initialize Configure error : %s", configErr)
		os.Exit(configure.KInitConfigError)
	}

	logErr := logging.InitLog()
	if logErr != nil {
		fmt.Printf("Initialize log error : %s", logErr)
		os.Exit(configure.KInitLogError)
	}

	logging.Trace("Initialized frame")
}

// 初始化mysql
func InitMysql() {
	mysqlErr := mysql.InitDB()
	if mysqlErr != nil {
		fmt.Printf("Initialize mysql error : %s", mysqlErr)
		os.Exit(configure.KInitMySQLError)
	}
}

// 初始化redis
func InitRedis() {
	redisErr := cache.InitRedis()
	if redisErr != nil {
		fmt.Printf("Initialize redis error : %s", redisErr)
		os.Exit(configure.KInitRedisError)
	}
}

// 初始化http
func InitHTTP() {
	server.InitHTTPServer()
}

// WaitSignal Wait signal
func WaitSignal() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)

	sig := <-sigChan

	logging.TraceF("signal: %d", sig)

	switch sig {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		logging.Trace("exit...")
		err := gracefulExit.GetExitList().Stop()
		if err != nil {
			fmt.Printf("gracefulExit error : %s", err)
		}
	case syscall.SIGUSR1:
		logging.Trace("catch the signal SIGUSR1")
	default:
		logging.Trace("signal do not know")
	}
}
