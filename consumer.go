package nebula_amqp_lib

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Consumer struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Tag     string
	Logger  *zap.Logger
	Done    chan error
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.Channel.Cancel(c.Tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.Conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer c.Logger.Info("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.Done
}
