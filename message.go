package nebula_amqp_lib

type Message struct {
	// 交换机
	Exchange string
	// 交换机类型
	ExchangeType string
	// 队列名称
	QueueName  string
	BindingKey string
	// MQCallback 消息队列回调
	CallBack func(body []byte)
	// 是否持久化
	Durable bool
	// 消费完成后是否自动删除
	AutoDelete bool
}
