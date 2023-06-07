package nebula_amqp_lib

import (
	"fmt"
	"github.com/streadway/amqp"
)

// 初始化RabbitMQ
func initRabbitMQ(c *Config) (err error) {
	if !c.Enabled {
		err = fmt.Errorf("[InitMQ] rabbitMQ disabled")
		return err
	}

	if len(c.Host) == 0 {
		err = fmt.Errorf("[InitMQ] rabbitMQ host invalid")
		return err
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		c.User,
		c.Password,
		c.Host,
		c.Port)
	conn, err = amqp.Dial(url)
	return
}
