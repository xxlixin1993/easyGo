package messageQueue

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/xxlixin1993/easyGo/logging"
)

type producerParam struct {
	exchange     string
	exchangeType string
	queueName    string
	routingKey   string
	reliable     bool
	publishing   amqp.Publishing
}

//创建发布者所需参数
func NewProducerParam(exchange, exchangeType, queueName, routingKey string, reliable bool, publishing amqp.Publishing) *producerParam {

	return &producerParam{
		exchange:     exchange,
		exchangeType: exchangeType,
		queueName:    queueName,
		routingKey:   routingKey,
		reliable:     reliable,
		publishing:   publishing,
	}
}

type producer struct {
	conn    *shareConn
	channel *safeChannel
}

// 创建生产者
func NewProducer() (*producer, error) {
	conn, err := GetConnection()
	if err != nil {
		logging.Warning("Get Connection Failed!", err)
		return nil, err
	}
	return &producer{conn: conn}, nil
}

// 消息发布(对创建交换器，队列，交换器队列绑定, 发布消息的封装)
func (p *producer) Publish(paramInfo *producerParam) error {
	var err error
	p.channel, err = p.conn.Channel()

	if err != nil {
		logging.Warning("Declare Channel Failed!", err)
		return err
	}

	if err = p.channel.ExchangeDeclare(
		paramInfo.exchange,
		paramInfo.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logging.Warning("Declare Exchange Failed!", err)
		return err
	}

	if paramInfo.reliable {
		logging.Info("Start Confirm Mode!")
		if err := p.channel.Confirm(false); err != nil {
			logging.WarningF("The Channel Failed To Be Set Confirm Mode, Reason Is: %s", err)
			return err
		}

		confirms := p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		if err = confirmOne(confirms); err != nil {
			return err
		}

	}

	if err := p.channel.Publish(
		paramInfo.exchange,   // 交换器名
		paramInfo.routingKey, // 绑定键
		false,                // mandatory, true:若没有一个队列与交换器绑定，则将消息返还给生产者 , false:若交换器没有匹配到队列，消息直接丢弃
		false,                // immediate , true:队列没有对应的消费者，则将消息返还给生产者,
		paramInfo.publishing,
	); err != nil {
		return err
	}

	return nil
}

func confirmOne(confirms chan amqp.Confirmation) error {
	logging.Info("Waiting RabbitMqServer Ack..")

	if confirmed := <-confirms; confirmed.Ack {
		logging.InfoF("Message: %d, Accept Successfully!", confirmed.DeliveryTag)
		return nil
	} else {
		errInfo := fmt.Sprintf("Message: %d, Failed Accept !", confirmed.DeliveryTag)
		logging.Fatal(errInfo)
		return errors.New(errInfo)
	}
}

func (p *producer) Shutdown() error {

	err := p.conn.conn.Close()
	if err != nil {
		logging.FatalF("连接 关闭失败!", err)
		return err
	}
	logging.Info("发布者连接成功关闭!")
	return nil
}
