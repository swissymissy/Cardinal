package pubsub 

import (
	"time"

	"github.com/google/uuid"
)
type QueueConfig struct {
	Name string
	Durable bool
	AutoDelete bool 
	Exclusive bool 
	NoWait bool
}

type QueueType int 
const (
	Durable QueueType = iota 
	Transient 
)

type ExchangeType string 
const (
	ExchangeFanout ExchangeType = "fanout"
	ExchangeDirect ExchangeType = "direct"
	ExchangeTopic ExchangeType = "topic"
)

type ChirpEvent struct {
	Body 		string 		`json:"body"`
	Triggerer	uuid.UUID	`json:"triggerer"`
	ChirpID		uuid.UUID	`json:"chirp_id"`
	CreatedAt	time.Time	`json:"created_at"`
}

type AckType int 
const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

