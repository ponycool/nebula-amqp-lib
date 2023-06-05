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
	Logger *zap.Logger
	Config *Config
}

// InitAmqp 初始化消息队列
func (mq *MQ) InitAmqp() {
	rwLock.Lock()
	defer rwLock.Unlock()

	var err error
	if initialized {
		err = fmt.Errorf("[Init] MQ already initialized")
		mq.Logger.Error(err.Error())
		return
	}

	switch mq.Config.Driver {
	case Rabbit:
		err = initRabbitMQ(mq.Config)
	default:
		err = initRabbitMQ(mq.Config)
	}

	if err != nil {
		defer mq.Logger.Error(err.Error())
		panic(err.Error())
	}

	mq.Logger.Info(fmt.Sprintf("%s connection successful", mq.Config.Driver.String()))
}

// GetAmqp 获取消息队列连接
func GetAmqp() *amqp.Connection {
	return conn
}
