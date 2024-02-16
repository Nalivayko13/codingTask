package service

import (
	"fmt"
	"github.com/streadway/amqp"
)

const exchangeKind = "fanout"

type RabbitMQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewRabbitMQConnection(cfg *RabbitMQConfig) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(url)
	return conn, err
}

func NewRabbitExchangeAndQueue(channel *amqp.Channel, exchangeName, queueName string) (string, error) {
	if err := channel.ExchangeDeclare(exchangeName, exchangeKind, true, false, false, false, nil); err != nil {
		return "", err
	}
	queue, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", err
	}
	return queue.Name, err
}

func NewRabbitChannel(connection *amqp.Connection) (*amqp.Channel, error) {
	channel, err := connection.Channel()

	return channel, err
}
