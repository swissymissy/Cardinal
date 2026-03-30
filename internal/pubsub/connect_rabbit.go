package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(url string) (*amqp.Connection, *amqp.Channel, error) {
	// establish a connection with RabbitMQ server
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("Can't connect to rabbitmq server: %w", err)
	}

	// open a channel from the connection
	channel, err := connection.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create new channel: %w", err)
	}
	fmt.Println("Successfully connect to RabbitMQ server")
	return connection, channel, nil
}
