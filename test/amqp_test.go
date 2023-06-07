package test

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ponycool/nebula-amqp-lib"
	"github.com/ponycool/nebula-lib/log"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"
)

var logger *zap.Logger

// 初始化日志
func logInit() {
	if logger != nil {
		return
	}
	logger = log.Init(
		log.SetAppName("nebula-amqp-lib"),
		log.SetDevelopment(true),
		log.SetLevel(zap.DebugLevel),
		log.SetMaxSize(2),
		log.SetMaxBackups(100),
		log.SetMaxAge(30),
	)
	logger.Info("logger initial successful")
}

func TestInit(t *testing.T) {
	logInit()
}

// 测试连接
func TestConnect(t *testing.T) {
	t.Helper()

	_ = godotenv.Load()

	port, _ := strconv.Atoi(os.Getenv("port"))
	// 创建测试用的Config结构体
	config := &nebula_amqp_lib.Config{
		Driver:   nebula_amqp_lib.Rabbit,
		Enabled:  true,
		Host:     os.Getenv("host"),
		Port:     port,
		User:     os.Getenv("user"),
		Password: os.Getenv("password"),
	}

	mq := &nebula_amqp_lib.MQ{
		Logger: logger,
		Config: config,
	}
	mq.InitAmqp()
}

// 测试获取MQ
func TestGetMQ(t *testing.T) {
	t.Helper()

	// 调用GetMQ函数
	fmt.Println("Getting MQ connection...")
	conn := nebula_amqp_lib.GetAmqp()
	if conn == nil {
		panic("MQ connection is nil")
	}
	fmt.Println("MQ connection obtained successfully")
}

// 测试消费队列消息
func TestConsume(t *testing.T) {
	t.Helper()

	conn := nebula_amqp_lib.GetAmqp()
	if conn == nil {
		panic("MQ connection is nil")
	}

	exchange := "amp.topic"
	exchangeType := "topic"
	queueName := "cms.sync.order"
	bindingKey := "cms.sync.order"
	consumerTag := "TestConsumer"

	msg := &nebula_amqp_lib.Message{
		Exchange:     exchange,
		ExchangeType: exchangeType,
		QueueName:    queueName,
		BindingKey:   bindingKey,
		Durable:      true,
		CallBack:     readMessage,
	}
	// 创建测试用的Consumer结构体
	consumer := nebula_amqp_lib.Consumer{
		Conn:        conn,
		Channel:     nil,
		ConsumerTag: consumerTag,
		Logger:      logger,
		Done:        make(chan error),
	}

	// 开始消费
	err := consumer.Consume(msg)
	if err != nil {
		panic(err)
	}

	// 创建一个通道用来接收信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGTSTP)

	// 创建一个标志变量用来控制循环是否继续
	running := true

	go func() {
		for running {
			select {
			case sig := <-sigChan:
				logger.Info("got signal. Aborting...",
					zap.Any("sig", sig),
				)
				// 调用Consumer结构体的Shutdown方法
				fmt.Println("Shutting down MQ...")
				err = consumer.Shutdown()
				if err != nil {
					panic(err)
				}
				fmt.Println("MQ shutdown successfully")
				// 设置标志变量为false
				running = false
				// 退出协程
				return
			default:
				break
			}
		}
	}()

	for running {
		if !running {
			// 退出循环
			break
		}
		select {
		default:
			// 每隔一秒钟检查一次running变量的值
			time.Sleep(1 * time.Second)
		}
	}

	logger.Info("process exit...")
	return

}

// 读取消息
func readMessage(msg []byte) {
	message := string(msg[:])
	logger.Info(message)
}
