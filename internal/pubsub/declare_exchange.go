package pubsub 

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareExchange(ch *amqp.Channel, exchangeName string, exchangeType ExchangeType) error {
	return ch.ExchangeDeclare(
		exchangeName,
		string(exchangeType),
		true,	// durable
		false,	// auto delete
		false,	// internal
		false,	// no wait 
		nil,	// arguments
	)
} 