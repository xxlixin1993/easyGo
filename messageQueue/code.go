package messageQueue

// 是对所有错误信息的统一管理
type ERRORSTRING string

var ERR_FAILED_RECREATE_CONNECTION ERRORSTRING = "Recreate Connection Failed!"
var ERR_FAILED_RECREATE_CHANNEL ERRORSTRING = "Recreate Channel Failed!"
var ERR_FAILED_CREATE ERRORSTRING = "Create Channel Failed"
var ERR_NO_INIT_CONNECTION_POOL ERRORSTRING = "Need Initialize RabbitMq Connection Pool Firstly!"
var ERR_CONNECTION_FAILED_CLOSE ERRORSTRING =  "Connection Close Failed!"
var ERR_CHANNEL_FAILED_CLOSE ERRORSTRING =  "Channel Close Failed!"
var ERR_DECLARE_QUEUE_FAILED ERRORSTRING = "Declare Queue Failed!";
var ERR_DECLARE_EXCHANGE_FAILED ERRORSTRING = "Declare Exchange Failed!";
var ERR_EXCHANGE_QUEUE_FAILED_BINDING ERRORSTRING = "Exchange Queue Binding Failed!"
var ERR_CONSUME_FAILED ERRORSTRING = "Consume Failed!"