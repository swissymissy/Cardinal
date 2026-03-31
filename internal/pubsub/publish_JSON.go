package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// publish to exchange
func PublishJSON[T any](ctx context.Context, ch *amqp.Channel, exchange, key string, val T) error {
	// convert val to json byte
	bytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Error convert val to json byte: %w", err)
	}
	// publish the message to exchange
	err = ch.PublishWithContext(
		ctx,
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         bytes,
		},
	)
	if err != nil {
		return fmt.Errorf("Failed to publish message to exchange: %w", err)
	}
	return nil
}
