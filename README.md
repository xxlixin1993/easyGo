# easyGo
A golang frame

# 代码规范

- 文件名用下划线命名
- 函数、变量等使用驼峰命名
- 常量以K或k开头
- 代码注释使用双斜线加空格，ex：`// 这是一个注释`
- 提交时 先用`go fmt`工具格式化

# 使用的开源项目
- orm: github.com/jinzhu/gorm
- redis: github.com/gomodule/redigo/redis
- grpc: google.golang.org/grpc
- proto: github.com/golang/protobuf/proto
- rabbitMq: github.com/streadway/amqp

# 如何使用
## 配置文件
在使用`easyGo`时，你需要先创建`app.ini`，在项目中有`app.ini.example`可以参考。在编译后可以使用选项去指定配置文件目录和模块，默认是加载当前目录下的`app.ini`的`local`模块配置
```
  -c string
        use -c <config file> (default "./app.ini")
  -m string
        Use -m <config mode> (default "local")
  -v    Use -v <current version>
```

### 配置说明
- 指定日志的名称
```
app.log_name = appLog
```

- http配置
```
; HTTP Config
; mode : debug,release,test
http.mode = debug // http日志支持的模式
http.host = 0.0.0.0 // http监听的ip
http.port = 8003 // http监听的port
http.read_timeout = 3 // http服务读取请求超时时间
http.write_timeout = 3 // http服务响应超时时间
http.quit_timeout = 30 // http服务平滑退出超时时间
```

- 日志配置

终端输出
```
log.output = stdout
```
文本输出
```
log.output = file
log.dir = /tmp
```
日志等级 即什么等级的日志需要输出
```
; Log Level
; LevelFatal = iota
; LevelError
; LevelWarn
; LevelNotice
; LevelInfo
; LevelTrace
; LevelDebug
log.level = 7
```

- GRPC Server配置
```
; The network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
; address is ip:port
grpc.server.network = tcp
grpc.server.address = 0.0.0.0:50051

grpc.server.accesslog.enabled = true // 是否开启server端访问日志
grpc.server.accesslog.path = /tmp/grpcserver.access.log // 访问日志输出到哪
```

- GRPC Client配置
```
grpc.client = auth,chat // 所有需要连接的grpc server的别名 如auth服务和chat服务
grpc.client_auth.address = 127.0.0.1:50051 // auth服务server的连接地址
grpc.client_chat.address = 127.0.0.1:50052

grpc.client.accesslog.enabled = true  // 是否开启client端访问日志
grpc.client.accesslog.path = /tmp/grpcclient.access.log
```

- MySQL 配置
```
mysql.names = mysql_first,mysql_second // 所有需要连接mysql server的别名 如需要连接first服务和second服务

mysql_first.read.count = 2 // first服务有几个读机器
mysql_first.write.count = 2 // first服务有几个写机器

// first server第一台配置
mysql_first.read.1.dsn = root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
mysql_first.read.1.max_idle_conn = 10 // 最大闲置的连接数
mysql_first.read.1.max_timeout = 300 // 最大超时时间
mysql_first.read.1.max_open_conn = 200 // 最大打开的连接数

mysql_first.read.2.dsn = root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
mysql_first.read.2.max_idle_conn = 10
mysql_first.read.2.max_timeout = 300
mysql_first.read.2.max_open_conn = 200

mysql_first.write.1.dsn = root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
mysql_first.write.1.max_idle_conn = 10
mysql_first.write.1.max_timeout = 300
mysql_first.write.1.max_open_conn = 200

mysql_first.write.2.dsn = root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
mysql_first.write.2.max_idle_conn = 10
mysql_first.write.2.max_life_time = 300
mysql_first.write.2.max_open_conn = 200
```

- Redis 配置
```
// 所有需要连接redis server的别名 如需要连接redis_group服务
redis.names = redis_group

// 读写各自多少
redis_group.read.count = 2
redis_group.write.count = 2

redis_group.read.1.addr = 127.0.0.1:6379
redis_group.read.1.password =
redis_group.read.1.max_idle_conn = 10
redis_group.read.1.max_timeout = 300
redis_group.read.1.max_open_conn = 200

redis_group.read.2.addr = 127.0.0.1:6379
redis_group.read.2.password =
redis_group.read.2.max_idle_conn = 10
redis_group.read.2.max_timeout = 300
redis_group.read.2.max_open_conn = 200

redis_group.write.1.addr = 127.0.0.1:6379
redis_group.write.1.password =
redis_group.write.1.max_idle_conn = 10
redis_group.write.1.max_timeout = 300
redis_group.write.1.max_open_conn = 200

redis_group.write.2.addr = 127.0.0.1:6379
redis_group.write.2.password =
redis_group.write.2.max_idle_conn = 10
redis_group.write.2.max_idle_timeout = 200
redis_group.write.2.max_life_timeout = 300
redis_group.write.2.max_open_conn = 200
```

## [GRPC example](https://github.com/xxlixin1993/easyGo/tree/master/examples/grpc)

## Redis 使用
[redis example](https://github.com/xxlixin1993/easyGo/tree/master/examples/redis.go)
支持的方法在`redis_command.go`可以直接去看

## [MySQL example](https://github.com/xxlixin1993/easyGo/tree/master/examples/db.go)

支持`gorm`的所有方法，[官方文档](http://gorm.io/docs/)