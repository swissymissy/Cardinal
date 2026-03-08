package handler

import (
	"github.com/swissymissy/Cardinal/internal/database"
)

// struct to hold stateful data
type ApiConfig struct {
	DB 		*database.Queries
	Port	string
}