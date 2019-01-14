package server

import (
	"time"
	"net/http"
	"context"
	"bytes"
	"net/http/httputil"
	"fmt"
	"runtime/debug"
	"net"

	"github.com/xxlixin1993/easyGo/configure"
	"github.com/xxlixin1993/easyGo/gracefulExit"
	"github.com/xxlixin1993/easyGo/logging"
	"github.com/gin-gonic/gin"
)

type EasyServer struct {
	host           string
	port           string
	readTimeout    time.Duration
	writeTimeout   time.Duration
	httpServer     *http.Server
	DispatchRouter func(engine *gin.Engine)
	Recover        func() gin.HandlerFunc
	NotFoundRouter func(c *gin.Context)
	Logger         func() gin.HandlerFunc
}

// GetModuleName Implement ExitInterface
func (easyServer *EasyServer) GetModuleName() string {
	return configure.KHTTPModuleName
}

// Stop Implement ExitInterface
func (easyServer *EasyServer) Stop() error {
	quitTimeout := configure.DefaultInt("http.quit_timeout", 30)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(quitTimeout)*time.Second)
	defer cancel()
	return easyServer.httpServer.Shutdown(ctx)
}

// 初始化
func InitHTTPServer(easyServer *EasyServer) {
	if easyServer == nil {
		easyServer = NewEasyServer()
	}
	easyServer.setHTTPServer(easyServer.getGinEngine())
	gracefulExit.GetExitList().UnShift(easyServer)
	go easyServer.listenAndServe()
}

func NewEasyServer() *EasyServer {
	host := configure.DefaultString("http.host", "")
	port := configure.DefaultString("http.port", "")
	readTimeout := configure.DefaultInt("http.read_timeout", 4)
	writeTimeout := configure.DefaultInt("http.write_timeout", 3)
	return &EasyServer{
		host:         host,
		port:         port,
		readTimeout:  time.Duration(readTimeout) * time.Second,
		writeTimeout: time.Duration(writeTimeout) * time.Second,
	}
}

// 设置http server
func (easyServer *EasyServer) setHTTPServer(engine *gin.Engine) {
	easyServer.httpServer = &http.Server{
		Addr:         easyServer.host + ":" + easyServer.port,
		Handler:      engine,
		ReadTimeout:  easyServer.readTimeout,
		WriteTimeout: easyServer.writeTimeout,
	}
}

// 获取Gin引擎
func (easyServer *EasyServer) getGinEngine() *gin.Engine {
	mode := configure.DefaultString("http.mode", "release")
	if mode == gin.ReleaseMode {
		gin.SetMode(mode)
	}
	engine := gin.New()

	if easyServer.Recover == nil {
		easyServer.Recover = easyServer.defaultRecover
	}

	if easyServer.NotFoundRouter == nil {
		easyServer.NotFoundRouter = func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{})
		}
	}

	if easyServer.Logger == nil {
		easyServer.Logger = easyServer.defaultLogger
	}

	if easyServer.DispatchRouter == nil {
		easyServer.DispatchRouter = easyServer.defaultDispatchRouter
	}

	engine.Use(easyServer.Logger(), easyServer.Recover())
	easyServer.DispatchRouter(engine)

	return engine
}

// 监听
func (easyServer *EasyServer) listenAndServe() {
	if err := easyServer.httpServer.ListenAndServe(); err != nil {
		logging.ErrorF("[http] listenAndServe server err:(%s)", err)
	}
}

// 默认路由
func (easyServer *EasyServer) defaultDispatchRouter(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status_code": 200,
			"message":     "hello world",
			"data":        "",
		})
	})
}

// 默认recover处理
func (easyServer *EasyServer) defaultRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				report := panicReport(c, err)
				if !filterErr(err) {
					// mail.SendMail(&mail.Message{
					// 	Subject: "http server panic",
					// 	Body:    report,
					// })
				}
				logging.FatalF("[http server] recovery panic recovered:\n%s", report)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

// 默认日志
func (easyServer *EasyServer) defaultLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)

		path := c.Request.URL.Path

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logging.InfoF("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}

// 过滤一些错误 做特殊处理
func filterErr(err interface{}) bool {
	_, ok := err.(*net.OpError)
	return ok
}

// 错误报告
func panicReport(c *gin.Context, err interface{}) string {
	var buf bytes.Buffer
	req, _ := httputil.DumpRequest(c.Request, false)
	if req != nil {
		buf.Write(req)
	}
	buf.WriteString(fmt.Sprintf("painc: %s\n\n", err))
	stack := debug.Stack()
	buf.Write(stack)
	return buf.String()
}
