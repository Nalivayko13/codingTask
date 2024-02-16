package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

type RabbitConfig struct {
	User   string
	Passwd string
	Host   string
	Port   string
}

func NewRabbitMQConnection(cfg *RabbitConfig) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Passwd, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(url)
	return conn, err
}

func NewRabbitChannel(connection *amqp.Connection) (*amqp.Channel, error) {
	channel, err := connection.Channel()

	return channel, err
}

func NewRabbitQueue(channel *amqp.Channel, queueName, exchangeName string) (amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	err = channel.QueueBind(queueName, "#", exchangeName, false, nil)
	if err != nil {
		return amqp.Queue{}, err
	}
	return queue, err
}
