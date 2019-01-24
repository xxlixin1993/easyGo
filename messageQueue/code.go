package messageQueue

// 是对所有错误信息的统一管理
var ErrFailedRecreateConnection = "recreate connection failed"
var ErrFailedRecreateChannel = "recreate channel failed"
var ErrFailedCreate = "create channel failed"
var ErrNoInitConnectionPool = "need initialize rabbitMq connection pool firstly"
var ErrConnectionFailedClose = "connection close failed"
