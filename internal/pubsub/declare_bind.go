package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// declare and bind a new queue
func DeclareAndBind(
	connection *amqp.Connection,
	exchangeName,
	queueName,
	key string,
	exchangeType ExchangeType,
	queueType QueueType,
) (*amqp.Channel, amqp.Queue, error) {
	// create new channel for queue and exchange
	ch, err := connection.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("Error create new channel: %w", err)
	}

	// queue config
	durable := false
	autoDel := false
	exclsv := false
	if queueType == Durable {
		durable = true
	} else {
		autoDel = true
		exclsv = true
	}
	qcfg := QueueConfig{
		Name:       queueName,
		Durable:    durable,
		AutoDelete: autoDel,
		Exclusive:  exclsv,
		NoWait:     false,
	}

	// declare exchange
	err = DeclareExchange(ch, exchangeName, exchangeType)
	if err != nil {
		ch.Close()
		return nil, amqp.Queue{}, fmt.Errorf("Failed to declare exchange: %w", err)
	}
	// declare new queue
	newQueue, err := ch.QueueDeclare(
		qcfg.Name,
		qcfg.Durable,
		qcfg.AutoDelete,
		qcfg.Exclusive,
		qcfg.NoWait,
		nil,
	)
	if err != nil {
		ch.Close()
		return nil, amqp.Queue{}, fmt.Errorf("Failed to declare new queue: %w", err)
	}

	// bind queue to exchange
	err = ch.QueueBind(
		qcfg.Name,
		key,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		return nil, amqp.Queue{}, fmt.Errorf("Failed to bind queue to exchange: %w", err)
	}
	return ch, newQueue, nil
}
