package nebula_amqp_lib

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Consumer struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	// 客户端的标识符
	ConsumerTag string
	Logger      *zap.Logger
	Done        chan error
}

func (c *Consumer) Consume(message *Message) error {
	var err error

	// 在连接关闭时打印出关闭的原因
	go func() {
		fmt.Printf("closing: %s", <-c.Conn.NotifyClose(make(chan *amqp.Error)))
	}()

	c.Logger.Info("got Connection, getting Channel")
	c.Channel, err = c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	c.Logger.Info(fmt.Sprintf("got Channel, declaring Exchange (%q)", message.Exchange))

	// 交换机声明
	if err = c.Channel.ExchangeDeclare(
		// name of the exchange
		message.Exchange,
		// type
		message.ExchangeType,
		// durable
		message.Durable,
		// delete when complete
		message.AutoDelete,
		// internal,默认false，
		// 设置true后表示当前Exchange是RabbitMQ内部使用，用户所创建的Queue不会消费该类型交换机下的消息
		false,
		// noWait
		false,
		// arguments
		nil,
	); err != nil {
		return fmt.Errorf("exchange Declare: %s", err)
	}

	c.Logger.Info(fmt.Sprintf("declared Exchange, declaring Queue %q", message.QueueName))

	// 声明队列
	queue, err := c.Channel.QueueDeclare(
		// name of the queue
		message.QueueName,
		// durable
		message.Durable,
		// delete when unused
		message.AutoDelete,
		// exclusive
		false,
		// noWait
		false,
		// arguments
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %s", err)
	}

	c.Logger.Info(
		fmt.Sprintf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
			queue.Name,
			queue.Messages,
			queue.Consumers,
			message.BindingKey,
		))

	if err = c.Channel.QueueBind(
		// name of the queue
		queue.Name,
		// bindingKey
		message.BindingKey,
		// sourceExchange
		message.Exchange,
		// noWait
		false,
		// arguments
		nil,
	); err != nil {
		return fmt.Errorf("queue Bind: %s", err)
	}

	c.Logger.Info(fmt.Sprintf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.ConsumerTag))

	deliveries, err := c.Channel.Consume(
		// name
		queue.Name,
		// consumerTag
		c.ConsumerTag,
		// noAck
		false,
		// exclusive
		false,
		// noLocal
		false,
		// noWait
		false,
		// arguments
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Consume: %s", err)
	}

	go c.handle(deliveries, message.CallBack, c.Done)
	return nil
}

// 消息队列处理函数
func (c *Consumer) handle(deliveries <-chan amqp.Delivery, callback func(body []byte), done chan error) {
	for d := range deliveries {
		c.Logger.Info(
			fmt.Sprintf("got %dB delivery: [%v]", len(d.Body), d.DeliveryTag),
			zap.ByteString("Body", d.Body),
		)
		callback(d.Body)
		err := d.Ack(false)
		if err != nil {
			c.Logger.Error(err.Error())
			return
		}
	}
	c.Logger.Info("mq handle: deliveries channel closed")
	done <- nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.Channel.Cancel(c.ConsumerTag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.Conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer c.Logger.Info("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.Done
}
