package handler

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/swissymissy/Cardinal/internal/database"
)

// struct to hold stateful data
type ApiConfig struct {
	DB           *database.Queries
	Port         string
	Platform     string
	JWTSecret    string
	MQConn       *amqp.Connection
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	BaseURL      string
}
