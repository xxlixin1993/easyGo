package messageQueue

import (
	"github.com/streadway/amqp"
	"github.com/xxlixin1993/easyGo/logging"
)


type consumerParam struct {
	exchange     string //交换器名称
	exchangeType string //交换器分发消息类型
	queueName    string //队列名称
	bindingKey   string // 交换器和队列的路由键
	consumerTag  string //消费者的标识
}

//创建消费者所需参数
func NewConsumerParam(exchange, exchangeType, queueName, bindingKey string) *consumerParam {

	return &consumerParam{
		exchange:     exchange,
		exchangeType: exchangeType,
		queueName:    queueName,
		bindingKey:   bindingKey,
	}
}

//创建消费者
func NewConsumer(consumerTag string) (*consumer, error) {
	shareConn, err := GetConnection()
	if err != nil {
		return nil, err
	}

	consumer := &consumer{conn: shareConn, tag: consumerTag}
	return consumer, nil
}

//消息消费(对创建交换器，队列，交换器队列绑定，消息消费的封装)
func (c *consumer) Consume(paramInfo *consumerParam) (<-chan amqp.Delivery, error) {
	var err error
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	if err = c.channel.ExchangeDeclare(
		paramInfo.exchange,     //交换器名称
		paramInfo.exchangeType, //交换器类型(fanout, direct, topic ,header)
		true,  // 是否持久化
		false, // 是否自动删除交换器(前提是至少有一个交换器/队列与之相连接)
		false, //是否为内部使用，不对外
		false, // 是否等待服务端的确认
		nil,   //额外参数
	); err != nil {
		return nil, err
	}

	queue, err := c.channel.QueueDeclare(
		paramInfo.queueName, //队列名
		true,                // 是否持久化
		false,               // 是否自动删除(前提是至少有一个消费者与队列相连，若消费者断开，则会触发自动删除)
		false,               //是否排他，(true:只有声明该队列的连接才能使用)
		false,               //是否等待确认,(true:假定服务端已创建该队列)
		nil,
	)
	if err != nil {

		return nil, err
	}

	if err = c.channel.QueueBind(
		queue.Name,           // 队列名
		paramInfo.bindingKey, // 绑定键
		paramInfo.exchange,   // 交换器名
		false,                // 不会等待服务端的确认
		nil,                  //额外参数
	); err != nil {
		return nil, err
	}

	deliveries, err := c.channel.Consume(
		queue.Name, // 队列名
		c.tag,      //消费者的标识
		false,      //需要业务段处理完逻辑后，需要自行确认
		false,      //是否排他, true:不会发给其他的消费者
		false,      //The noLocal flag is not supported by RabbitMQ.
		false,      // 是否等待服务端的确认
		nil)        //额外参数

	if err != nil {
		return nil, err
	}
	return deliveries, nil
}

type consumer struct {
	conn    *shareConn
	channel *safeChannel
	tag     string
}

func (c *consumer) shutdown() error {

	err := c.channel.originChannel.Cancel(c.tag, false)
	if err != nil {
		return err
	}

	err = c.conn.conn.Close()
	if err != nil {
		return err
	}

	logging.Info("消费者" + c.tag + "成功关闭!")
	return nil
}
