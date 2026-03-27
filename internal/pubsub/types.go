package pubsub 

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

