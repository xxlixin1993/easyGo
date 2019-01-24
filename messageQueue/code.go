package messageQueue

// 是对所有错误信息的统一管理
type ERRORSTRING string

var ErrFailedRecreateConnection ERRORSTRING = "Recreate Connection Failed!"
var ErrFailedRecreateChannel ERRORSTRING = "Recreate Channel Failed!"
var ErrFailedCreate ERRORSTRING = "Create Channel Failed"
var ErrNoInitConnectionPool ERRORSTRING = "Need Initialize RabbitMq Connection Pool Firstly!"
var ErrConnectionFailedClose ERRORSTRING = "Connection Close Failed!"
