package handler

import (
	"github.com/swissymissy/Cardinal/internal/database"
	amqp "github.com/rabbitmq/amqp091-go"
)

// struct to hold stateful data
type ApiConfig struct {
	DB 				*database.Queries
	Port			string
	Platform 		string
	JWTSecret		string
	MQConn			*amqp.Connection 
}