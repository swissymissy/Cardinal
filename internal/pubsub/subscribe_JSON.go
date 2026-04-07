package pubsub

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

// helper
func helperSubscribe[T any](
	conn *amqp.Connection,
	exchangeName,
	queueName,
	key string,
	queueType QueueType,
	exchangeType ExchangeType,
	handler func(T) AckType,
	decoder func([]byte) (T, error),
) error {
	// make sure queue - exchange bound
	ch, queue, err := DeclareAndBind(conn, exchangeName, queueName, key, exchangeType, queueType)
	if err != nil {
		return fmt.Errorf("Queue does not exist or not bound to exchange: %w", err)
	}

	// limit 10 unacknowledged msg at a time
	if err = ch.Qos(10, 0, false); err != nil {
		return fmt.Errorf("Error creating quality of service: %w", err)
	}

	// start consuming
	deliveryChan, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Error consuming msg from queue: %w", err)
	}

	const maxRetries = 3

	// process msg in background
	go func() {
		defer ch.Close()
		for msg := range deliveryChan {
			// decode msg
			data, decodeErr := decoder(msg.Body)
			if decodeErr != nil {
				fmt.Printf("Error decoding message: %s\n", decodeErr)
				msg.Nack(false, false) // discard malformed msg
				continue
			}

			// let handler function decide which ack type and its behavior for each msg
			ackType := handler(data)
			switch ackType {
			case Ack:
				fmt.Println("Ack: Message processed successfully.")
				msg.Ack(false)
			case NackRequeue:
				// get current retry count from headers
				retryCount := 0
				if msg.Headers != nil {
					if val, ok := msg.Headers["x-retry-count"]; ok {
						if count, ok := val.(int32); ok {
							retryCount = int(count)
						}
					}
				}

				if retryCount < maxRetries {
					// republish with incremented retry count
					err := ch.Publish(
						msg.Exchange,
						msg.RoutingKey,
						false,
						false,
						amqp.Publishing{
							ContentType: msg.ContentType,
							Body:        msg.Body,
							Headers: amqp.Table{
								"x-retry-count": int32(retryCount + 1),
							},
						},
					)
					if err != nil {
						fmt.Printf("Failed to republish for retry: %s\n", err)
						msg.Nack(false, false)
					} else {
						fmt.Printf("NackRequeue: retry %d%d\n", retryCount+1, maxRetries)
						msg.Ack(false) // ack original, new copy is in queue
					}
				}
			case NackDiscard:
				fmt.Println("NackDiscard: Handler failed, discarding message.")
				msg.Nack(false, false)
			}
		}
	}()
	return nil
}

// Subscribe wapper
func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchangeName,
	queueName,
	key string,
	queueType QueueType,
	exchangeType ExchangeType,
	handler func(T) AckType,
) error {
	return helperSubscribe(
		conn,
		exchangeName,
		queueName,
		key,
		queueType,
		exchangeType,
		handler,
		func(body []byte) (T, error) {
			var val T
			err := json.Unmarshal(body, &val)
			return val, err
		},
	)
}
