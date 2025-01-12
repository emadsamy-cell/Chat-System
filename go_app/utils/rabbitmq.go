package utils

import (
	"chat_with_go/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(config.RabbitMQURL)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	_, err = ch.QueueDeclare(config.ChatQueue, true, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	_, err = ch.QueueDeclare(config.MessageQueue, true, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}
