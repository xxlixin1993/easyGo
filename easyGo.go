package easyGo

import (
	"easyGo/configure"
	"easyGo/gracefulExit"
	"easyGo/logging"

	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

)

const (
	KVersion = "0.0.1"
)

func Run() {
	initFrame()

	waitSignal()
}

// 初始化框架
func initFrame() {
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

// waitSignal Wait signal
func waitSignal() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)

	sig := <-sigChan

	logging.TraceF("signal: %d", sig)

	switch sig {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		logging.Trace("exit...")
		gracefulExit.GetExitList().Stop()
	case syscall.SIGUSR1:
		logging.Trace("catch the signal SIGUSR1")
	default:
		logging.Trace("signal do not know")
	}
}
