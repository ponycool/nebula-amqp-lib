package nebula_amqp_lib

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sync"
)

var (
	rwLock      sync.RWMutex
	initialized bool
	conn        *amqp.Connection
)

type MQ struct {
	logger *zap.Logger
	config Config
}

// InitAmqp 初始化消息队列
func (mq *MQ) InitAmqp() {
	rwLock.Lock()
	defer rwLock.Unlock()

	var err error
	if initialized {
		err = fmt.Errorf("[Init] MQ already initialized")
		mq.logger.Error(err.Error())
		return
	}

	switch mq.config.Driver {
	case Rabbit:
		err = initRabbitMQ(mq.config)
	default:
		err = initRabbitMQ(mq.config)
	}

	if err != nil {
		defer mq.logger.Error(err.Error())
		panic(err.Error())
	}

	mq.logger.Info(fmt.Sprintf("%s connection successful", mq.config.Driver.String()))
}

// GetAmqp 获取消息队列连接
func GetAmqp() *amqp.Connection {
	return conn
}
